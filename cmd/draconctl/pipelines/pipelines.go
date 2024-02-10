package pipelines

import (
	"github.com/spf13/cobra"
)

func init() {

}

var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "A set of subcommands for building pipelines",
}

func RegisterPipelinesSubcommands(rootCmd *cobra.Command) {
	pipelinesCmd.AddCommand(buildSubCmd)

	rootCmd.AddCommand(pipelinesCmd)
}
