/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	lodestoneWrapper "github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// declutterCmd represents the declutter command
var declutterCmd = &cobra.Command{
	Use:     "declutter character_id",
	Short:   "Finds duplicate items if retainer inventory",
	Args:    cobra.ExactArgs(1),
	Example: "lodestone-cli retainer declutter 12345",
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]

		duplicateItems := false
		itemMap := make(map[string][]string)
		retainers := lodestoneWrapper.GetRetainers(character_id)
		for _, retainer := range retainers {
			for _, item := range retainer.Items {
				// don't consider a full stack, these can't be deduped
				if item.Quantity != "99" {
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