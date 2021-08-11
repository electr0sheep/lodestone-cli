package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// declutterCmd represents the declutter command
var declutterCmd = &cobra.Command{
	Use:     "declutter character_id",
	Short:   "Finds duplicate items if retainer inventory",
	Args:    cobra.MaximumNArgs(1),
	Example: "lodestone-cli retainer declutter 12345",
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

		duplicateItems := false
		itemMap := make(map[string][]string)
		c.GetRetainers()
		for _, retainer := range c.Retainers {
			for _, item := range retainer.Items {
				if item.IsStackable() {
					var name string
					if item.HighQuality {
						name = fmt.Sprintf("%s HQ", item.Name)
					} else {
						name = item.Name
					}
					itemMap[name] = append(itemMap[name], retainer.Name)
				}
			}
		}
		fmt.Println("==================================================")
		fmt.Println("            UNSTACKED RETAINER ITEMS")
		fmt.Println("==================================================")
		for itemName, retainerNames := range itemMap {
			if len(retainerNames) > 1 {
				duplicateItems = true
				fmt.Printf("%s was found in the following retainer inventories:\n", itemName)
				for _, retainerName := range retainerNames {
					fmt.Printf("%s\n", retainerName)
				}
				fmt.Printf("\n")
			}
		}
		if !duplicateItems {
			fmt.Printf("No duplicate items found in retainer inventories!")
		}
		fmt.Println("")
		fmt.Println("==================================================")
		fmt.Println("              PURCHASABLE ITEMS")
		fmt.Println("==================================================")
		for _, retainer := range c.Retainers {
			for _, item := range retainer.Items {
				if item.Purchasable {
					var name string
					if item.HighQuality {
						name = fmt.Sprintf("%s HQ", item.Name)
					} else {
						name = item.Name
					}
					fmt.Printf("%s was found in the inventory of %s and is purchasable with gil\n", name, retainer.Name)
				}
			}
		}
	},
}

func init() {
	retainerCmd.AddCommand(declutterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// declutterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// declutterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
