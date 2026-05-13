package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/config"
	"github.com/sylphy/git-switch/core/matcher"
)

func newProfileCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "profile", Short: "Manage profiles"}
	cmd.AddCommand(profileListCommand())
	cmd.AddCommand(profileShowCommand())
	cmd.AddCommand(profileAddCommand())
	cmd.AddCommand(profileRemoveCommand())
	cmd.AddCommand(profileUseCommand())
	cmd.AddCommand(profileCurrentCommand())
	cmd.AddCommand(profileEditCommand())
	return cmd
}

func profileListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := defaultStore()
			if err != nil {
				return err
			}
			profiles, err := store.ListProfiles(context.Background())
			if err != nil {
				return err
			}
			for _, profile := range profiles {
				fmt.Fprintf(cmd.OutOrStdout(), "%s - %s\n", profile.Profile.Name, profile.Profile.Description)
			}
			return nil
		},
	}
}

func profileShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show profile details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := defaultStore()
			if err != nil {
				return err
			}
			profile, err := store.GetProfile(context.Background(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Profile: %s\nUser: %s <%s>\n", profile.Profile.Name, profile.User.Name, profile.User.Email)
			if profile.SSH != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "SSH Key: %s\n", profile.SSH.KeyFile)
			}
			return nil
		},
	}
}

func profileAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name>",
		Short: "Add profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := defaultStore()
			if err != nil {
				return err
			}
			profile := config.Profile{Profile: config.ProfileMeta{Name: args[0]}, User: config.UserConfig{Name: "Your Name", Email: "your.email@example.com"}}
			if err := store.SaveProfile(context.Background(), profile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created profile %s\n", args[0])
			return nil
		},
	}
}

func profileRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := defaultStore()
			if err != nil {
				return err
			}
			if err := store.DeleteProfile(context.Background(), args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed profile %s\n", args[0])
			return nil
		},
	}
}

func profileUseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Set active profile for current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			// Write state file to .git/git-switch-profile
			gitDir := filepath.Join(dir, ".git")
			if err := os.MkdirAll(gitDir, 0o755); err != nil {
				return fmt.Errorf("create .git directory: %w", err)
			}
			stateFile := filepath.Join(gitDir, "git-switch-profile")
			if err := os.WriteFile(stateFile, []byte(args[0]+"\n"), 0o600); err != nil {
				return fmt.Errorf("write state file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set active profile to %s for %s\n", args[0], dir)
			return nil
		},
	}
}

func profileCurrentCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show current active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name := os.Getenv("GIT_SWITCH_PROFILE"); name != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s (env override)\n", name)
				return nil
			}

			// Check local state file
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			stateFile := filepath.Join(dir, ".git", "git-switch-profile")
			data, err := os.ReadFile(stateFile)
			if err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(data))
				return nil
			}

			// Fall back to resolver
			loader, err := defaultLoader()
			if err != nil {
				return err
			}
			cfg, err := loader.LoadMain(context.Background())
			if err != nil {
				return err
			}
			profiles, err := loader.LoadProfiles(context.Background())
			if err != nil {
				return err
			}
			resolver := matcher.NewResolver()
			result, err := resolver.Resolve(context.Background(), matcher.ResolveInput{
				CurrentDir:     dir,
				RemoteURL:      getRemoteURL(),
				Profiles:       profiles,
				DefaultProfile: cfg.General.DefaultProfile,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result.ProfileName)
			return nil
		},
	}
}

func profileEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit profile in default editor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.DefaultConfigDir()
			if err != nil {
				return err
			}
			path := filepath.Join(dir, "profiles", args[0]+".yaml")
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("profile %q not found at %s", args[0], path)
			}
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			c := exec.Command(editor, path)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
}

func defaultStore() (config.ConfigStore, error) {
	dir, err := config.DefaultConfigDir()
	if err != nil {
		return nil, err
	}
	return config.NewFileStore(dir), nil
}
