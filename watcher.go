package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type LogWatcher struct {
	rule         Rule
	stateManager *StateManager
	matcher      *PatternMatcher
	lineChan     chan<- MatchedLine
	done         chan struct{}
}

func NewLogWatcher(rule Rule, sm *StateManager, lineChan chan<- MatchedLine) *LogWatcher {
	return &LogWatcher{
		rule:         rule,
		stateManager: sm,
		matcher:      NewPatternMatcher(rule.Patterns, rule.SampleLength),
		lineChan:     lineChan,
		done:         make(chan struct{}),
	}
}

func (w *LogWatcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating fsnotify watcher: %w", err)
	}

	if err := watcher.Add(w.rule.LogDir); err != nil {
		_ = watcher.Close()
		return fmt.Errorf("watching directory %s: %w", w.rule.LogDir, err)
	}

	go w.watchLoop(watcher)
	go w.pollLoop()

	log.Printf("[%s] Watching directory: %s (glob: %s)",
		w.rule.Name, w.rule.LogDir, w.rule.FileGlob)
	return nil
}

func (w *LogWatcher) Stop() {
	close(w.done)
}

func (w *LogWatcher) watchLoop(watcher *fsnotify.Watcher) {
	defer func() {
		_ = watcher.Close()
	}()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				w.handleFileEvent(event.Name)
			}
			if event.Op&fsnotify.Remove != 0 {
				w.stateManager.RemoveFile(event.Name)
				log.Printf("[%s] File removed: %s", w.rule.Name, event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[%s] Watcher error: %v", w.rule.Name, err)

		case <-w.done:
			return
		}
	}
}

func (w *LogWatcher) pollLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.scanDirectory()
			w.readAllTrackedFiles()

		case <-w.done:
			return
		}
	}
}

func (w *LogWatcher) scanDirectory() {
	pattern := filepath.Join(w.rule.LogDir, w.rule.FileGlob)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("[%s] Glob error: %v", w.rule.Name, err)
		return
	}

	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}
		w.handleFileEvent(path)
	}
}

func (w *LogWatcher) handleFileEvent(path string) {
	if !w.matchesGlob(path) {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		w.stateManager.RemoveFile(path)
		return
	}

	if info.IsDir() {
		return
	}

	currentOffset := w.stateManager.GetOffset(path)
	if currentOffset == 0 && info.Size() > 0 {
		// New file with content - start from current end
		w.stateManager.SetOffset(path, info.Size())
		log.Printf("[%s] New file detected: %s (starting from end)", w.rule.Name, path)
	}

	w.readFile(path)
}

func (w *LogWatcher) readAllTrackedFiles() {
	for _, path := range w.stateManager.GetTrackedFiles() {
		w.readFile(path)
	}
}

func (w *LogWatcher) readFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		w.stateManager.RemoveFile(path)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	info, err := file.Stat()
	if err != nil {
		w.stateManager.RemoveFile(path)
		return
	}

	currentOffset := w.stateManager.GetOffset(path)
	if currentOffset > info.Size() {
		// File was truncated
		currentOffset = 0
	}

	if currentOffset == info.Size() {
		return
	}

	if _, err := file.Seek(currentOffset, io.SeekStart); err != nil {
		log.Printf("[%s] Seek error for %s: %v", w.rule.Name, path, err)
		return
	}

	buf := make([]byte, info.Size()-currentOffset)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("[%s] Read error for %s: %v", w.rule.Name, path, err)
		return
	}

	if n > 0 {
		w.stateManager.SetOffset(path, currentOffset+int64(n))
		content := string(buf[:n])
		w.processContent(path, content)
	}
}

func (w *LogWatcher) processContent(path, content string) {
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && !strings.HasSuffix(content, "\n") {
		lines = lines[:len(lines)-1]
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		result := w.matcher.Match(line)
		if result.Matched {
			w.lineChan <- MatchedLine{
				RuleName: w.rule.Name,
				FilePath: path,
				Line:     result.Sample, // now contains the matched portion
			}
		}
	}
}

func (w *LogWatcher) matchesGlob(path string) bool {
	dir := filepath.Dir(path)
	if dir != w.rule.LogDir {
		return false
	}

	matched, err := filepath.Match(w.rule.FileGlob, filepath.Base(path))
	if err != nil {
		return false
	}
	return matched
}
