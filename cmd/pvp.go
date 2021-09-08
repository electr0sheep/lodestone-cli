package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pvpCmd represents the pvp command
var pvpCmd = &cobra.Command{
	Use:                   "pvp",
	Short:                 "Gets pvp info",
	Args:                  cobra.MaximumNArgs(1),
	Example:               "lodestone-cli pvp 12345",
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

		c.GetPvpProfile()

		fmt.Printf("Faction   : %s\n", c.PvpProfile.Faction)
		fmt.Printf("Rank      : %s\n", c.PvpProfile.Rank)
		fmt.Printf("Total XP  : %s\n", c.PvpProfile.TotalXp)
		fmt.Printf("Current XP: %s / %s\n", c.PvpProfile.Xp, c.PvpProfile.NextXp)
		fmt.Printf("\nFrontline Overall Performance\n")
		fmt.Printf("Campaigns   : %s\n", c.PvpProfile.OverallPerformance.Campaigns)
		fmt.Printf("First Place : %s(Victory Rate: %s)\n", c.PvpProfile.OverallPerformance.FirstPlaceWins, c.PvpProfile.OverallPerformance.FirstPlaceWinPercentage)
		fmt.Printf("Second Place: %s(Victory Rate: %s)\n", c.PvpProfile.OverallPerformance.SecondPlaceWins, c.PvpProfile.OverallPerformance.SecondPlaceWinPercentage)
		fmt.Printf("Third Place : %s(Victory Rate: %s)\n", c.PvpProfile.OverallPerformance.ThirdPlaceWins, c.PvpProfile.OverallPerformance.ThirdPlaceWinPercentage)
		fmt.Printf("\nFrontline Weekly Performance\n")
		fmt.Printf("Campaigns   : %s\n", c.PvpProfile.WeeklyPerformance.Campaigns)
		fmt.Printf("First Place : %s(Victory Rate: %s)\n", c.PvpProfile.WeeklyPerformance.FirstPlaceWins, c.PvpProfile.WeeklyPerformance.FirstPlaceWinPercentage)
		fmt.Printf("Second Place: %s(Victory Rate: %s)\n", c.PvpProfile.WeeklyPerformance.SecondPlaceWins, c.PvpProfile.WeeklyPerformance.SecondPlaceWinPercentage)
		fmt.Printf("Third Place : %s(Victory Rate: %s)\n", c.PvpProfile.WeeklyPerformance.ThirdPlaceWins, c.PvpProfile.WeeklyPerformance.ThirdPlaceWinPercentage)
	},
}

func init() {
	rootCmd.AddCommand(pvpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pvpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pvpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
