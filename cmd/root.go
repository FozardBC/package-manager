package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "short application for test task",
	Long: `pm is a lightweight command-line package manager written in Go. It allows you to create, distribute, and install packages via SSH.
Packages are archived as .tar.gz files and can be versioned and filtered by version constraints.
Designed for simple deployment, configuration management, or internal tooling,
pm uses JSON or YAML configuration files to define what files to include in a package or which packages to install. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}
