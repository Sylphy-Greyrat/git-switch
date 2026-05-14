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

func TestURLMatcherMultipleWildcards(t *testing.T) {
	matcher := URLMatcher{}
	matched := matcher.Match("git@gitlab.com:company/team/project.git", []Rule{{Profile: "work", Pattern: "gitlab.com:*/team/*"}})
	if matched != "work" {
		t.Fatalf("expected work, got %q", matched)
	}
}

func TestDirectoryMatcherDoubleStarMiddle(t *testing.T) {
	matcher := DirectoryMatcher{}
	matched, err := matcher.Match("/projects/company/team/repo", []Rule{{Profile: "work", Pattern: "/projects/**/repo"}})
	if err != nil {
		t.Fatalf("match directory: %v", err)
	}
	if matched != "work" {
		t.Fatalf("expected work, got %q", matched)
	}
}

func TestURLMatcherWildcardMatchesNestedPath(t *testing.T) {
	matcher := URLMatcher{}
	matched := matcher.Match("git@github.com:sylphy/team/repo.git", []Rule{{Profile: "personal", Pattern: "github.com:sylphy/*"}})
	if matched != "personal" {
		t.Fatalf("expected personal, got %q", matched)
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
	if result.Rule != "/projects/work/*" {
		t.Fatalf("expected Rule '/projects/work/*', got %q", result.Rule)
	}
}

func TestResolverRuleFieldPopulation(t *testing.T) {
	resolver := NewResolver()
	profiles := []config.Profile{
		{
			Profile: config.ProfileMeta{Name: "personal"},
			User:    config.UserConfig{Name: "Sylphy", Email: "sylphy@example.com"},
			Rules:   config.RulesConfig{URL: []string{"github.com:sylphy/*"}},
		},
	}

	// URL match should populate Rule
	result, err := resolver.Resolve(context.Background(), ResolveInput{
		CurrentDir:     "/some/random/path",
		RemoteURL:      "git@github.com:sylphy/repo.git",
		Profiles:       profiles,
		DefaultProfile: "personal",
	})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if result.Source != SourceURL {
		t.Fatalf("expected URL match, got %s", result.Source)
	}
	if result.Rule != "github.com:sylphy/*" {
		t.Fatalf("expected Rule 'github.com:sylphy/*', got %q", result.Rule)
	}

	// Default should have empty Rule
	result2, err := resolver.Resolve(context.Background(), ResolveInput{
		CurrentDir:     "/some/random/path",
		RemoteURL:      "",
		Profiles:       profiles,
		DefaultProfile: "personal",
	})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if result2.Source != SourceDefault {
		t.Fatalf("expected default match, got %s", result2.Source)
	}
	if result2.Rule != "" {
		t.Fatalf("expected empty Rule for default, got %q", result2.Rule)
	}
}
