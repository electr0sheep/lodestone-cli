package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:                   "all [character_id]",
	Short:                 "Retrieves all collection data from Lodestone",
	Args:                  cobra.MaximumNArgs(1),
	Example:               "lodestone-cli all 12345",
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

		mounts := c.GetMounts()
		fmt.Println("MOUNTS")
		for _, mount := range mounts {
			fmt.Println(mount.Name)
		}

		fmt.Printf("\n\n\n")

		minions := c.GetMinions()
		fmt.Println("MINIONS")
		for _, minion := range minions {
			fmt.Println(minion.Name)
		}

		fmt.Printf("\n\n\n")

		orchestrions := c.GetOrchestrions()
		fmt.Println("ORCHESTRIONS")
		for _, orchestrion := range orchestrions {
			fmt.Println(orchestrion.Name)
		}

		fmt.Printf("\n\n\n")

		spells := c.GetSpells()
		fmt.Println("SPELLS")
		for _, spell := range spells {
			fmt.Println(spell.Name)
		}

		fmt.Printf("\n\n\n")

		achievements := c.GetAchievements()
		fmt.Println("ACHIEVEMENTS")
		for _, achievement := range achievements {
			fmt.Println(achievement.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
