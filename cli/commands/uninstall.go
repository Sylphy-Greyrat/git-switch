package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/config"
)

func newUninstallCommand() *cobra.Command {
	var keepConfig bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall git-switch and clean up",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Remove git alias
			out, err := exec.Command("git", "config", "--global", "--unset", "alias.sw").CombinedOutput()
			if err != nil {
				// Alias may not exist — report but don't fail
				fmt.Fprintf(cmd.ErrOrStderr(), "Note: %s\n", strings.TrimSpace(string(out)))
			}

			if !keepConfig {
				dir, err := config.DefaultConfigDir()
				if err != nil {
					return err
				}
				if err := os.RemoveAll(dir); err != nil {
					return fmt.Errorf("remove config dir: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Removed config directory: %s\n", dir)
			}

			// Remove generated SSH config file
			home, err := os.UserHomeDir()
			if err == nil {
				sshConfigFile := home + "/.ssh/config.d/git-switch"
				if err := os.Remove(sshConfigFile); err == nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Removed SSH config file: %s\n", sshConfigFile)
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Uninstalled git-switch\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&keepConfig, "keep-config", false, "Keep configuration files")
	return cmd
}
