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
	if !strings.Contains(pattern, "*") {
		return value == pattern
	}
	parts := strings.Split(pattern, "*")
	if len(parts) != 2 {
		return false
	}
	return strings.HasPrefix(value, parts[0]) && strings.HasSuffix(value, parts[1])
}
