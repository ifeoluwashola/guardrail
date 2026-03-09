package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prints the current Guardrail context for PS1 integration",
	Long: `Returns a formatted string representing the current environment context, 
designed to be embedded into your shell's PS1 or PROMPT variable.

Example usage in ~/.bashrc or ~/.zshrc:
  export PS1="\$(gr prompt) \u@\h:\w$ "`,
	Run: func(cmd *cobra.Command, args []string) {
		activeEnv := viper.GetString("active_environment")
		if activeEnv == "" {
			// If no environment is active, print nothing so the prompt isn't cluttered
			return
		}

		envKey := fmt.Sprintf("environments.%s", activeEnv)
		isProd := viper.GetBool(envKey + ".is_production")
		kubeContext := viper.GetString(envKey + ".kubernetes_context")

		// Shorten kubeContext if it's an ARN (e.g. arn:...:cluster/dev-cluster -> dev-cluster)
		parts := strings.Split(kubeContext, "/")
		shortContext := parts[len(parts)-1]
		if shortContext == "" {
			shortContext = "unknown-cluster"
		}

		// When a command is run inside a subshell like $(gr prompt), stdout is technically a pipe,
		// and most color libraries (including fatih/color) will automatically disable colors.
		// We explicitly force color enablement here.
		color.NoColor = false

		envDisplay := strings.ToUpper(activeEnv)
		promptStr := fmt.Sprintf("[%s | %s] ", envDisplay, shortContext)

		if isProd {
			fmt.Print(color.RedString(promptStr))
		} else {
			fmt.Print(color.CyanString(promptStr))
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
