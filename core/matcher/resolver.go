package matcher

import (
	"context"
	"os"

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
	if manual := os.Getenv("GIT_SWITCH_PROFILE"); manual != "" {
		return ResolveResult{ProfileName: manual, Source: SourceManual}, nil
	}

	directoryRules := make([]Rule, 0)
	for _, profile := range input.Profiles {
		for _, pattern := range profile.Rules.Directory {
			directoryRules = append(directoryRules, Rule{Profile: profile.Profile.Name, Pattern: pattern})
		}
	}
	if matched, err := r.directory.Match(input.CurrentDir, directoryRules); err != nil {
		return ResolveResult{}, err
	} else if matched != "" {
		return ResolveResult{ProfileName: matched, Source: SourceDirectory}, nil
	}

	if input.RemoteURL != "" {
		urlRules := make([]Rule, 0)
		for _, profile := range input.Profiles {
			for _, pattern := range profile.Rules.URL {
				urlRules = append(urlRules, Rule{Profile: profile.Profile.Name, Pattern: pattern})
			}
		}
		if matched := r.url.Match(input.RemoteURL, urlRules); matched != "" {
			return ResolveResult{ProfileName: matched, Source: SourceURL}, nil
		}
	}

	return ResolveResult{ProfileName: input.DefaultProfile, Source: SourceDefault}, nil
}
