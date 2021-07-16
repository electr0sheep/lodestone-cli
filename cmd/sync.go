package cmd

import (
	"fmt"

	ffxivcollectWrapper "github.com/electr0sheep/lodestone-cli/ffxivcollect"
	"github.com/electr0sheep/lodestone-cli/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:                   "sync character_id",
	Short:                 "Syncs private data to ffxivcollect.com",
	Args:                  cobra.MaximumNArgs(1),
	Example:               "lodestone-cli sync 12345",
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

		fmt.Printf("Syncing blue magic...\n")
		syncBlueMagic(character_id)

		fmt.Printf("Syncing orchestrions...\n")
		syncOrchestrions(character_id)

		fmt.Printf("Syncing triple triad cards...\n")
		syncCards(character_id)
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
	c := lib.Character{Id: character_id}
	spells := c.GetSpells()
	blueMagicSpellMap := ffxivcollectWrapper.GetBlueMagicSpells()

	noSpellsAdded := true
	for _, spell := range spells {
		blueMagicSpell := blueMagicSpellMap[spell.Name]
		if !blueMagicSpell.Obtained {
			noSpellsAdded = false
			spellSucessfullyAdded := ffxivcollectWrapper.AddBlueMagicSpell(spell.Name, blueMagicSpell.Id)
			if spellSucessfullyAdded {
				fmt.Printf("Checked %s\n", spell.Name)
			} else {
				fmt.Printf("Problem checking %s\n", spell.Name)
			}
		}
	}
	if noSpellsAdded {
		fmt.Printf("All blue magic data already synced\n")
	}
}

func syncOrchestrions(character_id string) {
	c := lib.Character{Id: character_id}
	orchestrions := c.GetOrchestrions()
	orchestrionMap := ffxivcollectWrapper.GetOrchestrions()

	noOrchestrionsAdded := true
	for _, orchestrion := range orchestrions {
		orchestrionName := orchestrion.Name
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

func syncCards(character_id string) {
	c := lib.Character{Id: character_id}
	cards := c.GetCards()
	cardMap := ffxivcollectWrapper.GetCards()

	if cardMap == nil {
		return
	}

	noCardsAdded := true
	for _, card := range cards {
		cardName := card.Name
		card := cardMap[cardName]
		if !card.Obtained {
			noCardsAdded = false
			cardSucessfullyAdded := ffxivcollectWrapper.AddCard(cardName, card.Id)
			if cardSucessfullyAdded {
				fmt.Printf("Checked %s\n", cardName)
			} else {
				fmt.Printf("Problem checking %s\n", cardName)
			}
		}
	}
	if noCardsAdded {
		fmt.Printf("All triple triad data already synced\n")
	}
}
