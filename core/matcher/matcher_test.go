package matcher

import (
	"context"
	"testing"

	"github.com/sylphy/git-switch/core/config"
)

func TestDirectoryMatcher(t *testing.T) {
	matcher := DirectoryMatcher{}
	matched, err := matcher.Match("/Users/sylphy/projects/work/repo", []Rule{{Profile: "work", Pattern: "/Users/sylphy/projects/work/*"}})
	if err != nil {
		t.Fatalf("match directory: %v", err)
	}
	if matched != "work" {
		t.Fatalf("expected work, got %q", matched)
	}
}

func TestURLMatcher(t *testing.T) {
	matcher := URLMatcher{}
	matched := matcher.Match("git@github.com:company/repo.git", []Rule{{Profile: "work", Pattern: "github.com:company/*"}})
	if matched != "work" {
		t.Fatalf("expected work, got %q", matched)
	}
}

func TestResolverPriority(t *testing.T) {
	resolver := NewResolver()
	profiles := []config.Profile{
		{
			Profile: config.ProfileMeta{Name: "personal"},
			User:    config.UserConfig{Name: "Sylphy", Email: "sylphy@example.com"},
			Rules:   config.RulesConfig{URL: []string{"github.com:sylphy/*"}},
		},
		{
			Profile: config.ProfileMeta{Name: "work"},
			User:    config.UserConfig{Name: "Zhang San", Email: "zhangsan@example.com"},
			Rules:   config.RulesConfig{Directory: []string{"/projects/work/*"}},
		},
	}

	result, err := resolver.Resolve(context.Background(), ResolveInput{
		CurrentDir:     "/projects/work/repo",
		RemoteURL:      "git@github.com:sylphy/repo.git",
		Profiles:       profiles,
		DefaultProfile: "personal",
	})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if result.ProfileName != "work" || result.Source != SourceDirectory {
		t.Fatalf("expected directory match work, got %#v", result)
	}
}
