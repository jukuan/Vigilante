package main

import (
	"regexp"
)

type PatternMatcher struct {
	patterns []*regexp.Regexp
}

func NewPatternMatcher(patterns []string) *PatternMatcher {
	pm := &PatternMatcher{
		patterns: make([]*regexp.Regexp, 0, len(patterns)),
	}

	for _, p := range patterns {
		// Escape special regex chars but allow simple patterns
		// Treat as case-insensitive fixed strings with regex support
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			// Fall back to exact match if regex compilation fails
			re = regexp.MustCompile("(?i)" + regexp.QuoteMeta(p))
		}
		pm.patterns = append(pm.patterns, re)
	}

	return pm
}

func (pm *PatternMatcher) Match(line string) bool {
	for _, pattern := range pm.patterns {
		if pattern.MatchString(line) {
			return true
		}
	}
	return false
}

// TruncateTo16 returns first 16 characters, used for alert message
func TruncateTo16(s string) string {
	if len(s) > 16 {
		return s[:16]
	}
	return s
}