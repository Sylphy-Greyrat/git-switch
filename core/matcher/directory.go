package matcher

import (
	"path/filepath"
	"strings"

	"github.com/sylphy/git-switch/core/config"
)

type Rule struct {
	Profile string
	Pattern string
}

type DirectoryMatcher struct{}

func (m DirectoryMatcher) Match(path string, rules []Rule) (string, error) {
	for _, rule := range rules {
		pattern, err := config.ExpandHome(rule.Pattern)
		if err != nil {
			return "", err
		}
		if matchesDirectory(path, pattern) {
			return rule.Profile, nil
		}
	}
	return "", nil
}

func matchesDirectory(path, pattern string) bool {
	if path == pattern {
		return true
	}
	if strings.Contains(pattern, "**") {
		prefix := strings.Split(pattern, "**")[0]
		return strings.HasPrefix(path, strings.TrimSuffix(prefix, string(filepath.Separator)))
	}
	matched, err := filepath.Match(pattern, path)
	return err == nil && matched
}
