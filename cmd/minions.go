package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// minionsCmd represents the minions command
var minionsCmd = &cobra.Command{
	Use:                   "minions character_id",
	Short:                 "Gets collected minions from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone minions 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		minions := lodestoneWrapper.GetMinions(character_id)
		for _, minion := range minions {
			fmt.Println(minion)
		}
	},
}

func init() {
	rootCmd.AddCommand(minionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// minionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// minionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
