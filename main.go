package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[vigilate] ")

	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Load configuration
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded %d rules from config", len(cfg.Rules))

	// Initialize state manager
	stateManager, err := NewStateManager(cfg.StateFile)
	if err != nil {
		log.Fatalf("Failed to initialize state manager: %v", err)
	}

	// Initialize alert manager
	alertManager := NewAlertManager()

	// Channel for matched lines
	lineChan := make(chan MatchedLine, 100)

	// Start watchers for each rule
	var watchers []*LogWatcher
	for _, rule := range cfg.Rules {
		watcher := NewLogWatcher(rule, stateManager, lineChan)
		if err := watcher.Start(); err != nil {
			log.Printf("Failed to start watcher for rule %s: %v", rule.Name, err)
			continue
		}
		watchers = append(watchers, watcher)
	}

	// Alert dispatcher goroutine
	go func() {
		for line := range lineChan {
			// Find the rule to get cooldown and actions
			for _, rule := range cfg.Rules {
				if rule.Name == line.RuleName {
					cooldown := time.Duration(rule.CooldownSeconds) * time.Second
					alertManager.AddMatch(line, cooldown, rule.Actions)
					break
				}
			}
		}
	}()

	// Periodic state save
	stateTicker := time.NewTicker(30 * time.Second)
	go func() {
		for range stateTicker.C {
			if err := stateManager.Save(); err != nil {
				log.Printf("Failed to save state: %v", err)
			}
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Service started successfully")

	<-sigChan
	log.Println("Shutting down...")

	// Stop all watchers
	for _, watcher := range watchers {
		watcher.Stop()
	}

	// Flush pending alerts
	alertManager.FlushAll()

	// Final state save
	if err := stateManager.Save(); err != nil {
		log.Printf("Failed to save state on shutdown: %v", err)
	}

	close(lineChan)
	stateTicker.Stop()

	log.Println("Service stopped")
}
