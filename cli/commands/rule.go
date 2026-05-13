package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/matcher"
)

func newRuleCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "rule", Short: "Manage matching rules"}
	cmd.AddCommand(ruleTestCommand())
	cmd.AddCommand(ruleAddCommand())
	cmd.AddCommand(ruleRemoveCommand())
	cmd.AddCommand(ruleListCommand())
	return cmd
}

func ruleTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test <path>",
		Short: "Test directory matching",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
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
			resolver := matcher.NewResolver()
			result, err := resolver.Resolve(ctx, matcher.ResolveInput{
				CurrentDir:     args[0],
				Profiles:       profiles,
				DefaultProfile: cfg.General.DefaultProfile,
			})
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Test path: %s\n", args[0])
			fmt.Fprintf(out, "Match result: %s (%s)\n", result.ProfileName, result.Source)
			return nil
		},
	}
}

func ruleAddCommand() *cobra.Command {
	var dirPattern, urlPattern, profileName string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add matching rule to a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if profileName == "" {
				return fmt.Errorf("--profile is required")
			}
			if dirPattern == "" && urlPattern == "" {
				return fmt.Errorf("--dir or --url is required")
			}
			store, err := defaultStore()
			if err != nil {
				return err
			}
			profile, err := store.GetProfile(context.Background(), profileName)
			if err != nil {
				return err
			}
			if dirPattern != "" {
				profile.Rules.Directory = append(profile.Rules.Directory, dirPattern)
			}
			if urlPattern != "" {
				profile.Rules.URL = append(profile.Rules.URL, urlPattern)
			}
			if err := store.SaveProfile(context.Background(), profile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added rule to profile %s\n", profileName)
			return nil
		},
	}
	cmd.Flags().StringVar(&dirPattern, "dir", "", "Directory pattern to add")
	cmd.Flags().StringVar(&urlPattern, "url", "", "URL pattern to add")
	cmd.Flags().StringVar(&profileName, "profile", "", "Profile name")
	return cmd
}

func ruleRemoveCommand() *cobra.Command {
	var dirPattern, urlPattern, profileName string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove matching rule from a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if profileName == "" {
				return fmt.Errorf("--profile is required")
			}
			if dirPattern == "" && urlPattern == "" {
				return fmt.Errorf("--dir or --url is required")
			}
			store, err := defaultStore()
			if err != nil {
				return err
			}
			profile, err := store.GetProfile(context.Background(), profileName)
			if err != nil {
				return err
			}
			if dirPattern != "" {
				profile.Rules.Directory = removePattern(profile.Rules.Directory, dirPattern)
			}
			if urlPattern != "" {
				profile.Rules.URL = removePattern(profile.Rules.URL, urlPattern)
			}
			if err := store.SaveProfile(context.Background(), profile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed rule from profile %s\n", profileName)
			return nil
		},
	}
	cmd.Flags().StringVar(&dirPattern, "dir", "", "Directory pattern to remove")
	cmd.Flags().StringVar(&urlPattern, "url", "", "URL pattern to remove")
	cmd.Flags().StringVar(&profileName, "profile", "", "Profile name")
	return cmd
}

func ruleListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all matching rules",
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
			for _, p := range profiles {
				for _, d := range p.Rules.Directory {
					fmt.Fprintf(out, "%s  dir: %s\n", p.Profile.Name, d)
				}
				for _, u := range p.Rules.URL {
					fmt.Fprintf(out, "%s  url: %s\n", p.Profile.Name, u)
				}
			}
			return nil
		},
	}
}

func removePattern(patterns []string, target string) []string {
	result := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if p != target {
			result = append(result, p)
		}
	}
	return result
}
