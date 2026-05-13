package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/hook"
)

func newHookCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "hook", Short: "Manage git alias and shell hooks"}
	cmd.AddCommand(hookInstallCommand())
	cmd.AddCommand(hookUninstallCommand())
	return cmd
}

func hookInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install git alias 'git sw'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hook.InstallGitAlias(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed git alias: git sw -> git-switch\n")
			return nil
		},
	}
}

func hookUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove git alias 'git sw'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hook.UninstallGitAlias(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed git alias: git sw\n")
			return nil
		},
	}
}
