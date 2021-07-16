package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// minionsCmd represents the minions command
var minionsCmd = &cobra.Command{
	Use:                   "minions character_id",
	Short:                 "Gets collected minions from Lodestone",
	Args:                  cobra.MaximumNArgs(1),
	Example:               "lodestone-cli minions 12345",
	DisableFlagsInUseLine: true,
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

		minions := c.GetMinions()
		for _, minion := range minions {
			fmt.Println(minion.Name)
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
