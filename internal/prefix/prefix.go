package prefix

import (
	"regexp"
	"strings"
)

var prefixTokenRE = regexp.MustCompile(`^\s*(\[[^\]]+\]|\([^)]+\)|\{[^}]+\})`)

// ExtractLeading returns all leading [X], (X), {X} tokens from name.
func ExtractLeading(name string) []string {
	s := strings.TrimSpace(name)
	var tokens []string
	for {
		m := prefixTokenRE.FindStringIndex(s)
		if m == nil {
			break
		}
		tok := strings.TrimSpace(s[m[0]:m[1]])
		tokens = append(tokens, tok)
		s = strings.TrimLeft(s[m[1]:], " ._-")
	}
	return tokens
}

// DetectRepeated returns tokens that appear in 2+ names.
func DetectRepeated(names []string) map[string]struct{} {
	counts := map[string]int{}
	for _, name := range names {
		for _, tok := range ExtractLeading(name) {
			counts[tok]++
		}
	}
	result := map[string]struct{}{}
	for tok, cnt := range counts {
		if cnt >= 2 {
			result[tok] = struct{}{}
		}
	}
	return result
}

// Strip removes removable prefix tokens from the front of name.
func Strip(name string, removable map[string]struct{}) string {
	if len(removable) == 0 {
		return strings.TrimSpace(name)
	}
	s := strings.TrimSpace(name)
	for {
		m := prefixTokenRE.FindStringIndex(s)
		if m == nil {
			break
		}
		tok := strings.TrimSpace(s[m[0]:m[1]])
		if _, ok := removable[tok]; !ok {
			break
		}
		s = strings.TrimLeft(s[m[1]:], " ._-")
	}
	return strings.TrimSpace(s)
}

// CleanTitle strips prefixes and known extension, returning a display title.
func CleanTitle(name string, removable map[string]struct{}, knownExts map[string]struct{}) string {
	s := Strip(name, removable)
	dot := strings.LastIndex(s, ".")
	if dot >= 0 {
		ext := strings.ToLower(s[dot:])
		if _, ok := knownExts[ext]; ok {
			s = s[:dot]
		}
	}
	return strings.TrimSpace(s)
}
