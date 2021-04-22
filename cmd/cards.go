package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"
	"github.com/spf13/cobra"
)

// cardsCmd represents the cards command
var cardsCmd = &cobra.Command{
	Use:   "cards character_id",
	Short: "Gets collected triple triad cards from Lodestone",
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		achievements := lodestoneWrapper.GetCards(character_id)
		for _, achievement := range achievements {
			fmt.Println(achievement)
		}
	},
}

func init() {
	rootCmd.AddCommand(cardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
