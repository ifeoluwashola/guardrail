package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listContextsCmd represents the list-contexts command
var listContextsCmd = &cobra.Command{
	Use:     "list-contexts",
	Aliases: []string{"lc"},
	Short:   "List all available Guardrail environments",
	Long: `Reads your ~/.guardrail/config.yaml and displays all configured environments.
The currently active environment is highlighted.`,
	Run: func(cmd *cobra.Command, args []string) {
		environments := viper.GetStringMap("environments")
		if len(environments) == 0 {
			color.Yellow("No environments found. Run 'guardrail create-config' to set them up.")
			return
		}

		activeEnv := viper.GetString("active_environment")

		fmt.Println("Available Environments:")
		fmt.Println(strings.Repeat("-", 40))

		for envName, envDataRaw := range environments {
			envData, ok := envDataRaw.(map[string]interface{})
			if !ok {
				continue
			}

			// Extract data safely
			isProd := false
			if val, exists := envData["is_production"]; exists {
				isProd, _ = val.(bool)
			}
			kubeContext := ""
			if val, exists := envData["kubernetes_context"]; exists {
				kubeContext, _ = val.(string)
			}

			// Format the output string
			displayStr := fmt.Sprintf("  %s (K8s: %s)", envName, kubeContext)
			if isProd {
				displayStr += color.RedString(" [PROD]")
			}

			// Highlight if active
			if envName == activeEnv {
				color.Green("→ %s *Active*", displayStr[2:]) // Strip the leading spaces, prepend arrow
			} else {
				fmt.Println(displayStr)
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	},
}

func init() {
	rootCmd.AddCommand(listContextsCmd)
}
