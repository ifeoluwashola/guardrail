package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editConfigCmd represents the edit-config command
var editConfigCmd = &cobra.Command{
	Use:     "edit-config",
	Aliases: []string{"edit", "e", "ec"},
	Short:   "Open the config.yaml file directly in your preferred text editor",
	Long:    `Opens ~/.guardrail/config.yaml using the $EDITOR environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			color.Red("Could not find user home directory: %v", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".guardrail", "config.yaml")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			color.Red("Config file not found at %s. Run `guardrail create-config` first.", configPath)
			os.Exit(1)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano" // Fallback if $EDITOR is not defined
		}

		fmt.Printf("Opening %s in %s...\n", configPath, editor)

		editCmd := exec.Command(editor, configPath)
		// Connect the standard streams directly to shell so terminal GUI editors work correctly
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr

		err = editCmd.Run()
		if err != nil {
			color.Red("Editor exited with error: %v", err)
		} else {
			color.Green("Config file saved.")
		}
	},
}

// createConfigCmd represents the create-config command
var createConfigCmd = &cobra.Command{
	Use:     "create-config",
	Aliases: []string{"create", "c", "cc"},
	Short:   "Interactively generates a base ~/.guardrail/config.yaml configuration",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".guardrail")
		configPath := filepath.Join(configDir, "config.yaml")

		if _, err := os.Stat(configPath); err == nil {
			color.Yellow("A configuration file already exists at %s", configPath)
			return
		}

		fmt.Println(color.CyanString("Welcome to Guardrail! Let's scaffold your first environment context."))

		envName := promptText("Environment Name (e.g., dev, staging, prod)", "dev")
		cloudProvider := promptSelect("Cloud Provider", []string{"aws", "gcp"})
		cloudProfile := promptText("Cloud Profile String (e.g., aws-dev, gcp-staging-project)", "aws-dev")
		kubeContext := promptText("Kubernetes Context Name (as it appears in ~/.kube/config)", "minikube")
		isProd := promptSelect("Is this a production environment? (Will trigger safety interceptors)", []string{"false", "true"})

		isProductionBool := false
		if isProd == "true" {
			isProductionBool = true
		}

		// Ensure directory exists
		os.MkdirAll(configDir, os.ModePerm)

		// Create default dummy struct
		viper.Set("environments."+envName+".cloud_provider", cloudProvider)
		viper.Set("environments."+envName+".cloud_profile", cloudProfile)
		viper.Set("environments."+envName+".kubernetes_context", kubeContext)
		viper.Set("environments."+envName+".is_production", isProductionBool)

		err := viper.SafeWriteConfigAs(configPath)
		if err != nil {
			color.Red("Failed to create config file: %v", err)
			os.Exit(1)
		}

		color.Green("Successfully created Guardrail config at %s", configPath)
		fmt.Println("You can add more environments by running: `guardrail set-context` or `guardrail edit-config`")
	},
}

// setContextCmd represents the set-context command
var setContextCmd = &cobra.Command{
	Use:     "set-context",
	Aliases: []string{"set", "s", "sc"},
	Short:   "Interactively modify an existing environment or create a new one",
	Run: func(cmd *cobra.Command, args []string) {
		envs := viper.GetStringMap("environments")

		var envKeys []string
		for k := range envs {
			envKeys = append(envKeys, k)
		}
		envKeys = append(envKeys, "[Create New Environment]")

		selectedEnv := promptSelect("Select an environment to modify", envKeys)

		if selectedEnv == "[Create New Environment]" {
			selectedEnv = promptText("Enter new environment name", "sandbox")
		}

		envKey := "environments." + selectedEnv

		// Load Existing if available
		currProvider := viper.GetString(envKey + ".cloud_provider")
		currProfile := viper.GetString(envKey + ".cloud_profile")
		currKube := viper.GetString(envKey + ".kubernetes_context")
		currProd := viper.GetBool(envKey + ".is_production")
		if currProvider == "" {
			currProvider = "aws"
		}

		fmt.Printf("\n--- Modifying %s ---\n", selectedEnv)
		newProvider := promptSelect("Cloud Provider", []string{"aws", "gcp"})
		newProfile := promptText("Cloud Profile", currProfile)
		newKube := promptText("Kubernetes Context", currKube)
		newProdStr := promptSelect(fmt.Sprintf("Is Production? (Current: %v)", currProd), []string{"false", "true"})

		newProd := false
		if newProdStr == "true" {
			newProd = true
		}

		viper.Set(envKey+".cloud_provider", newProvider)
		viper.Set(envKey+".cloud_profile", newProfile)
		viper.Set(envKey+".kubernetes_context", newKube)
		viper.Set(envKey+".is_production", newProd)

		err := viper.WriteConfig()
		if err != nil {
			color.Red("Failed to update config file: %v", err)
		} else {
			color.Green("Successfully updated %s context in %s", selectedEnv, viper.ConfigFileUsed())
		}
	},
}

// promptText is a helper function to wrap PromptUI text inputs
func promptText(label string, defaultVal string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

// promptSelect is a helper function to wrap PromptUI select dropdowns
func promptSelect(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

func init() {
	rootCmd.AddCommand(editConfigCmd)
	rootCmd.AddCommand(createConfigCmd)
	rootCmd.AddCommand(setContextCmd)
}
