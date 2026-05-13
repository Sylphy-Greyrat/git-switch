package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/applier"
	"github.com/sylphy/git-switch/core/config"
	"github.com/sylphy/git-switch/core/matcher"
)

func newSSHCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "ssh", Short: "Manage SSH configuration"}
	cmd.AddCommand(sshConfigCommand())
	cmd.AddCommand(sshAddCommand())
	cmd.AddCommand(sshListCommand())
	cmd.AddCommand(sshTestCommand())
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

func sshAddCommand() *cobra.Command {
	var keyPath, profileName, hostAlias string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add SSH key configuration to a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if profileName == "" {
				return fmt.Errorf("--profile is required")
			}
			if keyPath == "" {
				return fmt.Errorf("--key is required")
			}

			store, err := defaultStore()
			if err != nil {
				return err
			}
			profile, err := store.GetProfile(context.Background(), profileName)
			if err != nil {
				return err
			}

			if profile.SSH == nil {
				profile.SSH = &config.SSHConfig{}
			}
			profile.SSH.KeyFile = keyPath
			if hostAlias != "" {
				profile.SSH.HostAlias = hostAlias
			}

			if err := store.SaveProfile(context.Background(), profile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added SSH key %s to profile %s\n", keyPath, profileName)
			return nil
		},
	}
	cmd.Flags().StringVar(&keyPath, "key", "", "Path to SSH private key")
	cmd.Flags().StringVar(&profileName, "profile", "", "Profile name")
	cmd.Flags().StringVar(&hostAlias, "host-alias", "", "SSH Host alias")
	return cmd
}

func sshListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List SSH keys across profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			loader, err := defaultLoader()
			if err != nil {
				return err
			}
			profiles, err := loader.LoadProfiles(context.Background())
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			found := false
			for _, p := range profiles {
				if p.SSH != nil {
					found = true
					fmt.Fprintf(out, "%s: key=%s", p.Profile.Name, p.SSH.KeyFile)
					if p.SSH.HostAlias != "" {
						fmt.Fprintf(out, " alias=%s", p.SSH.HostAlias)
					}
					fmt.Fprintf(out, " hosts=%v\n", p.SSH.Hosts)
				}
			}
			if !found {
				fmt.Fprintf(out, "No SSH keys configured\n")
			}
			return nil
		},
	}
}

func sshTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test <host>",
		Short: "Test SSH connection to host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := exec.Command("ssh", "-T", "git@"+args[0])
			c.Stdin = os.Stdin
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			return c.Run()
		},
	}
}
