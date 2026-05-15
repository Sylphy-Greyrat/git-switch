package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/applier"
	"github.com/sylphy/git-switch/core/config"
	"github.com/sylphy/git-switch/core/matcher"
)

func newStatusCommand() *cobra.Command {
	var quiet bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current configuration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			loader, err := defaultLoader()
			if err != nil {
				return err
			}

			cfg, err := loader.LoadMain(ctx)
			if err != nil {
				return err
			}

			profiles, err := loader.LoadProfiles(ctx)
			if err != nil {
				return err
			}

			remoteURL := getRemoteURL()

			resolver := matcher.NewResolver()
			result, err := resolver.Resolve(ctx, matcher.ResolveInput{
				CurrentDir:     dir,
				RemoteURL:      remoteURL,
				Profiles:       profiles,
				DefaultProfile: cfg.General.DefaultProfile,
			})
			if err != nil {
				return err
			}

			// Apply the matched profile
			for _, p := range profiles {
				if p.Profile.Name == result.ProfileName {
					home, err := os.UserHomeDir()
					if err != nil {
						return err
					}
					profileApplier := applier.NewProfileApplier(
						filepath.Join(home, ".gitconfig"),
						filepath.Join(home, ".ssh"),
					)
					if err := profileApplier.ApplyGitConfig(ctx, p); err != nil {
						return fmt.Errorf("apply git config: %w", err)
					}
					if p.SSH != nil {
						if err := profileApplier.ApplySSHConfig(ctx, p); err != nil {
							return fmt.Errorf("apply ssh config: %w", err)
						}
					}
					break
				}
			}

			if quiet {
				return nil
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Working directory: %s\n", dir)
			if remoteURL != "" {
				fmt.Fprintf(out, "Remote URL: %s\n", remoteURL)
			}
			fmt.Fprintf(out, "\nActive profile: %s\n", result.ProfileName)
			fmt.Fprintf(out, "Match source: %s\n", result.Source)

			for _, p := range profiles {
				if p.Profile.Name == result.ProfileName {
					fmt.Fprintf(out, "\nUser config:\n")
					fmt.Fprintf(out, "  user.name = %s\n", p.User.Name)
					fmt.Fprintf(out, "  user.email = %s\n", p.User.Email)
					if p.SSH != nil {
						fmt.Fprintf(out, "\nSSH config:\n")
						fmt.Fprintf(out, "  Key file: %s\n", p.SSH.KeyFile)
						if p.SSH.HostAlias != "" {
							fmt.Fprintf(out, "  Host alias: %s\n", p.SSH.HostAlias)
						}
					}
					break
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Silent mode, apply profile without output")
	return cmd
}

func getRemoteURL() string {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func defaultLoader() (*config.Loader, error) {
	dir, err := config.DefaultConfigDir()
	if err != nil {
		return nil, err
	}
	return config.NewLoader(dir), nil
}
