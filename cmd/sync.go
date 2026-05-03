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

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync repositories from source to targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.Load(configPath)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", config)

		err = worker.Sync(config)

		return err
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
