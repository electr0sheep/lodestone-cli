package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// currenciesCmd represents the currencies command
var currenciesCmd = &cobra.Command{
	Use:     "currencies",
	Short:   "Gets currency and reputation information from Lodestone",
	Args:    cobra.MaximumNArgs(1),
	Example: "lodestone-cli currencies 12345",
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
		c := lib.Character{Id: character_id}

		c.GetCurrenciesAndRep()
		for _, currency := range c.Currencies {
			fmt.Printf("Type: %s, Name: %s, Current Amount: %s, Maximum: %s\n", currency.Type, currency.Name, currency.CurrentAmount, currency.Maximum)
		}
		for _, reputation := range c.Reputations {
			fmt.Printf("Type: %s, Name: %s, Current Amount: %s, Maximum: %s\n", reputation.Type, reputation.Name, reputation.CurrentAmount, reputation.Maximum)
		}
	},
}

func init() {
	rootCmd.AddCommand(currenciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// currenciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// currenciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
