package cmd

import (
	"github.com/spf13/cobra"
)

// retainerCmd represents the retainer command
var retainerCmd = &cobra.Command{
	Use:                   "retainer",
	Short:                 "Gets retainer info from Lodestone",
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(retainerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// retainerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// retainerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
