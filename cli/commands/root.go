package commands

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-switch",
		Short: "Manage multiple git users and SSH keys",
	}
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newProfileCommand())
	return cmd
}
