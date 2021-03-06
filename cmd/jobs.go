package cmd

import (
	"fmt"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:     "jobs",
	Short:   "Gets name and level of jobs from Lodestone",
	Args:    cobra.MaximumNArgs(1),
	Example: "lodestone-cli jobs 12345",
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

		c.GetJobs()
		for _, job := range c.Jobs {
			fmt.Printf("Role: %s, Job: %s, Level: %s, Xp: %s\n", job.Role, job.Name, job.Level, job.Xp)
		}
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
