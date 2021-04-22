package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// achievementsCmd represents the achievements command
var achievementsCmd = &cobra.Command{
	Use:                   "achievements character_id",
	Short:                 "Gets collected achievements from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone achievements 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		achievements := lodestoneWrapper.GetAchievements(character_id)
		for _, achievement := range achievements {
			fmt.Println(achievement)
		}
	},
}

func init() {
	rootCmd.AddCommand(achievementsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// achievementsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// achievementsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
