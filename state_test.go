package main

import (
	"os"
	"testing"
)

func TestStateManagerBasic(t *testing.T) {
	tmpFile := "test_state.json"
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	}()

	sm, err := NewStateManager(tmpFile)
	if err != nil {
		t.Fatalf("NewStateManager failed: %v", err)
	}

	sm.SetOffset("/var/log/app.log", 1024)
	sm.SetOffset("/var/log/other.log", 2048)

	if off := sm.GetOffset("/var/log/app.log"); off != 1024 {
		t.Errorf("expected offset 1024, got %d", off)
	}

	sm.RemoveFile("/var/log/other.log")
	if off := sm.GetOffset("/var/log/other.log"); off != 0 {
		t.Errorf("expected offset 0 after removal, got %d", off)
	}

	// Save and reload
	if err := sm.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	sm2, err := NewStateManager(tmpFile)
	if err != nil {
		t.Fatalf("NewStateManager reload failed: %v", err)
	}
	if off := sm2.GetOffset("/var/log/app.log"); off != 1024 {
		t.Errorf("reloaded offset mismatch: %d", off)
	}
}
