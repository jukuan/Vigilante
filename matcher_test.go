package main

import (
	"testing"
)

func TestPatternMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		line     string
		want     bool
	}{
		{
			name:     "exact match",
			patterns: []string{"FATAL"},
			line:     "FATAL: something broke",
			want:     true,
		},
		{
			name:     "case insensitive",
			patterns: []string{"error"},
			line:     "ERROR: disk full",
			want:     true,
		},
		{
			name:     "regex special chars escaped",
			patterns: []string{"[INFO]"},
			line:     "this [INFO] message",
			want:     true,
		},
		{
			name:     "no match",
			patterns: []string{"FATAL"},
			line:     "DEBUG: all good",
			want:     false,
		},
		{
			name:     "multiple patterns, OR logic",
			patterns: []string{"FATAL", "ERROR"},
			line:     "ERROR: something",
			want:     true,
		},
		{
			name:     "pattern with regex (starts with)",
			patterns: []string{"^start"},
			line:     "start of line",
			want:     true,
		},
		{
			name:     "pattern with regex (end with)",
			patterns: []string{"end$"},
			line:     "this is the end",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewPatternMatcher(tt.patterns)
			if got := pm.Match(tt.line); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestTruncateTo16(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"short", "short"},
		{"this is a long string", "this is a long s"},
		{"", ""},
		{"exactly 16 chars!!", "exactly 16 chars"},
	}
	for _, tt := range tests {
		got := TruncateTo16(tt.input)
		if got != tt.want {
			t.Errorf("TruncateTo16(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
