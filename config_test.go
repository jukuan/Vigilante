package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
inactivity_seconds: 600
state_file: test_state.json
rules:
  - name: test-rule
    log_dir: /var/log
    file_glob: "*.log"
    patterns:
      - "ERROR"
      - "FATAL"
    actions:
      - scripts/email.sh
    cooldown_seconds: 120
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	rule := cfg.Rules[0]
	if rule.Name != "test-rule" {
		t.Errorf("expected name 'test-rule', got %q", rule.Name)
	}
	if rule.CooldownSeconds != 120 {
		t.Errorf("expected cooldown 120, got %d", rule.CooldownSeconds)
	}
	if rule.InactivitySeconds != 600 {
		t.Errorf("expected inactivity 600, got %d", rule.InactivitySeconds)
	}
	if rule.LogDir != "/var/log" {
		t.Errorf("expected log_dir /var/log, got %q", rule.LogDir)
	}
	if len(rule.Patterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(rule.Patterns))
	}
	if len(rule.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(rule.Actions))
	}
}

func TestLoadConfigNoRules(t *testing.T) {
	yamlContent := `inactivity_seconds: 100`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Fatal("expected error for missing rules")
	}
}
