package cmd

import (
	"github.com/electr0sheep/lodestone-cli/lodestone"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Lodestone",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// usernamePrompt := promptui.Prompt{
		// 	Label: "Lodestone Username",
		// }
		// username, err := usernamePrompt.Run()
		// if err != nil {
		// 	panic(err)
		// }

		// passwordPrompt := promptui.Prompt{
		// 	Label: "Lodestone Password",
		// 	Mask:  '*',
		// }
		// password, err := passwordPrompt.Run()
		// if err != nil {
		// 	panic(err)
		// }
		lodestone.Login()

		// select example, can use something similar to pick character
		// prompt := promptui.Select{
		// 	Label: "Select Day",
		// 	Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		// 		"Saturday", "Sunday"},
		// 	Size: 7,
		// }

		// _, result, err := prompt.Run()

		// if err != nil {
		// 	fmt.Printf("Prompt failed %v\n", err)
		// 	return
		// }

		// fmt.Printf("You choose %q\n", result)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
