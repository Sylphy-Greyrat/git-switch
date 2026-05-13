package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/config"
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
			dir, err := config.DefaultConfigDir()
			if err != nil {
				return err
			}
			templatesDir := filepath.Join(dir, "templates")
			entries, err := os.ReadDir(templatesDir)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "No templates found\n")
				return nil
			}
			for _, entry := range entries {
				if entry.IsDir() {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", entry.Name())
				}
			}
			return nil
		},
	}
}
