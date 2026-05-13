package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	configDir string
}

func NewLoader(configDir string) *Loader {
	return &Loader{configDir: configDir}
}

func (l *Loader) LoadMain(ctx context.Context) (MainConfig, error) {
	if err := ctx.Err(); err != nil {
		return MainConfig{}, err
	}
	path := filepath.Join(l.configDir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return MainConfig{}, fmt.Errorf("read main config %s: %w", path, err)
	}

	cfg := DefaultMainConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return MainConfig{}, fmt.Errorf("parse main config %s: %w", path, err)
	}
	return cfg, nil
}

func (l *Loader) SaveMain(ctx context.Context, cfg MainConfig) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal main config: %w", err)
	}
	path := filepath.Join(l.configDir, "config.yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write main config %s: %w", path, err)
	}
	return nil
}

func (l *Loader) LoadProfiles(ctx context.Context) ([]Profile, error) {
	cfg, err := l.LoadMain(ctx)
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	for _, include := range cfg.Include {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		pattern := filepath.Join(l.configDir, include)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid include pattern %q: %w", include, err)
		}
		for _, match := range matches {
			profile, err := l.loadProfileFile(ctx, match)
			if err != nil {
				return nil, err
			}
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

func (l *Loader) loadProfileFile(ctx context.Context, path string) (Profile, error) {
	if err := ctx.Err(); err != nil {
		return Profile{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, fmt.Errorf("read profile %s: %w", path, err)
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return Profile{}, fmt.Errorf("parse profile %s: %w", path, err)
	}

	baseDir := filepath.Dir(path)
	for _, include := range profile.Include {
		includePattern := filepath.Join(baseDir, include)
		matches, err := filepath.Glob(includePattern)
		if err != nil {
			return Profile{}, fmt.Errorf("invalid profile include pattern %q: %w", include, err)
		}
		for _, match := range matches {
			if err := mergeProfileInclude(ctx, &profile, match); err != nil {
				return Profile{}, err
			}
		}
	}
	return profile, nil
}

func mergeProfileInclude(ctx context.Context, profile *Profile, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read profile include %s: %w", path, err)
	}

	var included Profile
	if err := yaml.Unmarshal(data, &included); err != nil {
		return fmt.Errorf("parse profile include %s: %w", path, err)
	}

	profile.Rules.Directory = append(profile.Rules.Directory, included.Rules.Directory...)
	profile.Rules.URL = append(profile.Rules.URL, included.Rules.URL...)
	return nil
}
