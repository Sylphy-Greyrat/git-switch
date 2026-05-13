package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/config"
)

func newProfileCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "profile", Short: "Manage profiles"}
	cmd.AddCommand(profileListCommand())
	cmd.AddCommand(profileShowCommand())
	cmd.AddCommand(profileAddCommand())
	cmd.AddCommand(profileRemoveCommand())
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

func defaultStore() (config.ConfigStore, error) {
	dir, err := config.DefaultConfigDir()
	if err != nil {
		return nil, err
	}
	return config.NewFileStore(dir), nil
}
