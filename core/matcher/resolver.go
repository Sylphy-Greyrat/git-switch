package matcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/sylphy/git-switch/core/config"
)

type MatchSource string

const (
	SourceManual    MatchSource = "manual"
	SourceDirectory MatchSource = "directory"
	SourceURL       MatchSource = "url"
	SourceDefault   MatchSource = "default"
)

type ResolveInput struct {
	CurrentDir     string
	RemoteURL      string
	Profiles       []config.Profile
	DefaultProfile string
}

type ResolveResult struct {
	ProfileName string
	Source      MatchSource
	Rule        string
}

type Resolver interface {
	Resolve(ctx context.Context, input ResolveInput) (ResolveResult, error)
}

type defaultResolver struct {
	directory DirectoryMatcher
	url       URLMatcher
}

func NewResolver() Resolver {
	return &defaultResolver{directory: DirectoryMatcher{}, url: URLMatcher{}}
}

func (r *defaultResolver) Resolve(ctx context.Context, input ResolveInput) (ResolveResult, error) {
	if err := ctx.Err(); err != nil {
		return ResolveResult{}, err
	}

	// 1. Environment variable override
	if manual := os.Getenv("GIT_SWITCH_PROFILE"); manual != "" {
		return ResolveResult{ProfileName: manual, Source: SourceManual, Rule: "env:GIT_SWITCH_PROFILE"}, nil
	}

	// 2. Local state file override
	if stateName := readStateFile(input.CurrentDir); stateName != "" {
		return ResolveResult{ProfileName: stateName, Source: SourceManual, Rule: "state:.git/git-switch-profile"}, nil
	}

	// 3. Directory rules
	directoryRules := make([]Rule, 0)
	for _, profile := range input.Profiles {
		for _, pattern := range profile.Rules.Directory {
			directoryRules = append(directoryRules, Rule{Profile: profile.Profile.Name, Pattern: pattern})
		}
	}
	if matched, err := r.directory.Match(input.CurrentDir, directoryRules); err != nil {
		return ResolveResult{}, err
	} else if matched != "" {
		return ResolveResult{ProfileName: matched, Source: SourceDirectory, Rule: matchedDirRule(input.CurrentDir, directoryRules)}, nil
	}

	// 4. URL rules
	if input.RemoteURL != "" {
		urlRules := make([]Rule, 0)
		for _, profile := range input.Profiles {
			for _, pattern := range profile.Rules.URL {
				urlRules = append(urlRules, Rule{Profile: profile.Profile.Name, Pattern: pattern})
			}
		}
		if matched := r.url.Match(input.RemoteURL, urlRules); matched != "" {
			return ResolveResult{ProfileName: matched, Source: SourceURL, Rule: matchedURLRule(input.RemoteURL, urlRules)}, nil
		}
	}

	// 5. Default
	return ResolveResult{ProfileName: input.DefaultProfile, Source: SourceDefault, Rule: ""}, nil
}

func readStateFile(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, ".git", "git-switch-profile"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func matchedDirRule(dir string, rules []Rule) string {
	for _, rule := range rules {
		pattern, err := config.ExpandHome(rule.Pattern)
		if err != nil {
			continue
		}
		if matchesDirectory(dir, pattern) {
			return rule.Pattern
		}
	}
	return ""
}

func matchedURLRule(remoteURL string, rules []Rule) string {
	m := URLMatcher{}
	for _, rule := range rules {
		if m.Match(remoteURL, []Rule{rule}) != "" {
			return rule.Pattern
		}
	}
	return ""
}
