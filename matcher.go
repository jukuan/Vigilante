package main

import (
	"regexp"
)

type PatternMatcher struct {
	patterns []*regexp.Regexp
}

type MatchResult struct {
	Matched     bool
	Sample      string
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

func (pm *PatternMatcher) Match(line string) MatchResult {
	for _, pattern := range pm.patterns {
		loc := pattern.FindStringIndex(line)
		if loc != nil {
			// Extract the matching substring
			matched := line[loc[0]:loc[1]]
			// If very short, grab some context after it
			sample := matched
			if len(sample) < 16 && loc[1] < len(line) {
				end := loc[1] + (16 - len(sample))
				if end > len(line) {
					end = len(line)
				}
				sample = line[loc[0]:end]
			}
			return MatchResult{
				Matched: true,
				Sample:  TruncateTo16(sample),
			}
		}
	}
	return MatchResult{Matched: false}
}

func TruncateTo16(s string) string {
	if len(s) > 16 {
		return s[:16]
	}
	return s
}