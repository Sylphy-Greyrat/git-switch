package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/applier"
	"github.com/sylphy/git-switch/core/matcher"
)

func newSSHCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "ssh", Short: "Manage SSH configuration"}
	cmd.AddCommand(sshConfigCommand())
	return cmd
}

func sshConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Regenerate SSH config.d/git-switch",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			store, err := defaultStore()
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

			profiles, err := store.ListProfiles(ctx)
			if err != nil {
				return err
			}

			dir, err := os.Getwd()
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

			for _, p := range profiles {
				if p.Profile.Name == result.ProfileName {
					sshApplier := applier.NewSSHApplier(filepath.Join(home, ".ssh"))
					if err := sshApplier.ApplySSHConfig(ctx, p); err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Updated SSH config for profile %s\n", p.Profile.Name)
					return nil
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "No matching profile found\n")
			return nil
		},
	}
}
