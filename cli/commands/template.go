package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/template"
)

func newTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "template", Short: "Manage project templates"}
	cmd.AddCommand(templateListCommand())
	return cmd
}

func templateListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := template.ListTemplates()
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No templates found\n")
				return nil
			}
			for _, name := range names {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
			}
			return nil
		},
	}
}
