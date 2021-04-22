package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// spellsCmd represents the spells command
var spellsCmd = &cobra.Command{
	Use:                   "spells character_id",
	Short:                 "Gets collected blue mage spells from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone spells 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		spells := lodestoneWrapper.GetSpells(character_id)
		for _, spell := range spells {
			fmt.Println(spell)
		}
	},
}

func init() {
	rootCmd.AddCommand(spellsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spellsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spellsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
