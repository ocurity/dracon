package pipelines

import (
	"github.com/spf13/cobra"
)

var pipelinesCmd = &cobra.Command{
	Use:     "pipelines",
	Short:   "A set of subcommands for building pipelines",
	GroupID: "top-level",
}

func RegisterPipelinesSubcommands(rootCmd *cobra.Command) {
	pipelinesCmd.AddGroup(&cobra.Group{
		ID:    "Pipelines",
		Title: "Pipeline Commands:",
	})
	pipelinesCmd.AddCommand(newBuildSubCmd())

	rootCmd.AddCommand(pipelinesCmd)
}
