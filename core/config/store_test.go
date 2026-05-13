package config

import (
	"context"
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
