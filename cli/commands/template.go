package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/template"
)

func newTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "template", Short: "Manage project templates"}
	cmd.AddCommand(templateListCommand())
	cmd.AddCommand(templateCreateCommand())
	cmd.AddCommand(templateApplyCommand())
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

func templateCreateCommand() *cobra.Command {
	var profileName, description, gitignore string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create project template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if profileName == "" {
				return fmt.Errorf("--profile is required")
			}
			if err := template.CreateTemplate(args[0], profileName, description, gitignore); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created template %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&profileName, "profile", "", "Profile name to use")
	cmd.Flags().StringVar(&description, "description", "", "Template description")
	cmd.Flags().StringVar(&gitignore, "gitignore", "", ".gitignore content")
	return cmd
}

func templateApplyCommand() *cobra.Command {
	var targetDir string
	cmd := &cobra.Command{
		Use:   "apply <name>",
		Short: "Apply template to project directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if targetDir == "" {
				return fmt.Errorf("--dir is required")
			}
			if err := template.ApplyTemplate(args[0], targetDir); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Applied template %s to %s\n", args[0], targetDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&targetDir, "dir", "", "Target directory")
	return cmd
}
