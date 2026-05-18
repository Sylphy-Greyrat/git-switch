package commands

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sylphy/git-switch/core/hook"
)

func newHookCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "hook", Short: "Manage git alias and shell hooks"}
	cmd.AddCommand(hookInstallCommand())
	cmd.AddCommand(hookUninstallCommand())
	cmd.AddCommand(hookStatusCommand())
	return cmd
}

func hookInstallCommand() *cobra.Command {
	var shell string
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install git alias 'git sw'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hook.InstallGitAlias(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed git alias: git sw -> git-switch\n")

			effectiveShell := shell
			if effectiveShell == "" {
				detectedShell, err := hook.DetectCurrentShell()
				if err != nil {
					return err
				}
				effectiveShell = detectedShell
			}

			switch effectiveShell {
			case "powershell", "pwsh":
				if err := hook.InstallPowerShellHook(); err != nil {
					return fmt.Errorf("powershell hook: %w", err)
				}
			default:
				if err := hook.InstallShellHook(effectiveShell); err != nil {
					return fmt.Errorf("shell hook: %w", err)
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed shell hook for %s\n", effectiveShell)

			// Install completion
			completionScript, err := genCompletionScript(effectiveShell, cmd.Root())
			if err != nil {
				return fmt.Errorf("generate completion: %w", err)
			}
			if err := hook.WriteCompletionScript(effectiveShell, completionScript); err != nil {
				return fmt.Errorf("write completion: %w", err)
			}
			if err := hook.InjectCompletionIntoRC(effectiveShell); err != nil {
				return fmt.Errorf("inject completion: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed completion for %s\n", effectiveShell)
			fmt.Fprintf(cmd.OutOrStdout(), "Run 'source %s' or open a new terminal for completion to take effect\n",
				hook.RCPathHint(effectiveShell))

			return nil
		},
	}
	cmd.Flags().StringVar(&shell, "shell", "", "Install shell hook (bash, zsh, powershell, pwsh)")
	return cmd
}

func hookUninstallCommand() *cobra.Command {
	var shell string
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove git alias and shell hooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hook.UninstallGitAlias(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed git alias: git sw\n")

			if shell == "" {
				detectedShell, err := hook.DetectCurrentShell()
				if err != nil {
					return err
				}
				shell = detectedShell
			}

			switch shell {
			case "powershell", "pwsh":
				if err := hook.UninstallPowerShellHook(); err != nil {
					return fmt.Errorf("powershell hook: %w", err)
				}
			default:
				if err := hook.UninstallShellHook(shell); err != nil {
					return fmt.Errorf("shell hook: %w", err)
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed shell hook for %s\n", shell)

			// Remove completion
			if err := hook.RemoveCompletionFromRC(shell); err != nil {
				return fmt.Errorf("remove completion: %w", err)
			}
			if err := hook.RemoveCompletionScript(shell); err != nil {
				return fmt.Errorf("remove completion: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed completion for %s\n", shell)

			return nil
		},
	}
	cmd.Flags().StringVar(&shell, "shell", "", "Remove shell hook (bash, zsh, powershell, pwsh)")
	return cmd
}

func hookStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show hook installation status",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			installed, val, err := hook.StatusGitAlias()
			if err != nil {
				return err
			}
			if installed {
				fmt.Fprintf(out, "git alias 'sw': installed (%s)\n", val)
			} else {
				fmt.Fprintf(out, "git alias 'sw': not installed\n")
			}

			for _, sh := range []string{"bash", "zsh"} {
				installed, err := hook.IsShellHookInstalled(sh)
				if err != nil {
					continue
				}
				if installed {
					fmt.Fprintf(out, "shell hook (%s): installed\n", sh)
				} else {
					fmt.Fprintf(out, "shell hook (%s): not installed\n", sh)
				}
			}

			installed, err = hook.IsPowerShellHookInstalled()
			if err != nil {
				fmt.Fprintf(out, "powershell hook: error checking: %v\n", err)
			} else if installed {
				fmt.Fprintf(out, "powershell hook: installed\n")
			} else {
				fmt.Fprintf(out, "powershell hook: not installed\n")
			}

			for _, sh := range []string{"bash", "zsh"} {
				installed, err := hook.IsCompletionInstalledForShell(sh)
				if err != nil {
					continue
				}
				if installed {
					fmt.Fprintf(out, "completion (%s): installed\n", sh)
				} else {
					fmt.Fprintf(out, "completion (%s): not installed\n", sh)
				}
			}

			installed, err = hook.IsCompletionInstalledForShell("pwsh")
			if err != nil {
				fmt.Fprintf(out, "completion (pwsh): error checking: %v\n", err)
			} else if installed {
				fmt.Fprintf(out, "completion (pwsh): installed\n")
			} else {
				fmt.Fprintf(out, "completion (pwsh): not installed\n")
			}

			return nil
		},
	}
}

func genCompletionScript(shell string, root *cobra.Command) (string, error) {
	var buf bytes.Buffer
	switch shell {
	case "bash":
		if err := root.GenBashCompletionV2(&buf, true); err != nil {
			return "", err
		}
	case "zsh":
		if err := root.GenZshCompletion(&buf); err != nil {
			return "", err
		}
	case "powershell", "pwsh":
		if err := root.GenPowerShellCompletion(&buf); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
	return buf.String(), nil
}
