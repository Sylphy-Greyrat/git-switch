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
	return cmd
}

func ruleTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test <path>",
		Short: "Test directory matching",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
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
