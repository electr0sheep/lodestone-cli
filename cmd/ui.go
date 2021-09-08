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
	Short: "Provides a UI to navigate Lodestone data",
	Run: func(cmd *cobra.Command, args []string) {
		CHARACTER = lib.Character{Id: viper.GetString("character_id")}
		CHARACTER.GetProfile()
		g, err := gocui.NewGui(gocui.Output256)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()

		g.SetManagerFunc(makeMainMenuLayout)

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
	case "minion":
		menuLength = len(CHARACTER.Minions)
	case "mount":
		menuLength = len(CHARACTER.Mounts)
	case "achievement":
		menuLength = len(CHARACTER.Achievements)
	case "rep":
		menuLength = len(CHARACTER.Currencies) + len(CHARACTER.Reputations)
	case "spell":
		menuLength = len(CHARACTER.Spells)
	case "orchestrion":
		menuLength = len(CHARACTER.Orchestrions)
	case "card":
		menuLength = len(CHARACTER.Cards)
	case "trust":
		menuLength = len(CHARACTER.TrustCompanions)
	case "quest":
		menuLength = len(CHARACTER.Quests)
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

func switchMenu(g *gocui.Gui, nextMenu string) error {
	currentMenu := g.CurrentView().Name()
	if currentMenu != "main" {
		if err := g.DeleteView(currentMenu); err != nil {
			return err
		}
	}
	switch nextMenu {
	case "main":
		g.SetCurrentView("main")
	case "character":
		return makeCharacterMenuLayout(g)
	case "profile":
		return makeProfileMenuLayout(g)
	case "job":
		return showJobMenu(g)
	case "minion":
		return showMinionMenu(g)
	case "mount":
		return showMountMenu(g)
	case "achievement":
		return showAchievementMenu(g)
	case "rep":
		return showRepMenu(g)
	case "spell":
		return showSpellMenu(g)
	case "orchestrion":
		return showOrchestrionMenu(g)
	case "card":
		return showCardMenu(g)
	case "goldsaucer":
		return showGoldSaucerMenu(g)
	case "trust":
		return showTrustMenu(g)
	case "pvp":
		return showPvpMenu(g)
	case "quest":
		return showQuestMenu(g)
	}
	return nil
}

func makeProfileMenuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("profile", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Profile", CHARACTER.Name)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		fmt.Fprintf(v, "Name         : %s\n", CHARACTER.Name)
		fmt.Fprintf(v, "Title        : %s\n", CHARACTER.Title)
		fmt.Fprintf(v, "World        : %s\n", CHARACTER.World)
		fmt.Fprintf(v, "Race         : %s\n", CHARACTER.Race)
		fmt.Fprintf(v, "Clan         : %s\n", CHARACTER.Clan)
		fmt.Fprintf(v, "Gender       : %s\n", CHARACTER.Gender)
		fmt.Fprintf(v, "Nameday      : %s\n", CHARACTER.Nameday)
		fmt.Fprintf(v, "Guardian     : %s\n", CHARACTER.Guardian)
		fmt.Fprintf(v, "City-state   : %s\n", CHARACTER.CityState)
		fmt.Fprintf(v, "Grand Company: %s\n", CHARACTER.GrandCompany)
		fmt.Fprintf(v, "Free Company : %s\n", CHARACTER.FreeCompany)

		if _, err := g.SetCurrentView("profile"); err != nil {
			return err
		}
	}
	return nil
}

func processMainMenuSelection(g *gocui.Gui, selection string) {
	switch selection {
	case "Character":
		switchMenu(g, "character")
	case "Companions":
		showMessage(g, selection)
	case "Retainers":
		showMessage(g, selection)
	}
}

func showJobMenu(g *gocui.Gui) error {
	if CHARACTER.Jobs == nil {
		CHARACTER.GetJobs()
	}
	return makeJobMenuLayout(g)
}

func showMinionMenu(g *gocui.Gui) error {
	if CHARACTER.Minions == nil {
		CHARACTER.GetMinions()
	}
	return makeMinionMenuLayout(g)
}

func showMountMenu(g *gocui.Gui) error {
	if CHARACTER.Mounts == nil {
		CHARACTER.GetMounts()
	}
	return makeMountMenuLayout(g)
}

func showCardMenu(g *gocui.Gui) error {
	if CHARACTER.Cards == nil {
		CHARACTER.GetCards()
	}
	return makeCardMenuLayout(g)
}

func showPvpMenu(g *gocui.Gui) error {
	if (CHARACTER.PvpProfile == lib.PvpProfile{}) {
		CHARACTER.GetPvpProfile()
	}
	return makePvpMenuLayout(g)
}

func showQuestMenu(g *gocui.Gui) error {
	if CHARACTER.Quests == nil {
		CHARACTER.GetQuests()
	}
	return makeQuestMenuLayout(g)
}

func showTrustMenu(g *gocui.Gui) error {
	if CHARACTER.TrustCompanions == nil {
		CHARACTER.GetTrustCompanions()
	}
	return makeTrustMenuLayout(g)
}

func showGoldSaucerMenu(g *gocui.Gui) error {
	if (CHARACTER.GoldSaucer == lib.GoldSaucer{}) {
		CHARACTER.GetGoldSaucer()
	}
	return makeGoldSaucerMenuLayout(g)
}

func showOrchestrionMenu(g *gocui.Gui) error {
	if CHARACTER.Orchestrions == nil {
		CHARACTER.GetOrchestrions()
	}
	return makeOrchestrionMenuLayout(g)
}

func showSpellMenu(g *gocui.Gui) error {
	if CHARACTER.Spells == nil {
		CHARACTER.GetSpells()
	}
	return makeSpellMenuLayout(g)
}

func showRepMenu(g *gocui.Gui) error {
	if CHARACTER.Currencies == nil || CHARACTER.Reputations == nil {
		CHARACTER.GetCurrenciesAndRep()
	}
	return makeRepMenuLayout(g)
}

func showAchievementMenu(g *gocui.Gui) error {
	if CHARACTER.Achievements == nil {
		CHARACTER.GetAchievements()
	}
	return makeAchievementMenuLayout(g)
}

func processCharacterMenuSelection(g *gocui.Gui, selection string) {
	switch selection {
	case "Profile":
		switchMenu(g, "profile")
	case "Class/Job":
		switchMenu(g, "job")
	case "Minions":
		switchMenu(g, "minion")
	case "Mounts":
		switchMenu(g, "mount")
	case "Currencies/Reputation":
		switchMenu(g, "rep")
	case "Quests":
		switchMenu(g, "quest")
	case "Achievements":
		switchMenu(g, "achievement")
	case "Orchestrion Roll":
		switchMenu(g, "orchestrion")
	case "PvP Profile":
		switchMenu(g, "pvp")
	case "Blue Magic Spellbook":
		switchMenu(g, "spell")
	case "Trust":
		switchMenu(g, "trust")
	case "The Gold Saucer":
		switchMenu(g, "goldsaucer")
	case "Triple Triad":
		switchMenu(g, "card")
	}
}

func processJobMenuChange(g *gocui.Gui, job lib.Job) {
	makeJobDetailView(g, job)
}

func processMinionMenuChange(g *gocui.Gui, minion *lib.Minion) {
	makeMinionDetailView(g, minion)
}

func processMountMenuChange(g *gocui.Gui, mount *lib.Mount) {
	makeMountDetailView(g, mount)
}

func processAchievementMenuChange(g *gocui.Gui, achievement *lib.Achievement) {
	makeAchievementDetailView(g, achievement)
}

func processCurrencyMenuChange(g *gocui.Gui, currency lib.Currency) {
	makeCurrencyDetailView(g, currency)
}

func processReputationMenuChange(g *gocui.Gui, reputation lib.Reputation) {
	makeReputationDetailView(g, reputation)
}

func processQuestMenuChange(g *gocui.Gui, quest lib.Quest) {
	makeQuestDetailView(g, quest)
}

func processTrustMenuChange(g *gocui.Gui, trustCompanion lib.TrustCompanion) {
	makeTrustDetailView(g, trustCompanion)
}

func processCardMenuChange(g *gocui.Gui, card lib.Card) {
	makeCardDetailView(g, card)
}

func processOrchestrionMenuChange(g *gocui.Gui, orchestrion lib.Orchestrion) {
	makeOrchestrionDetailView(g, orchestrion)
}

func processSpellMenuChange(g *gocui.Gui, spell lib.Spell) {
	makeSpellDetailView(g, spell)
}

func getMountDetails(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()

	selectedMount := &CHARACTER.Mounts[oy+cy]

	if selectedMount.Description == "" {
		CHARACTER.GetMountDetails(selectedMount)
	}

	return makeMountDetailView(g, selectedMount)
}

func getMinionDetails(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()

	selectedMinion := &CHARACTER.Minions[oy+cy]

	if selectedMinion.Description == "" {
		CHARACTER.GetMinionDetails(selectedMinion)
	}

	return makeMinionDetailView(g, selectedMinion)
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
	case "minion":
		selectedMinion := &CHARACTER.Minions[oy+cy]
		processMinionMenuChange(g, selectedMinion)
	case "mount":
		selectedMount := &CHARACTER.Mounts[oy+cy]
		processMountMenuChange(g, selectedMount)
	case "achievement":
		selectedAchievement := &CHARACTER.Achievements[oy+cy]
		processAchievementMenuChange(g, selectedAchievement)
	case "rep":
		if len(CHARACTER.Currencies) > oy+cy {
			selectedCurrency := CHARACTER.Currencies[oy+cy]
			processCurrencyMenuChange(g, selectedCurrency)
		} else {
			selectedReputation := CHARACTER.Reputations[oy+cy-len(CHARACTER.Currencies)]
			processReputationMenuChange(g, selectedReputation)
		}
	case "spell":
		selectedSpell := CHARACTER.Spells[oy+cy]
		processSpellMenuChange(g, selectedSpell)
	case "orchestrion":
		selectedOrchestrion := CHARACTER.Orchestrions[oy+cy]
		processOrchestrionMenuChange(g, selectedOrchestrion)
	case "card":
		selectedCard := CHARACTER.Cards[oy+cy]
		processCardMenuChange(g, selectedCard)
	case "trust":
		selectedTrustCompanion := CHARACTER.TrustCompanions[oy+cy]
		processTrustMenuChange(g, selectedTrustCompanion)
	case "quest":
		selectedQuest := CHARACTER.Quests[oy+cy]
		processQuestMenuChange(g, selectedQuest)
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
			return switchMenu(g, "main")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("profile", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("pvp", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("goldsaucer", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return switchMenu(g, "character")
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
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("minion", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("minion", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("minion", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("minion_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("minion", gocui.KeyArrowRight, gocui.ModNone, getMinionDetails); err != nil {
		return err
	}
	if err := g.SetKeybinding("mount", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("mount", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("mount", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("mount_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("mount", gocui.KeyArrowRight, gocui.ModNone, getMountDetails); err != nil {
		return err
	}
	if err := g.SetKeybinding("achievement", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("achievement", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("achievement", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("achievement_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("achievement", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("rep", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("rep", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("reputation_detail")
			g.DeleteView("currency_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("rep", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("spell", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("spell", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("spell_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("spell", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("orchestrion", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("orchestrion", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("orchestrion_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("orchestrion", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("card", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("card", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("card_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("card", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("quest", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("quest", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("quest", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("quest_detail")
			return switchMenu(g, "character")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("trust", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorDown(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("trust", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cursorUp(g, v)
			return processMenuSelection(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("trust", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("trust_detail")
			return switchMenu(g, "character")
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
	if v, err := g.SetView("job_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = job.Name
		fmt.Fprintf(v, "Role : %s\n", job.Role)
		fmt.Fprintf(v, "Level: %s\n", job.Level)
		fmt.Fprintf(v, "XP   : %s\n", job.Xp)
	}
	return nil
}

func makeQuestDetailView(g *gocui.Gui, quest lib.Quest) error {
	g.DeleteView("quest_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("quest_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = quest.Name
		fmt.Fprintf(v, "Completion Date: %s\n", quest.CompletionDate)
		fmt.Fprintf(v, "Type           : %s\n", quest.Type)
	}

	return nil
}

func makeTrustDetailView(g *gocui.Gui, trustCompanion lib.TrustCompanion) error {
	g.DeleteView("trust_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("trust_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = trustCompanion.Name
		fmt.Fprintf(v, "Job  : %s\n", trustCompanion.Job)
		fmt.Fprintf(v, "Level: %s\n", trustCompanion.Level)
		fmt.Fprintf(v, "XP   : %s\n", trustCompanion.Xp)
		fmt.Fprintf(v, "Next : %s\n", trustCompanion.NextLevelXp)
	}

	return nil
}

func makeCardDetailView(g *gocui.Gui, card lib.Card) error {
	g.DeleteView("card_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("card_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sizeX, _ := v.Size()
		v.Title = card.Name
		fmt.Fprintf(v, "Acquired   : %t\n", card.Acquired)
		fmt.Fprintf(v, "Rarity     : %d\n", card.Rarity)
		fmt.Fprintf(v, "Type       : %s\n", card.Type)
		fmt.Fprintf(v, "Strength   :\n  %s  \n%s   %s\n  %s  \n", card.AttackUp, card.AttackLeft, card.AttackRight, card.AttackDown)
		fmt.Fprintf(v, "Description:\n%s\n", wrapStringToSize(card.Description, sizeX))
	}

	return nil
}

func makeOrchestrionDetailView(g *gocui.Gui, orchestrion lib.Orchestrion) error {
	g.DeleteView("orchestrion_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("orchestrion_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sizeX, _ := v.Size()
		v.Title = orchestrion.Name
		fmt.Fprintf(v, "Acquired     : %t\n", orchestrion.Acquired)
		fmt.Fprintf(v, "Categroy     : %s\n", orchestrion.Category)
		fmt.Fprintf(v, "Where to Find:\n%s\n", wrapStringToSize(orchestrion.WhereToFind, sizeX))
	}

	return nil
}

func makeSpellDetailView(g *gocui.Gui, spell lib.Spell) error {
	g.DeleteView("spell_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("spell_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sizeX, _ := v.Size()
		v.Title = spell.Name
		fmt.Fprintf(v, "Aspect        : %s\n", spell.Aspect)
		fmt.Fprintf(v, "Rank          : %s\n", spell.Rank)
		fmt.Fprintf(v, "Type          : %s\n", spell.Type)
		fmt.Fprintf(v, "Where to Learn:\n%s\n", spell.WhereToLearn)
		fmt.Fprintf(v, "Description   :\n%s\n", wrapStringToSize(spell.Description, sizeX))
	}

	return nil
}

func makeReputationDetailView(g *gocui.Gui, reputation lib.Reputation) error {
	g.DeleteView("currency_detail")
	g.DeleteView("reputation_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("reputation_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = reputation.Name
		fmt.Fprintf(v, "Type            : %s\n", reputation.Type)
		if reputation.ReputationLevel != "" {
			fmt.Fprintf(v, "Reputation Level: %s\n", reputation.ReputationLevel)
		}
		fmt.Fprintf(v, "Current Amount  : %s\n", reputation.CurrentAmount)
		if reputation.Maximum != "" {
			fmt.Fprintf(v, "Maximum         : %s\n", reputation.Maximum)
		}
	}

	return nil
}

func makeCurrencyDetailView(g *gocui.Gui, currency lib.Currency) error {
	g.DeleteView("currency_detail")
	g.DeleteView("reputation_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("currency_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if currency.Name != "" {
			v.Title = currency.Name
		} else {
			v.Title = currency.Type
		}
		fmt.Fprintf(v, "Type          : %s\n", currency.Type)
		fmt.Fprintf(v, "Current Amount: %s\n", currency.CurrentAmount)
		if currency.Maximum != "" {
			fmt.Fprintf(v, "Maximum       : %s\n", currency.Maximum)
		}
	}

	return nil
}

func makeAchievementDetailView(g *gocui.Gui, achievement *lib.Achievement) error {
	g.DeleteView("achievement_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("achievement_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = achievement.Name
		fmt.Fprintf(v, "Acquistion Date:\n%s\n", achievement.AcquistionDate)
	}

	return nil
}

func makeMountDetailView(g *gocui.Gui, mount *lib.Mount) error {
	g.DeleteView("mount_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("mount_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sizeX, _ := v.Size()
		v.Title = mount.Name
		if mount.Description != "" {
			fmt.Fprintf(v, "Movement       : %s\n", mount.Movement)
			fmt.Fprintf(v, "Max Riders     : %s\n", mount.MaxRiders)
			fmt.Fprintf(v, "Acquistion Date: %s\n", mount.AcquistionDate)
			fmt.Fprintf(v, "Description    :\n%s\n", wrapStringToSize(mount.Description, sizeX))
		}
	}
	return nil
}

func makeMinionDetailView(g *gocui.Gui, minion *lib.Minion) error {
	g.DeleteView("minion_detail")
	maxX, maxY := g.Size()
	if v, err := g.SetView("minion_detail", VIEW_RIGHT_BOUND, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sizeX, _ := v.Size()
		v.Title = minion.Name
		if minion.Description != "" {
			fmt.Fprintf(v, "Behavior       : %s\n", minion.Behavior)
			fmt.Fprintf(v, "Acquistion Date: %s\n", minion.AcquistionDate)
			fmt.Fprintf(v, "Description    :\n%s\n", wrapStringToSize(minion.Description, sizeX))
		}
	}
	return nil
}

func makeCardMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 19
	if v, err := g.SetView("card", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Cards", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, card := range CHARACTER.Cards {
			fmt.Fprintln(v, card.Name)
		}
		if _, err := g.SetCurrentView("card"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeGoldSaucerMenuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("goldsaucer", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Gold Saucer", CHARACTER.Name)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, "Manderville Gold Saucer Points (MGP): %s\n\n", CHARACTER.GoldSaucer.MGP)
		fmt.Fprintf(v, "Recent Triple Triad Tournament Results:\n%s\n\n", CHARACTER.GoldSaucer.TripleTriadTournamentResult)
		fmt.Fprintf(v, "Jumbo Cactpot: %s\n", CHARACTER.GoldSaucer.JumboCactpot)
		fmt.Fprintf(v, "Your Numbers: %s, %s, %s", CHARACTER.GoldSaucer.JumboCactpotNumberOne, CHARACTER.GoldSaucer.JumboCactpotNumberTwo, CHARACTER.GoldSaucer.JumboCactpotNumberThree)
		if _, err := g.SetCurrentView("goldsaucer"); err != nil {
			return err
		}
	}
	return nil
}

func makePvpMenuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 17
	if v, err := g.SetView("pvp", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>PvP", CHARACTER.Name)
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, "Faction   : %s\n", CHARACTER.PvpProfile.Faction)
		fmt.Fprintf(v, "Rank      : %s\n", CHARACTER.PvpProfile.Rank)
		fmt.Fprintf(v, "Total XP  : %s\n", CHARACTER.PvpProfile.TotalXp)
		fmt.Fprintf(v, "Current XP: %s / %s\n", CHARACTER.PvpProfile.Xp, CHARACTER.PvpProfile.NextXp)

		fmt.Fprintf(v, "\nFrontline Overall Performance\n")
		fmt.Fprintf(v, "Campaigns   : %s\n", CHARACTER.PvpProfile.OverallPerformance.Campaigns)
		fmt.Fprintf(v, "First Place : %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.OverallPerformance.FirstPlaceWins, CHARACTER.PvpProfile.OverallPerformance.FirstPlaceWinPercentage)
		fmt.Fprintf(v, "Second Place: %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.OverallPerformance.SecondPlaceWins, CHARACTER.PvpProfile.OverallPerformance.SecondPlaceWinPercentage)
		fmt.Fprintf(v, "Third Place : %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.OverallPerformance.ThirdPlaceWins, CHARACTER.PvpProfile.OverallPerformance.ThirdPlaceWinPercentage)

		fmt.Fprintf(v, "\nFrontline Weekly Performance\n")
		fmt.Fprintf(v, "Campaigns   : %s\n", CHARACTER.PvpProfile.WeeklyPerformance.Campaigns)
		fmt.Fprintf(v, "First Place : %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.WeeklyPerformance.FirstPlaceWins, CHARACTER.PvpProfile.WeeklyPerformance.FirstPlaceWinPercentage)
		fmt.Fprintf(v, "Second Place: %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.WeeklyPerformance.SecondPlaceWins, CHARACTER.PvpProfile.WeeklyPerformance.SecondPlaceWinPercentage)
		fmt.Fprintf(v, "Third Place : %s(Victory Rate: %s)\n", CHARACTER.PvpProfile.WeeklyPerformance.ThirdPlaceWins, CHARACTER.PvpProfile.WeeklyPerformance.ThirdPlaceWinPercentage)
		if _, err := g.SetCurrentView("pvp"); err != nil {
			return err
		}
	}
	return nil
}

func makeQuestMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 19
	if v, err := g.SetView("quest", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Quest", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, quest := range CHARACTER.Quests {
			fmt.Fprintln(v, quest.Name)
		}
		if _, err := g.SetCurrentView("quest"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeTrustMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 19
	if v, err := g.SetView("trust", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Trust", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, companion := range CHARACTER.TrustCompanions {
			fmt.Fprintln(v, companion.Name)
		}
		if _, err := g.SetCurrentView("trust"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeOrchestrionMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 26
	if v, err := g.SetView("orchestrion", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Orchestrions", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, orchestrion := range CHARACTER.Orchestrions {
			fmt.Fprintln(v, orchestrion.Name)
		}
		if _, err := g.SetCurrentView("orchestrion"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeSpellMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestSpellNameLength := 0
	for _, spell := range CHARACTER.Spells {
		if len(spell.Name) > longestSpellNameLength {
			longestSpellNameLength = len(spell.Name)
		}
	}
	if longestSpellNameLength+1 > len(CHARACTER.Name)+24 {
		VIEW_RIGHT_BOUND = longestSpellNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 24
	}
	if v, err := g.SetView("spell", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Blue Magic", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, spell := range CHARACTER.Spells {
			fmt.Fprintln(v, spell.Name)
		}
		if _, err := g.SetCurrentView("spell"); err != nil {
			return err
		}
		processMenuSelection(g, v)
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
	if longestJobNameLength+1 > len(CHARACTER.Name)+18 {
		VIEW_RIGHT_BOUND = longestJobNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 18
	}
	if v, err := g.SetView("job", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Jobs", CHARACTER.Name)
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

func makeRepMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestRepOrCurrencyNameLength := 0
	for _, currency := range CHARACTER.Currencies {
		if len(currency.Name) > longestRepOrCurrencyNameLength {
			longestRepOrCurrencyNameLength = len(currency.Name)
		}
	}
	for _, reputation := range CHARACTER.Reputations {
		if len(reputation.Name) > longestRepOrCurrencyNameLength {
			longestRepOrCurrencyNameLength = len(reputation.Name)
		}
	}
	if longestRepOrCurrencyNameLength+1 > len(CHARACTER.Name)+26 {
		VIEW_RIGHT_BOUND = longestRepOrCurrencyNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 26
	}
	if v, err := g.SetView("rep", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Currency/Rep", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, currency := range CHARACTER.Currencies {
			if currency.Name != "" {
				fmt.Fprintln(v, currency.Name)
			} else {
				fmt.Fprintln(v, currency.Type)
			}
		}
		for _, reputation := range CHARACTER.Reputations {
			fmt.Fprintln(v, reputation.Name)
		}
		if _, err := g.SetCurrentView("rep"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeAchievementMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestAchievementNameLength := 0
	for _, achievement := range CHARACTER.Achievements {
		if len(achievement.Name) > longestAchievementNameLength {
			longestAchievementNameLength = len(achievement.Name)
		}
	}
	if longestAchievementNameLength+1 > len(CHARACTER.Name)+26 {
		VIEW_RIGHT_BOUND = longestAchievementNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 26
	}
	if v, err := g.SetView("achievement", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Achievements", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, achievement := range CHARACTER.Achievements {
			fmt.Fprintln(v, achievement.Name)
		}
		if _, err := g.SetCurrentView("achievement"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeMountMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestMountNameLength := 0
	for _, mount := range CHARACTER.Mounts {
		if len(mount.Name) > longestMountNameLength {
			longestMountNameLength = len(mount.Name)
		}
	}
	if longestMountNameLength+1 > len(CHARACTER.Name)+20 {
		VIEW_RIGHT_BOUND = longestMountNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 20
	}
	if v, err := g.SetView("mount", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Mounts", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, mount := range CHARACTER.Mounts {
			fmt.Fprintln(v, mount.Name)
		}
		if _, err := g.SetCurrentView("mount"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeMinionMenuLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	longestMinionNameLength := 0
	for _, minion := range CHARACTER.Minions {
		if len(minion.Name) > longestMinionNameLength {
			longestMinionNameLength = len(minion.Name)
		}
	}
	if longestMinionNameLength+1 > len(CHARACTER.Name)+21 {
		VIEW_RIGHT_BOUND = longestMinionNameLength + 1
	} else {
		VIEW_RIGHT_BOUND = len(CHARACTER.Name) + 21
	}
	if v, err := g.SetView("minion", 0, 0, VIEW_RIGHT_BOUND, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character>Minions", CHARACTER.Name)
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, minion := range CHARACTER.Minions {
			fmt.Fprintln(v, minion.Name)
		}
		if _, err := g.SetCurrentView("minion"); err != nil {
			return err
		}
		processMenuSelection(g, v)
	}
	return nil
}

func makeCharacterMenuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("character", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s>Character", CHARACTER.Name)
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

func makeMainMenuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		VIEW_RIGHT_BOUND = 10
		v.Title = CHARACTER.Name
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

func wrapStringToSize(s string, maxLength int) string {
	stringWords := strings.Split(s, " ")
	var stringLine string
	var formattedString []string
	for _, word := range stringWords {
		if len(stringLine)+len(word)+1 > maxLength {
			formattedString = append(formattedString, strings.TrimSpace(stringLine))
			stringLine = word
		} else {
			stringLine += " " + word
		}
	}
	formattedString = append(formattedString, strings.TrimSpace(stringLine))

	return strings.Join(formattedString, "\n")
}
