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
	lodestoneWrapper "lodestone-cli/lodestone"

	"github.com/spf13/cobra"
)

// spellsCmd represents the spells command
var spellsCmd = &cobra.Command{
	Use:                   "spells character_id",
	Short:                 "Gets collected blue mage spells from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone spells 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		spells := lodestoneWrapper.GetSpells(character_id)
		for _, spell := range spells {
			fmt.Println(spell)
		}
	},
}

func init() {
	rootCmd.AddCommand(spellsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spellsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spellsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
