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

// mountsCmd represents the mounts command
var mountsCmd = &cobra.Command{
	Use:                   "mounts character_id",
	Short:                 "Gets collected mounts from Lodestone",
	Args:                  cobra.ExactArgs(1),
	Example:               "lodestone mounts 12345",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		mounts := lodestoneWrapper.GetMounts(character_id)
		for _, mount := range mounts {
			fmt.Println(mount)
		}
	},
}

func init() {
	rootCmd.AddCommand(mountsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mountsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mountsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
