package cmd

import (
	"fmt"
	"strings"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type command struct {
	Label    string
	Helptext string
}

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := lib.Character{Id: viper.GetString("character_id")}
		c.GetProfile()

		mainMenu(c)
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// interactiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// interactiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mainMenu(c lib.Character) {
	commands := []command{
		{Label: "Character", Helptext: ""},
		{Label: "Companions", Helptext: ""},
		{Label: "Retainers", Helptext: ""},
		{Label: "Exit", Helptext: ""},
	}

	templates := &promptui.SelectTemplates{
		Label:    c.Name,
		Active:   "\U000025B8 {{ .Label | cyan }}",
		Inactive: "  {{ .Label | cyan }}",
		Selected: "\U000025B8 {{ .Label | red | cyan }}",
		Details:  "{{ .Helptext }}",
	}

	prompt := promptui.Select{
		Items:        commands,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		characterMenu(c)
	case 1:
		companionMenu(c)
	case 2:
		retainerMenu(c)
	case 3:
		return
	}
}

func characterMenu(c lib.Character) {
	commands := []command{
		{Label: "Profile", Helptext: ""},
		{Label: "Class/Job", Helptext: ""},
		{Label: "Minions", Helptext: ""},
		{Label: "Mounts", Helptext: ""},
		{Label: "Currencies/Reputation", Helptext: ""},
		{Label: "Quests", Helptext: ""},
		{Label: "Achievements", Helptext: ""},
		{Label: "Orchestrion Roll", Helptext: ""},
		{Label: "PvP Profile", Helptext: ""},
		{Label: "Blue Magic Spellbook", Helptext: ""},
		{Label: "Trust", Helptext: ""},
		{Label: "The Gold Saucer", Helptext: ""},
		{Label: "Triple Triad", Helptext: ""},
		{Label: "Return", Helptext: ""},
		{Label: "Exit", Helptext: ""},
	}

	templates := &promptui.SelectTemplates{
		Label:    c.Name,
		Active:   "\U000025B8 {{ .Label | cyan }}",
		Inactive: "  {{ .Label | cyan }}",
		Selected: "\U000025B8 {{ .Label | red | cyan }}",
		Details:  "{{ .Helptext }}",
	}

	prompt := promptui.Select{
		Items:        commands,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		profileMenu(c)
	case 1:
		jobMenu(c)
	case 2:
		minionMenu(c)
	case 3:
		mountMenu(c)
	case 4:
		currencyMenu(c)
	case 5:
		questMenu(c)
	case 6:
		achievementMenu(c)
	case 7:
		orchestrionMenu(c)
	case 8:
		pvpMenu(c)
	case 9:
		blueMagicMenu(c)
	case 10:
		trustMenu(c)
	case 11:
		goldSaucerMenu(c)
	case 12:
		tripleTriadMenu(c)
	case 13:
		characterMenu(c)
	case 14:
		return
	}
}

func companionMenu(c lib.Character) {

}

func retainerMenu(c lib.Character) {

}

func profileMenu(c lib.Character) {
	type attribute struct {
		Name  string
		Value string
	}

	characterAttributes := []attribute{
		{Name: "Title        ", Value: c.Title},
		{Name: "World        ", Value: c.World},
		{Name: "Race         ", Value: c.Race},
		{Name: "Clan         ", Value: c.Clan},
		{Name: "Gender       ", Value: c.Gender},
		{Name: "Nameday      ", Value: c.Nameday},
		{Name: "Guardian     ", Value: c.Guardian},
		{Name: "City-state   ", Value: c.CityState},
		{Name: "Grand Company", Value: c.GrandCompany},
		{Name: "Free Company ", Value: c.FreeCompany},
	}

	for _, linkshell := range c.Linkshells {
		characterAttributes = append(characterAttributes, attribute{Name: "Linkshell    ", Value: linkshell})
	}

	templates := &promptui.SelectTemplates{
		Label:    c.Name,
		Active:   "\U000025B8 {{ .Name | cyan }}: {{ .Value }}",
		Inactive: "  {{ .Name | cyan }}: {{ .Value }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}: {{ .Value }}",
	}

	prompt := promptui.Select{
		Items:        characterAttributes,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	characterMenu(c)
}

func jobMenu(c lib.Character) {
	if c.Jobs == nil {
		c.GetJobs()
	}

	templates := &promptui.SelectTemplates{
		Label:    "Class/Job",
		Active:   "\U000025B8 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}",
		Details: `
Level: {{ .Level }}
Role: {{ .Role }}
Xp: {{ .Xp }}`,
	}

	prompt := promptui.Select{
		Items:        c.Jobs,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	characterMenu(c)
}

func minionMenu(c lib.Character) {
	if c.Minions == nil {
		c.GetMinions()
	}

	templates := &promptui.SelectTemplates{
		Label:    "Minions",
		Active:   "\U000025B8 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}",
	}

	prompt := promptui.Select{
		Items:        c.Minions,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if c.Minions[i].Description == "" {
		c.GetMinionDetails(&c.Minions[i])
		// also need to format description for use here
		c.Minions[i].Description = formatForTerminal(c.Minions[i].Description)
	}

	minionDetailMenu(c, c.Minions[i])
}

func minionDetailMenu(c lib.Character, m lib.Minion) {
	type attribute struct {
		Name  string
		Value string
	}

	minionAttributes := []attribute{
		{Name: "Behavior        ", Value: m.Behavior},
		{Name: "Acquisition Date", Value: m.AcquistionDate},
		{Name: "Description     ", Value: ""},
	}

	templates := &promptui.SelectTemplates{
		Label:    m.Name,
		Active:   "\U000025B8 {{ .Name | cyan }}: {{ .Value }}",
		Inactive: "  {{ .Name | cyan }}: {{ .Value }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}: {{ .Value }}",
		Details:  m.Description,
	}

	prompt := promptui.Select{
		Items:        minionAttributes,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	characterMenu(c)
}

func mountMenu(c lib.Character) {
	if c.Mounts == nil {
		c.GetMounts()
	}

	templates := &promptui.SelectTemplates{
		Label:    "Mounts",
		Active:   "\U000025B8 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}",
	}

	prompt := promptui.Select{
		Items:        c.Mounts,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if c.Mounts[i].Description == "" {
		c.GetMountDetails(&c.Mounts[i])
		// also need to format description for use here
		c.Mounts[i].Description = formatForTerminal(c.Mounts[i].Description)
	}

	mountDetailMenu(c, c.Mounts[i])
}

func mountDetailMenu(c lib.Character, m lib.Mount) {
	type attribute struct {
		Name  string
		Value string
	}

	mountAttributes := []attribute{
		{Name: "Movement        ", Value: m.Movement},
		{Name: "Acquisition Date", Value: m.AcquistionDate},
		{Name: "Description     ", Value: ""},
	}

	templates := &promptui.SelectTemplates{
		Label:    m.Name,
		Active:   "\U000025B8 {{ .Name | cyan }}: {{ .Value }}",
		Inactive: "  {{ .Name | cyan }}: {{ .Value }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}: {{ .Value }}",
		Details:  m.Description,
	}

	prompt := promptui.Select{
		Items:        mountAttributes,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	characterMenu(c)
}

func currencyMenu(c lib.Character) {

}

func questMenu(c lib.Character) {

}

func achievementMenu(c lib.Character) {

}

func orchestrionMenu(c lib.Character) {

}

func pvpMenu(c lib.Character) {

}

func blueMagicMenu(c lib.Character) {
	if c.Spells == nil {
		c.GetSpells()
		// also need to format description for use here
		for i, spell := range c.Spells {
			c.Spells[i].Description = formatForTerminal(spell.Description)
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "Blue Magic",
		Active:   "\U000025B8 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000025B8 {{ .Name | red | cyan }}",
		Details: `
{{ .Description }}
Type: {{ .Type }}
Aspect: {{ .Aspect }}
Rank: {{ .Rank }}
How to learn: {{ .WhereToLearn }}`,
	}

	prompt := promptui.Select{
		Items:        c.Spells,
		Templates:    templates,
		Size:         10,
		HideSelected: true,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	characterMenu(c)
}

func trustMenu(c lib.Character) {

}

func goldSaucerMenu(c lib.Character) {

}

func tripleTriadMenu(c lib.Character) {

}

func formatForTerminal(s string) string {
	stringWords := strings.Split(s, " ")
	var stringLine string
	var formattedString []string
	for _, word := range stringWords {
		if len(stringLine)+len(word)+1 > 80 {
			formattedString = append(formattedString, strings.TrimSpace(stringLine))
			stringLine = word
		} else {
			stringLine += " " + word
		}
	}
	return strings.Join(formattedString, "\n")
}
