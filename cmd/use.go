package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [environment]",
	Short: "Switch to a specified environment",
	Long: `Switches your local developer environment (Kubernetes context, cloud profile, etc.) 
to the one specified in your configuration.

Example:
  gr use dev
  gr use prod`,
	Args: cobra.ExactArgs(1), // Ensure exactly one argument is passed
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]

		runUseConfig(env)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUseConfig(env string) {
	fmt.Printf("Attempting to switch to environment: %s\n", env)

	// Viper allows us to query configurations using dot notation dict-style
	envKey := fmt.Sprintf("environments.%s", env)

	if !viper.IsSet(envKey) {
		// In Go, red.Printf is similar to using formatting strings with colored output libraries in Python
		color.Red("Error: Environment '%s' not found in configuration.", env)
		os.Exit(1)
	}

	// Fetch specific fields from the environment block
	// This replaces dict.get("key") from Python
	kubeContext := viper.GetString(envKey + ".kubernetes_context")
	cloudProfile := viper.GetString(envKey + ".cloud_profile")
	cloudProvider := viper.GetString(envKey + ".cloud_provider")
	isProduction := viper.GetBool(envKey + ".is_production")

	// 1. Switch Kubernetes Context (with Validation)
	if kubeContext != "" {
		fmt.Printf("Validating Kubernetes context: %s\n", kubeContext)

		// In Go, client-go acts like the kubernetes python client. We load the raw config file.
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.LoadFromFile(kubeconfigPath)
		if err != nil {
			color.Yellow("Warning: Could not load ~/.kube/config to validate context: %v", err)
		} else {
			// Check if the mapped context actually exists in the file
			if _, exists := config.Contexts[kubeContext]; !exists {
				color.Red("Error: Kubernetes context '%s' does not exist in your ~/.kube/config!", kubeContext)
				color.Red("Aborting environment switch to prevent corrupted shell state.")
				os.Exit(1)
			}
		}

		fmt.Printf("Switching Kubernetes context to: %s\n", kubeContext)
		// os/exec package is Go's equivalent to Python's subprocess.
		// exec.Command returns an *exec.Cmd struct, it does not run it immediately.
		// It's similar to preparing a list like ["kubectl", "config", "use-context", kubeContext]
		kubectlCmd := exec.Command("kubectl", "config", "use-context", kubeContext)

		// In Go, you can redirect stdout/stderr by assigning fields on the Cmd struct.
		// Here, CombinedOutput() runs the command and returns its combined standard output and error.
		// This is analogous to subprocess.run(..., capture_output=True) in Python.
		out, err := kubectlCmd.CombinedOutput()
		if err != nil {
			color.Red("Failed to switch K8s context: %v\nOutput: %s", err, string(out))
		} else {
			color.Green("✓ K8s context switched successfully.")
		}
	} else {
		color.Yellow("No Kubernetes context defined for this environment.")
	}

	// 2. Switch Cloud Profile
	if cloudProfile != "" {
		fmt.Printf("Switching Cloud profile to: %s [%s]\n", cloudProfile, cloudProvider)

		var cloudCmd *exec.Cmd

		switch strings.ToLower(cloudProvider) {
		case "aws":
			cloudCmd = exec.Command("aws", "configure", "set", "profile", cloudProfile)
		case "gcp":
			cloudCmd = exec.Command("gcloud", "config", "set", "project", cloudProfile)
		default:
			// If missing or unhandled, default to generic aws mock for backward compatibility
			cloudCmd = exec.Command("aws", "configure", "set", "profile", cloudProfile)
		}

		cloudCmd.Stdout = os.Stdout
		cloudCmd.Stderr = os.Stderr

		err := cloudCmd.Run()
		if err != nil {
			color.Yellow("Note: Cloud command failed, likely because CLI tool is not installed or profile does not exist locally. Error: %v", err)
		} else {
			color.Green("✓ Cloud profile switched successfully.")
		}
	} else {
		color.Yellow("No Cloud profile defined for this environment.")
	}

	fmt.Println(strings.Repeat("-", 40))
	if isProduction {
		// Using fatih/color to print bright red for production
		c := color.New(color.FgHiWhite, color.BgRed, color.Bold)
		c.Println(" WARNING: NOW OPERATING IN PRODUCTION ")
	} else {
		color.Green("Environment is safe for standard operations.")
	}

	// Persist the active environment to the configuration file so `gr run` has context.
	viper.Set("active_environment", env)
	err := viper.WriteConfig()
	if err != nil {
		color.Red("Warning: Failed to persist active environment to %s: %v", viper.ConfigFileUsed(), err)
	} else {
		color.Green("Active environment state saved.")
	}
}
