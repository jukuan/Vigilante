package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type StateManager struct {
	mu       sync.RWMutex
	offsets  map[string]int64
	filePath string
}

func NewStateManager(filePath string) (*StateManager, error) {
	sm := &StateManager{
		offsets:  make(map[string]int64),
		filePath: filePath,
	}

	if err := sm.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading state: %w", err)
	}

	return sm, nil
}

func (sm *StateManager) GetOffset(path string) int64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.offsets[path]
}

func (sm *StateManager) SetOffset(path string, offset int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.offsets[path] = offset
}

func (sm *StateManager) RemoveFile(path string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.offsets, path)
}

func (sm *StateManager) GetTrackedFiles() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	files := make([]string, 0, len(sm.offsets))
	for path := range sm.offsets {
		files = append(files, path)
	}
	return files
}

func (sm *StateManager) load() error {
	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	return json.Unmarshal(data, &sm.offsets)
}

func (sm *StateManager) Save() error {
	sm.mu.RLock()
	data, err := json.MarshalIndent(sm.offsets, "", "  ")
	sm.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}

	dir := filepath.Dir(sm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating state directory: %w", err)
	}

	tmpFile := sm.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("writing temp state file: %w", err)
	}

	if err := os.Rename(tmpFile, sm.filePath); err != nil {
		return fmt.Errorf("renaming state file: %w", err)
	}

	return nil
}
