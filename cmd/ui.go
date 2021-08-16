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
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

var CURRENT_MENU string = ""
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

// func nextView(g *gocui.Gui, v *gocui.View) error {
// 	if v == nil || v.Name() == "side" {
// 		_, err := g.SetCurrentView("main")
// 		return err
// 	}
// 	_, err := g.SetCurrentView("side")
// 	return err
// }

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	menuLength := 0
	switch g.CurrentView().Name() {
	case "main":
		menuLength = len(MAIN_MENU_OPTIONS)
	case "character":
		menuLength = len(CHARACTER_MENU_OPTIONS)
	}

	if v != nil {
		cx, cy := v.Cursor()
		if cy < menuLength-1 {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
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
		// showMessage(g, selection)
	case "Companions":
		showMessage(g, selection)
	case "Retainers":
		showMessage(g, selection)
	}
}

func processCharacterMenuSelection(g *gocui.Gui, selection string) {
	switch selection {
	case "Profile":
		showMessage(g, selection)
	case "Class/Job":
		showMessage(g, selection)
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

func processMenuSelection(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	switch g.CurrentView().Name() {
	case "main":
		selectedOption := MAIN_MENU_OPTIONS[cy]
		processMainMenuSelection(g, selectedOption)
	case "character":
		selectedOption := CHARACTER_MENU_OPTIONS[cy]
		processCharacterMenuSelection(g, selectedOption)
	}
	// showMessage(g, selectedOption)
	// showMessage(g, "TEST\nTEST1234\nTEST\nTEST")
	// var l string
	// var err error

	// _, cy := v.Cursor()
	// if l, err = v.Line(cy); err != nil {
	// 	l = ""
	// }

	// maxX, maxY := g.Size()
	// if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	fmt.Fprintln(v, l)
	// 	if _, err := g.SetCurrentView("msg"); err != nil {
	// 		return err
	// 	}
	// }
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
	// if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
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
	if err := g.SetKeybinding("character", gocui.KeyArrowLeft, gocui.ModNone, returnToMainMenu); err != nil {
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
	// if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
	// 	return err
	// }
	return nil
}

// func saveMain(g *gocui.Gui, v *gocui.View) error {
// 	f, err := ioutil.TempFile("", "gocui_demo_")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	p := make([]byte, 5)
// 	v.Rewind()
// 	for {
// 		n, err := v.Read(p)
// 		if n > 0 {
// 			if _, err := f.Write(p[:n]); err != nil {
// 				return err
// 			}
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
// 	f, err := ioutil.TempFile("", "gocui_demo_")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	vb := v.ViewBuffer()
// 	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
// 		return err
// 	}
// 	return nil
// }

func makeCharacterMenuLayout(g *gocui.Gui) error {
	// g.CurrentView().Highlight = false
	_, maxY := g.Size()
	if v, err := g.SetView("character", 10, -1, 32, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
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
	// maxX, maxY := g.Size()
	_, maxY := g.Size()
	if v, err := g.SetView("main", -1, -1, 10, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, option := range MAIN_MENU_OPTIONS {
			fmt.Fprintln(v, option)
		}
		// fmt.Fprintln(v, "Character")
		// fmt.Fprintln(v, "Companions")
		// fmt.Fprintln(v, "Retainers")
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	// if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Fprintf(v, "%s", b)
	// 	v.Editable = true
	// 	v.Wrap = true
	// 	if _, err := g.SetCurrentView("side"); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func returnToMainMenu(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(g.CurrentView().Name()); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}
