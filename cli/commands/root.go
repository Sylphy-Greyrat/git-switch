package commands

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "git-switch",
		Short:   "Manage multiple git users and SSH keys",
		Version: getVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(args)
		},
	}
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newProfileCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newRuleCommand())
	cmd.AddCommand(newSSHCommand())
	cmd.AddCommand(newHookCommand())
	cmd.AddCommand(newCompletionCommand())
	cmd.AddCommand(newTemplateCommand())
	cmd.AddCommand(newUninstallCommand())
	return cmd
}
