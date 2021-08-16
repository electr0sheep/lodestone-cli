package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/electr0sheep/lodestone-cli/lib"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CHARACTER lib.Character
var CURRENT_MENU string = ""
var VIEW_RIGHT_BOUND int
var MAIN_MENU_OPTIONS = []string{
	"Character",
	"Companions",
	"Retainers"}
var CHARACTER_MENU_OPTIONS = []string{
	"Profile",
	"Class/Job",
	"Minions",
	"Mounts",
	"Currencies/Reputation",
	"Quests",
	"Achievements",
	"Orchestrion Roll",
	"PvP Profile",
	"Blue Magic Spellbook",
	"Trust",
	"The Gold Saucer",
	"Triple Triad"}

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		CHARACTER = lib.Character{Id: viper.GetString("character_id")}
		CHARACTER.GetProfile()
		g, err := gocui.NewGui(gocui.Output256)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()

		g.SetManagerFunc(layout)

		if err := keybindings(g); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	menuLength := 0
	switch v.Name() {
	case "main":
		menuLength = len(MAIN_MENU_OPTIONS)
	case "character":
		menuLength = len(CHARACTER_MENU_OPTIONS)
	case "job":
		menuLength = len(CHARACTER.Jobs)
	}

	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy+oy < menuLength-1 {
			if err := v.SetCursor(cx, cy+1); err != nil {
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func showMessage(g *gocui.Gui, message string) error {
	CURRENT_MENU = g.CurrentView().Name()
	extraXSpace := 0
	extraYSpace := 0
	lines := strings.Split(message, "\n")
	maxLineLength := 0
	for _, line := range lines {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}
	stringYSize := strings.Count(message, "\n") + 1
	if stringYSize%2 != 0 {
		extraYSpace = 1
	}
	if maxLineLength%2 != 0 {
		extraXSpace = 1
	}
	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-maxLineLength/2, maxY/2-stringYSize/2, maxX/2+maxLineLength/2+1+extraXSpace, maxY/2+stringYSize/2+1+extraYSpace); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, message)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func processMainMenuSelection(g *gocui.Gui, selection string) {
	switch selection {
	case "Character":
		makeCharacterMenuLayout(g)
	case "Companions":
		showMessage(g, selection)
	case "Retainers":
		showMessage(g, selection)
	}
}

func showProfilePopup(g *gocui.Gui) {
	showMessage(g,
		fmt.Sprintf(`Title        : %s
World        : %s
Race         : %s
Clan         : %s
Gender       : %s
Nameday      : %s
Guardian     : %s
City-state   : %s
Grand Company: %s
Free Company : %s`,
			CHARACTER.Title,
			CHARACTER.World,
			CHARACTER.Race,
			CHARACTER.Clan,
			CHARACTER.Gender,
			CHARACTER.Nameday,
			CHARACTER.Guardian,
			CHARACTER.CityState,
			CHARACTER.GrandCompany,
			CHARACTER.FreeCompany))

}

func showJobMenu(g *gocui.Gui) {
	if CHARACTER.Jobs == nil {
		CHARACTER.GetJobs()
	}
	makeJobMenuLayout(g)
}

func processCharacterMenuSelection(g *gocui.Gui, selection string) {
	switch selection {
	case "Profile":
		showProfilePopup(g)
	case "Class/Job":
		showJobMenu(g)
	case "Minions":
		showMessage(g, selection)
	case "Mounts":
		showMessage(g, selection)
	case "Currencies/Reputation":
		showMessage(g, selection)
	case "Quests":
		showMessage(g, selection)
	case "Achievements":
		showMessage(g, selection)
	case "Orchestrion Roll":
		showMessage(g, selection)
	case "PvP Profile":
		showMessage(g, selection)
	case "Blue Magic Spellbook":
		showMessage(g, selection)
	case "Trust":
		showMessage(g, selection)
	case "The Gold Saucer":
		showMessage(g, selection)
	case "Triple Triad":
		showMessage(g, selection)
	}
}

func processJobMenuChange(g *gocui.Gui, job lib.Job) {
	makeJobDetailView(g, job)
}

func processMenuSelection(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()

	switch v.Name() {
	case "main":
		selectedOption := MAIN_MENU_OPTIONS[oy+cy]
		processMainMenuSelection(g, selectedOption)
	case "character":
		selectedOption := CHARACTER_MENU_OPTIONS[oy+cy]
		processCharacterMenuSelection(g, selectedOption)
	case "job":
		selectedJob := CHARACTER.Jobs[oy+cy]
		processJobMenuChange(g, selectedJob)
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(CURRENT_MENU); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, processMenuSelection); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, processMenuSelection); err != nil {
		return err
	}
	if err := g.SetKeybinding("character", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("character", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("character", gocui.KeyEnter, gocui.ModNone, processMenuSelection); err != nil {
		return err
	}
	if err := g.SetKeybinding("character", gocui.KeyArrowRight, gocui.ModNone, processMenuSelection); err != nil {
		return err
	}
	if err := g.SetKeybinding("character", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return returnToMenu(g, v, "main")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("job", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("job", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("job", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("job_detail")
			return returnToMenu(g, v, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyArrowLeft, gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}
func makeJobDetailView(g *gocui.Gui, job lib.Job) error {
	g.DeleteView("job_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("job_detail", VIEW_RIGHT_BOUND, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v,
			`Name : %s
Role : %s
Level: %s
XP   : %s`,
			job.Name,
			job.Role,
			job.Level,
			job.Xp)
	}
	return nil
}

func makeJobMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestJobNameLength := 0
	for _, job := range CHARACTER.Jobs {
		if len(job.Name) > longestJobNameLength {
			longestJobNameLength = len(job.Name)
		}
	}
	if v, err := g.SetView("job", VIEW_RIGHT_BOUND, -1, VIEW_RIGHT_BOUND+1+longestJobNameLength, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		VIEW_RIGHT_BOUND = VIEW_RIGHT_BOUND + longestJobNameLength + 1
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, job := range CHARACTER.Jobs {
			fmt.Fprintln(v, job.Name)
		}
		if _, err := g.SetCurrentView("job"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeCharacterMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("character", VIEW_RIGHT_BOUND, -1, VIEW_RIGHT_BOUND+22, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		VIEW_RIGHT_BOUND = VIEW_RIGHT_BOUND + 22
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, option := range CHARACTER_MENU_OPTIONS {
			fmt.Fprintln(v, option)
		}
		if _, err := g.SetCurrentView("character"); err != nil {
			return err
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("main", -1, -1, 10, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		VIEW_RIGHT_BOUND = 10
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, option := range MAIN_MENU_OPTIONS {
			fmt.Fprintln(v, option)
		}
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func returnToMenu(g *gocui.Gui, v *gocui.View, menuName string) error {
	viewX, _ := v.Size()
	if err := g.DeleteView(g.CurrentView().Name()); err != nil {
		return err
	}
	VIEW_RIGHT_BOUND = VIEW_RIGHT_BOUND - viewX - 1
	if _, err := g.SetCurrentView(menuName); err != nil {
		return err
	}
	return nil
}
