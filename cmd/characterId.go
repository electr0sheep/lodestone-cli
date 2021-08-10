package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// characterIdCmd represents the characterId command
var characterIdCmd = &cobra.Command{
	Use:     "characterId",
	Short:   "Sets a default characterId so you don't have to specify it.",
	Args:    cobra.ExactArgs(1),
	Example: "lodestone-cli set characterId 12345",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("character_id", args[0])
		viper.WriteConfig()
	},
}

func init() {
	setCmd.AddCommand(characterIdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// characterIdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// characterIdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
