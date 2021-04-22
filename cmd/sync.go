package cmd

import (
	"fmt"

	ffxivcollectWrapper "github.com/electr0sheep/lodestone-cli/ffxivcollect"
	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:                   "sync character_id",
	Short:                 "Syncs private data to ffxivcollect.com",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone sync 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]

		fmt.Printf("Syncing blue magic...\n")
		syncBlueMagic(character_id)

		fmt.Printf("Syncing orchestrions...\n")
		syncOrchestrions(character_id)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func syncBlueMagic(character_id string) {
	spells := lodestoneWrapper.GetSpells(character_id)
	blueMagicSpellMap := ffxivcollectWrapper.GetBlueMagicSpells()

	noSpellsAdded := true
	for _, spell := range spells {
		blueMagicSpell := blueMagicSpellMap[spell]
		if !blueMagicSpell.Obtained {
			noSpellsAdded = false
			spellSucessfullyAdded := ffxivcollectWrapper.AddBlueMagicSpell(spell, blueMagicSpell.Id)
			if spellSucessfullyAdded {
				fmt.Printf("Checked %s\n", spell)
			} else {
				fmt.Printf("Problem checking %s\n", spell)
			}
		}
	}
	if noSpellsAdded {
		fmt.Printf("All blue magic data already synced\n")
	}
}

func syncOrchestrions(character_id string) {
	orchestrions := lodestoneWrapper.GetOrchestrions(character_id)
	orchestrionMap := ffxivcollectWrapper.GetOrchestrions()

	noOrchestrionsAdded := true
	for _, orchestrionName := range orchestrions {
		orchestrion := orchestrionMap[orchestrionName]
		if !orchestrion.Obtained {
			noOrchestrionsAdded = false
			orchestrionSucessfullyAdded := ffxivcollectWrapper.AddOrchestrion(orchestrionName, orchestrion.Id)
			if orchestrionSucessfullyAdded {
				fmt.Printf("Checked %s\n", orchestrionName)
			} else {
				fmt.Printf("Problem checking %s\n", orchestrionName)
			}
		}
	}
	if noOrchestrionsAdded {
		fmt.Printf("All orchestrion data already synced\n")
	}
}
