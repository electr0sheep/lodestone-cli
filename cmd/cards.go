package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cardsCmd represents the cards command
var cardsCmd = &cobra.Command{
	Use:     "cards character_id",
	Short:   "Gets collected triple triad cards from Lodestone",
	Args:    cobra.MaximumNArgs(1),
	Example: "lodestone-cli cards 12345",
	Run: func(cmd *cobra.Command, args []string) {
		character_id := ""
		if len(args) == 0 {
			character_id = viper.GetString("character_id")
			if character_id == "" {
				fmt.Println("Character ID not set. Pleaes run set character_id first!")
				return
			}
		} else if len(args) == 1 {
			character_id = args[0]
		}

		c := lodestoneWrapper.Character{Id: character_id}

		cards := c.GetCards()
		fmt.Printf("You have %d cards:\n", len(cards))
		for _, card := range cards {
			fmt.Println(card.Name)
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
