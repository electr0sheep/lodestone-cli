/*
Copyright © 2021 electr0sheep electr0sheep@electr0sheep.com

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
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync [character ID] [session token]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		character_id := args[0]
		session_token := args[1]

		mounts := getMounts(character_id)
		fmt.Println("MOUNTS")
		for _, mount := range mounts {
			fmt.Println(mount)
		}

		fmt.Printf("\n\n\n")

		minions := getMinions(character_id)
		fmt.Println("MINIONS")
		for _, minion := range minions {
			fmt.Println(minion)
		}

		fmt.Printf("\n\n\n")

		orchestrions := getOrchestrions(character_id, session_token)
		fmt.Println("ORCHESTRIONS")
		for _, orchestrion := range orchestrions {
			fmt.Println(orchestrion)
		}

		fmt.Printf("\n\n\n")

		spells := getSpells(character_id, session_token)
		fmt.Println("SPELLS")
		for _, spell := range spells {
			fmt.Println(spell)
		}

		fmt.Printf("\n\n\n")

		achievements := getAchievements(character_id, session_token)
		fmt.Println("ACHIEVEMENTS")
		for _, achievement := range achievements {
			fmt.Println(achievement)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Gets mounts from Lodestone
func getMounts(character_id string) []string {
	client := &http.Client{}
	req := setupRequest("mount", character_id, "")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		panic("NOOOOO")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	mountElements := doc.Find(".mount__name")

	var mounts []string

	mountElements.Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		mounts = append(mounts, name)
	})

	return mounts
}

// Gets minions from Lodestone
func getMinions(character_id string) []string {
	client := &http.Client{}
	req := setupRequest("minion", character_id, "")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		panic("NOOOOO")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	minionElements := doc.Find(".minion__name")

	var minions []string

	minionElements.Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		minions = append(minions, name)
	})

	return minions
}

// Gets orchestrions from Lodestone
func getOrchestrions(character_id string, session_token string) []string {
	client := &http.Client{}
	req := setupRequest("orchestrion", character_id, session_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		panic("NOOOOO")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	orchestrionElements := doc.Find("li:not([class])").Find(".orchestrion-list__name")

	var orchestrions []string

	orchestrionElements.Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		orchestrions = append(orchestrions, name)
	})

	return orchestrions
}

// Gets blue magic spellbook from Lodestone
func getSpells(character_id string, session_token string) []string {
	client := &http.Client{}
	req := setupRequest("bluemage", character_id, session_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		panic("NOOOOO")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	spellElements := doc.Find(".sys-reward").Find(".bluemage-action__name")

	var spells []string

	spellElements.Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		spells = append(spells, name)
	})

	return spells
}

// Gets achievements from Lodestone
func getAchievements(character_id string, session_token string) []string {
	client := &http.Client{}
	req := setupRequest("achievement", character_id, session_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		panic("NOOOOO")
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	achievementsElements := doc.Find(".entry__achievement--complete")

	var achievements []string

	achievementsElements.Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		achievements = append(achievements, name)
	})

	return achievements
}

// Sets up a request
func setupRequest(endpoint string, character_id string, session_token string) *http.Request {
	USER_AGENT := "Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.76 Mobile Safari/537.36"

	req, err := http.NewRequest("GET", fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s/%s", character_id, endpoint), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)

	if session_token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("ldst_sess=%s;", session_token))
	}

	return req
}
