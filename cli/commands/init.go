package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/config"
	"gopkg.in/yaml.v3"
)

func newInitCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize git-switch configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.DefaultConfigDir()
			if err != nil {
				return err
			}

			configPath := filepath.Join(dir, "config.yaml")
			if _, err := os.Stat(configPath); err == nil && !force {
				return fmt.Errorf("config already exists at %s (use --force to overwrite)", dir)
			}

			if err := os.MkdirAll(filepath.Join(dir, "profiles"), 0o700); err != nil {
				return err
			}

			mainConfig, err := yaml.Marshal(config.DefaultMainConfig())
			if err != nil {
				return err
			}
			if err := os.WriteFile(configPath, mainConfig, 0o600); err != nil {
				return err
			}

			profile := config.Profile{
				Profile: config.ProfileMeta{Name: "personal", Description: "Personal GitHub account"},
				User:    config.UserConfig{Name: "Your Name", Email: "your.email@example.com"},
				SSH:     &config.SSHConfig{KeyFile: "~/.ssh/id_rsa", Hosts: []string{"github.com"}, HostAlias: "github.com-personal"},
				Rules:   config.RulesConfig{Directory: []string{"~/projects/personal/*"}, URL: []string{"github.com:yourusername/*"}},
			}
			data, err := yaml.Marshal(profile)
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dir, "profiles", "personal.yaml"), data, 0o600); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Initialized git-switch config at %s\n", dir)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing configuration")
	return cmd
}
