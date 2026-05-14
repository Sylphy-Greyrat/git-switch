package matcher

import (
	"path/filepath"
	"regexp"
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
	matched, err := filepath.Match(pattern, path)
	if err == nil && matched {
		return true
	}
	return matchGlobPattern(path, pattern)
}

func matchGlobPattern(value, pattern string) bool {
	quoted := regexp.QuoteMeta(filepath.Clean(pattern))
	quoted = strings.ReplaceAll(quoted, `\*\*`, `.*`)
	quoted = strings.ReplaceAll(quoted, `\*`, `[^`+regexp.QuoteMeta(string(filepath.Separator))+`]*`)
	matched, err := regexp.MatchString(`^`+quoted+`$`, filepath.Clean(value))
	return err == nil && matched
}
