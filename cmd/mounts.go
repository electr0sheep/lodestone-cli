package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mountsCmd represents the mounts command
var mountsCmd = &cobra.Command{
	Use:                   "mounts character_id",
	Short:                 "Gets collected mounts from Lodestone",
	Args:                  cobra.MaximumNArgs(1),
	Example:               "lodestone-cli mounts 12345",
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

		c.GetMounts()
		c.GetMountDetails(&c.Mounts[1])
		for _, mount := range c.Mounts {
			fmt.Println(mount.Name)
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
