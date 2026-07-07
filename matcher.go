package main

import (
	"regexp"
)

type PatternMatcher struct {
	patterns []*regexp.Regexp
	maxLen   int
}

type MatchResult struct {
	Matched bool
	Sample  string
}

func NewPatternMatcher(patterns []string, maxLen int) *PatternMatcher {
	if maxLen == 0 {
		maxLen = 16
	}
	pm := &PatternMatcher{
		patterns: make([]*regexp.Regexp, 0, len(patterns)),
		maxLen:   maxLen,
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
				Sample:  TruncateLogStr(sample, pm.maxLen),
			}
		}
	}
	return MatchResult{Matched: false}
}

func TruncateLogStr(s string, maxLen ...int) string {
	limit := 16
	if len(maxLen) > 0 && maxLen[0] > 0 {
		limit = maxLen[0]
	}
	if len(s) > limit {
		return s[:limit]
	}
	return s
}
