package lodestoneWrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Gets a session token from Lodestone to read private data
func getSessionToken() {
	// Some crap I wrote to see if we could get our own tokens
	// Maybe I'll work more on this later?

	// idPrompt := promptui.Prompt{
	// 	Label: "Square Enix ID",
	// }
	// id, err := idPrompt.Run()
	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return ""
	// }
	// if id == "" {
	// 	return ""
	// }

	// passwordPrompt := promptui.Prompt{
	// 	Label: "Square Enix Password",
	// 	Mask:  '*',
	// }
	// password, err := passwordPrompt.Run()
	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return ""
	// }
	// if password == "" {
	// 	return ""
	// }

	// fmt.Printf("Square Enix ID: %s\n", id)
	// fmt.Printf("Square Enix Password: %s\n", password)
	// return ""
	tokenPrompt := promptui.Prompt{
		Label: "Lodestone Session Token",
	}
	lodestone_session_token, err := tokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	viper.Set("lodestone_session_token", lodestone_session_token)
	viper.WriteConfig()
}

// Sets up a request
func setupRequest(endpoint string, character_id string) *http.Request {
	USER_AGENT := "Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.76 Mobile Safari/537.36"

	lodestone_session_token := viper.Get("lodestone_session_token")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s/%s", character_id, endpoint), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)

	if lodestone_session_token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("ldst_sess=%s;", lodestone_session_token))
	}

	return req
}

// Gets achievements from Lodestone
func GetAchievements(character_id string) []string {
	client := &http.Client{}
	morePages := true
	var achievements []string

	for page := 1; morePages; page++ {
		progressIndicator := map[int]string{
			0: ".  ",
			1: ".. ",
			2: "...",
		}
		fmt.Printf("\rGetting achievement page %d%s", page, progressIndicator[(page-1)%3])
		req := setupRequest(fmt.Sprintf("achievement/?page=%d", page), character_id)

		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}
		achievementsElements := doc.Find(".entry__achievement").Find(".entry__activity__txt")

		if achievementsElements.Length() == 0 {
			fmt.Printf("\r                                 \r")
			morePages = false
		} else {
			achievementsElements.Each(func(_ int, achievementsElement *goquery.Selection) {
				name := achievementsElements.Text()
				name = strings.Split(name, "\"")[1]
				achievements = append(achievements, name)
			})
		}
	}

	return achievements
}

// Gets minions from Lodestone
func GetMinions(character_id string) []string {
	client := &http.Client{}

	req := setupRequest("minion", character_id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	minionElements := doc.Find(".minion__name")

	var minions []string

	minionElements.Each(func(_ int, minionElement *goquery.Selection) {
		name := minionElement.Text()
		minions = append(minions, name)
	})

	return minions
}

// Gets mounts from Lodestone
func GetMounts(character_id string) []string {
	client := &http.Client{}
	req := setupRequest("mount", character_id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	mountElements := doc.Find(".mount__name")

	var mounts []string

	mountElements.Each(func(_ int, mountElement *goquery.Selection) {
		name := mountElement.Text()
		mounts = append(mounts, name)
	})

	return mounts
}

// Gets orchestrions from Lodestone
func GetOrchestrions(character_id string) []string {
	client := &http.Client{}
	req := setupRequest("orchestrion", character_id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	orchestrionElements := doc.Find("li:not([class])").Find(".orchestrion-list__name")

	var orchestrions []string

	orchestrionElements.Each(func(_ int, orchestrionElement *goquery.Selection) {
		name := orchestrionElement.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		orchestrions = append(orchestrions, name)
	})

	return orchestrions
}

// Gets blue magic spellbook from Lodestone
func GetSpells(character_id string) []string {
	client := &http.Client{}
	req := setupRequest("bluemage", character_id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		getSessionToken()
		req = setupRequest("bluemage", character_id)
		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == 404 {
			cobra.CheckErr("There was an error retrieiving data from Lodestone. Is your session token correct?")
		}
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	spellElements := doc.Find(".sys-reward").Find(".bluemage-action__name")

	var spells []string

	spellElements.Each(func(_ int, spellElement *goquery.Selection) {
		name := spellElement.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		spells = append(spells, name)
	})

	return spells
}
