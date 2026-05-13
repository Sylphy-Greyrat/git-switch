package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderLoadsMainConfigAndIncludedProfiles(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "config.yaml"), `
general:
  default_profile: personal
  auto_switch: true
  ssh_config_path: ~/.ssh/config

git:
  alias_prefix: sw

include:
  - profiles/*.yaml
`)

	profilesDir := filepath.Join(dir, "profiles")
	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		t.Fatalf("create profiles dir: %v", err)
	}
	writeFile(t, filepath.Join(profilesDir, "personal.yaml"), `
profile:
  name: personal
  description: Personal account
user:
  name: Sylphy
  email: sylphy@example.com
rules:
  directory:
    - ~/projects/personal/*
  url:
    - "github.com:sylphy/*"
`)

	loader := NewLoader(dir)
	cfg, err := loader.LoadMain(ctx)
	if err != nil {
		t.Fatalf("load main config: %v", err)
	}
	if cfg.General.DefaultProfile != "personal" {
		t.Fatalf("expected default profile personal, got %q", cfg.General.DefaultProfile)
	}

	profiles, err := loader.LoadProfiles(ctx)
	if err != nil {
		t.Fatalf("load profiles: %v", err)
	}
	if len(profiles) != 1 || profiles[0].Profile.Name != "personal" {
		t.Fatalf("unexpected profiles: %#v", profiles)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
