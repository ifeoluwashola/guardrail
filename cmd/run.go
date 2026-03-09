package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [flags] -- [command...]",
	Short: "Execute a command under the active environment's safety context",
	Long: `Executes arbitrary shell commands within the current environment context.
It intercepts commands targeting production and checks them for high-risk operations 
(e.g., delete, uninstall, destroy). If a risk is detected, it will block execution 
until explicit terminal confirmation is provided.

Example:
  gr use prod
  gr run -- kubectl delete pods --all`,
	// Cobra's standard behavior is to parse flags across the whole string.
	// DisableFlagParsing ensures anything after 'run' (or '--' depending on config)
	// is treated purely as positional arguments, exactly as the user typed them.
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// If '--' is the first argument, shift it off
		if len(args) > 0 && args[0] == "--" {
			args = args[1:]
		}

		if len(args) == 0 {
			color.Red("No command provided to execute. Try: gr run -- <command>")
			os.Exit(1)
		}

		executeInterceptor(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// In Python, you use *args in function signatures to capture arbitrary arguments into a list.
// In Go, variadic functions use `...type`. Below, `args []string` is just a standard slice argument,
// but we will use the `...` operator later to UPpack the slice into another variadic function.
func executeInterceptor(args []string) {
	// Reconstruct the full command string for checking keywords
	fullCmd := strings.Join(args, " ")

	activeEnv := viper.GetString("active_environment")
	if activeEnv == "" {
		color.Red("No active environment set. Run `gr use <env>` first.")
		os.Exit(1)
	}

	envKey := fmt.Sprintf("environments.%s", activeEnv)
	isProduction := viper.GetBool(envKey + ".is_production")
	kubeContext := viper.GetString(envKey + ".kubernetes_context")

	if isProduction {
		if isHighRisk(fullCmd) {
			color.Yellow("Safety Engine: High-risk command detected in production context!")
			c := color.New(color.FgHiWhite, color.BgRed, color.Bold)
			c.Println(" DANGER: You are targeting PRODUCTION. ")

			prompt := promptui.Prompt{
				Label: fmt.Sprintf("Type the exact cluster context name ('%s') to confirm", kubeContext),
				Validate: func(input string) error {
					if input != kubeContext {
						return errors.New("input does not match the cluster context")
					}
					return nil
				},
			}

			result, err := prompt.Run()

			if err != nil {
				color.Red("\nPrompt failed or cancelled: %v. Aborting execution.", err)
				os.Exit(1)
			}

			color.Green("Confirmation accepted ('%s'). Proceeding...", result)
		} else {
			color.Blue("Safety Engine: Command does not contain high-risk keywords. Passing through.")
		}
	}

	// EXECUTION PHASE
	// args[0] is the command (e.g., 'kubectl')
	// args[1:] are the arguments (e.g., ['get', 'pods'])
	executable := args[0]
	cmdArgs := []string{}
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	// This is the Go variadic unpacking in action:
	// We pass the slice `cmdArgs` followed by `...` to unpack it as individual arguments into exec.Command,
	// which has the signature `func Command(name string, arg ...string) *Cmd`.
	// This is analogous to doing `subprocess.run([executable, *cmdArgs])` in Python.
	execCmd := exec.Command(executable, cmdArgs...)

	// Bind standard standard streams so output streams directly to the terminal
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	err := execCmd.Run()
	if err != nil {
		// Exit with the same status code if possible, or 1
		color.Red("Command exited with error: %v", err)
		os.Exit(1)
	}
}

func isHighRisk(cmdStr string) bool {
	// Standard substring tests
	keywords := []string{"delete", "uninstall", "destroy", "apply"}
	for _, kw := range keywords {
		if strings.Contains(cmdStr, kw) {
			return true
		}
	}
	return false
}
