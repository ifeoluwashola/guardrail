package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "guardrail",
	Short: "Guardrail - Safe environment switching CLI",
	Long: `Guardrail is a context-aware CLI tool designed to safely manage 
environment switching across various clusters and cloud accounts. It 
prevents accidental destructive commands in production environments.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// In Go, uppercase functions (like Execute) are public (exported). Lowercase are private.
func Execute() {
	// Error handling in Go uses explicit checks rather than try/except blocks like in Python.
	err := rootCmd.Execute()
	if err != nil {
		// os.Exit(1) is similar to sys.exit(1) in Python
		os.Exit(1)
	}
}

func init() {
	// init() is a special function in Go that runs before main()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// root.PersistentFlags() returns a pointer to the FlagSet object,
	// similar to getting a reference to an object in Python.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.guardrail/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		// Another example of standard Go error handling: if the error is not nil (None in Python), handle it.
		cobra.CheckErr(err)

		// Search config in ~/.guardrail directory with name "config" (without extension).
		viper.AddConfigPath(home + "/.guardrail")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
