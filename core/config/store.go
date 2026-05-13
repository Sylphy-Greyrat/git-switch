package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrProfileNotFound = errors.New("profile not found")

type ConfigStore interface {
	ListProfiles(ctx context.Context) ([]Profile, error)
	GetProfile(ctx context.Context, name string) (Profile, error)
	SaveProfile(ctx context.Context, profile Profile) error
	DeleteProfile(ctx context.Context, name string) error
}

type FileStore struct {
	configDir string
}

func NewFileStore(configDir string) *FileStore {
	return &FileStore{configDir: configDir}
}

func (s *FileStore) ListProfiles(ctx context.Context) ([]Profile, error) {
	profilesDir := filepath.Join(s.configDir, "profiles")
	entries, err := os.ReadDir(profilesDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read profiles directory %s: %w", profilesDir, err)
	}

	profiles := make([]Profile, 0, len(entries))
	for _, entry := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}
		profile, err := s.GetProfile(ctx, entry.Name()[:len(entry.Name())-len(".yaml")])
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (s *FileStore) GetProfile(ctx context.Context, name string) (Profile, error) {
	if err := ctx.Err(); err != nil {
		return Profile{}, err
	}
	path := filepath.Join(s.configDir, "profiles", name+".yaml")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	if err != nil {
		return Profile{}, fmt.Errorf("read profile %s: %w", path, err)
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return Profile{}, fmt.Errorf("parse profile %s: %w", path, err)
	}
	return profile, nil
}

func (s *FileStore) SaveProfile(ctx context.Context, profile Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if profile.Profile.Name == "" {
		return errors.New("profile name is required")
	}

	profilesDir := filepath.Join(s.configDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0o700); err != nil {
		return fmt.Errorf("create profiles directory %s: %w", profilesDir, err)
	}

	data, err := yaml.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshal profile %s: %w", profile.Profile.Name, err)
	}

	path := filepath.Join(profilesDir, profile.Profile.Name+".yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write profile %s: %w", path, err)
	}
	return nil
}

func (s *FileStore) DeleteProfile(ctx context.Context, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path := filepath.Join(s.configDir, "profiles", name+".yaml")
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrProfileNotFound, name)
		}
		return fmt.Errorf("delete profile %s: %w", path, err)
	}
	return nil
}
