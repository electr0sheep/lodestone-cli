package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// orchestrionsCmd represents the orchestrions command
var orchestrionsCmd = &cobra.Command{
	Use:                   "orchestrions character_id",
	Short:                 "Gets collected orchestrions from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone-cli orchestrions 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		orchestrions := lodestoneWrapper.GetOrchestrions(character_id)
		for _, orchestrion := range orchestrions {
			fmt.Println(orchestrion)
		}
	},
}

func init() {
	rootCmd.AddCommand(orchestrionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orchestrionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orchestrionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
