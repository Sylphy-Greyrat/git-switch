package matcher

import "strings"

type URLMatcher struct{}

func (m URLMatcher) Match(remoteURL string, rules []Rule) string {
	normalized := normalizeURL(remoteURL)
	for _, rule := range rules {
		pattern := strings.ReplaceAll(rule.Pattern, ":", "/")
		if matchSimpleWildcard(normalized, pattern) {
			return rule.Profile
		}
	}
	return ""
}

func normalizeURL(remoteURL string) string {
	remoteURL = strings.TrimSuffix(remoteURL, ".git")
	remoteURL = strings.TrimPrefix(remoteURL, "ssh://")
	remoteURL = strings.TrimPrefix(remoteURL, "https://")
	remoteURL = strings.TrimPrefix(remoteURL, "http://")
	remoteURL = strings.TrimPrefix(remoteURL, "git@")
	return strings.ReplaceAll(remoteURL, ":", "/")
}

func matchSimpleWildcard(value, pattern string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return value == pattern
	}
	if !strings.HasPrefix(value, parts[0]) {
		return false
	}
	pos := len(parts[0])
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		idx := strings.Index(value[pos:], part)
		if idx == -1 {
			return false
		}
		pos += idx + len(part)
	}
	last := parts[len(parts)-1]
	return last == "" || strings.HasSuffix(value, last)
}
