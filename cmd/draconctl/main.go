package main

import (
	"os"

	"github.com/ocurity/dracon/cmd/draconctl/migrations"
	"github.com/ocurity/dracon/cmd/draconctl/pipelines"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "draconctl",
	Short: "A CLI to manage all things related to Dracon",
}

func main() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "top-level",
		Title: "Top-level Commands:",
	})

	pipelines.RegisterPipelinesSubcommands(rootCmd)
	migrations.RegisterMigrationsSubcommands(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
