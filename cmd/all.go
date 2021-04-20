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

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:                   "all character_id",
	Short:                 "Retrieves all collection data from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone all 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]

		mounts := lodestoneWrapper.GetMounts(character_id)
		fmt.Println("MOUNTS")
		for _, mount := range mounts {
			fmt.Println(mount)
		}

		fmt.Printf("\n\n\n")

		minions := lodestoneWrapper.GetMinions(character_id)
		fmt.Println("MINIONS")
		for _, minion := range minions {
			fmt.Println(minion)
		}

		fmt.Printf("\n\n\n")

		orchestrions := lodestoneWrapper.GetOrchestrions(character_id)
		fmt.Println("ORCHESTRIONS")
		for _, orchestrion := range orchestrions {
			fmt.Println(orchestrion)
		}

		fmt.Printf("\n\n\n")

		spells := lodestoneWrapper.GetSpells(character_id)
		fmt.Println("SPELLS")
		for _, spell := range spells {
			fmt.Println(spell)
		}

		fmt.Printf("\n\n\n")

		achievements := lodestoneWrapper.GetAchievements(character_id)
		fmt.Println("ACHIEVEMENTS")
		for _, achievement := range achievements {
			fmt.Println(achievement)
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
