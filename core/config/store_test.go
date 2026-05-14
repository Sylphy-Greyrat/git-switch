package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStoreSaveGetListDeleteProfile(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	store := NewFileStore(dir)

	profile := Profile{
		Profile: ProfileMeta{Name: "work", Description: "Work account"},
		User:    UserConfig{Name: "Zhang San", Email: "zhangsan@example.com"},
		Rules:   RulesConfig{},
	}

	if err := store.SaveProfile(ctx, profile); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	loaded, err := store.GetProfile(ctx, "work")
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if loaded.User.Name != "Zhang San" {
		t.Fatalf("expected user Zhang San, got %q", loaded.User.Name)
	}

	profiles, err := store.ListProfiles(ctx)
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}

	if err := store.DeleteProfile(ctx, "work"); err != nil {
		t.Fatalf("delete profile: %v", err)
	}

	if _, err := store.GetProfile(ctx, "work"); err == nil {
		t.Fatal("expected error for deleted profile")
	}
}

func TestFileStoreRejectsProfilePathTraversal(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	store := NewFileStore(dir)
	outsidePath := filepath.Join(dir, "..", "evil.yaml")

	profile := Profile{
		Profile: ProfileMeta{Name: "../evil"},
		User:    UserConfig{Name: "Evil", Email: "evil@example.com"},
	}

	if err := store.SaveProfile(ctx, profile); !errors.Is(err, ErrInvalidProfileName) {
		t.Fatalf("expected ErrInvalidProfileName from SaveProfile, got %v", err)
	}
	if _, err := os.Stat(outsidePath); !os.IsNotExist(err) {
		t.Fatalf("path traversal created file outside profiles dir: %v", err)
	}

	if _, err := store.GetProfile(ctx, "../evil"); !errors.Is(err, ErrInvalidProfileName) {
		t.Fatalf("expected ErrInvalidProfileName from GetProfile, got %v", err)
	}
	if err := store.DeleteProfile(ctx, "../evil"); !errors.Is(err, ErrInvalidProfileName) {
		t.Fatalf("expected ErrInvalidProfileName from DeleteProfile, got %v", err)
	}
}

func TestFileStoreRejectsControlCharacterProfileName(t *testing.T) {
	ctx := context.Background()
	store := NewFileStore(t.TempDir())
	profile := Profile{
		Profile: ProfileMeta{Name: "bad\nHost *"},
		User:    UserConfig{Name: "Evil", Email: "evil@example.com"},
	}

	if err := store.SaveProfile(ctx, profile); !errors.Is(err, ErrInvalidProfileName) {
		t.Fatalf("expected ErrInvalidProfileName for control character, got %v", err)
	}
}
