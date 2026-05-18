package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validShells = []string{"bash", "zsh", "powershell", "pwsh"}

func newCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "completion <shell>",
		Short:     "Generate shell completion script",
		Long:      "Generate shell completion script. Output to stdout, ready to source in your RC file.\n\nSupported shells: bash, zsh, powershell, pwsh",
		Args:      cobra.ExactArgs(1),
		ValidArgs: validShells,
		RunE: func(cmd *cobra.Command, args []string) error {
			script, err := genCompletionScript(args[0], cmd.Root())
			if err != nil {
				return err
			}
			if _, err := fmt.Fprint(os.Stdout, script); err != nil {
				return fmt.Errorf("write completion: %w", err)
			}
			return nil
		},
	}
}
