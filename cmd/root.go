/*
Copyright © 2026 0xPierre
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Git Saver",
	Short: "Keep in sync, save, backups repositories and organisation",
	Long:  "Keep in sync, save, backups repositories and organisation",
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
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yml", "path to config file")
}
