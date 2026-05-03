/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/worker"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull repositories from source",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.Load(configPath)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", config)

		_, err = worker.Pull(config)

		return err
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
