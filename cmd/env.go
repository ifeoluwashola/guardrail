package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Prints shell export statements for the active environment",
	Long: `Prints export statements for cloud profiles and kubernetes contexts, 
designed to be evaluated by your parent shell.

Example usage:
  eval $(gr env)

You can create an alias in your ~/.bashrc or ~/.zshrc:
  alias gr-use='gr use $1 && eval $(gr env)'`,
	Run: func(cmd *cobra.Command, args []string) {
		activeEnv := viper.GetString("active_environment")
		if activeEnv == "" {
			// Do not print anything to standard out that isn't valid bash to avoid breaking eval
			fmt.Fprintln(os.Stderr, "No active environment set. Run `gr use <env>` first.")
			os.Exit(1)
		}

		envKey := fmt.Sprintf("environments.%s", activeEnv)
		kubeContext := viper.GetString(envKey + ".kubernetes_context")
		cloudProfile := viper.GetString(envKey + ".cloud_profile")
		cloudProvider := viper.GetString(envKey + ".cloud_provider")

		if kubeContext != "" {
			fmt.Printf("export KUBECONFIG_CONTEXT=\"%s\"\n", kubeContext)
			// In some setups, you might literally swap KUBECONFIG here instead
		}

		if cloudProfile != "" {
			switch strings.ToLower(cloudProvider) {
			case "gcp":
				fmt.Printf("export CLOUDSDK_CORE_PROJECT=\"%s\"\n", cloudProfile)
				fmt.Printf("export GOOGLE_CLOUD_PROJECT=\"%s\"\n", cloudProfile)
			default:
				// AWS is standard fallback
				fmt.Printf("export AWS_PROFILE=\"%s\"\n", cloudProfile)
				fmt.Printf("export AWS_DEFAULT_PROFILE=\"%s\"\n", cloudProfile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
