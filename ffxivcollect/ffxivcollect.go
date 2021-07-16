package ffxivcollect

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

type BlueMagicSpell struct {
	Name     string
	Id       int
	Obtained bool
}

type Orchestrion struct {
	Name     string
	Id       int
	Obtained bool
}

type TripleTriadCard struct {
	Name     string
	Id       int
	Obtained bool
}

func validateAndUpdateTriadSessionToken(scripts []string) bool {
	if len(scripts) >= 10 {
		ffxiv_triple_triad_authenticity_token := scripts[9]
		viper.Set("ffxiv_triple_triad_authenticity_token", ffxiv_triple_triad_authenticity_token)
		viper.WriteConfig()
		return true
	} else {
		viper.Set("ffxiv_triple_triad_authenticity_token", "")
		viper.Set("ffxiv_triple_triad_session_token", "")
		viper.WriteConfig()
		return false
	}
}

func getSessionToken() {
	tokenPrompt := promptui.Prompt{
		Label: "FFXIV Collect Session Token",
	}
	ffxiv_collect_session_token, err := tokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	viper.Set("ffxiv_collect_session_token", ffxiv_collect_session_token)
	viper.WriteConfig()
}

func getTripleTriadSessionToken() {
	tokenPrompt := promptui.Prompt{
		Label: "Triad Session Token",
	}
	triad_session_token, err := tokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	viper.Set("ffxiv_triple_triad_session_token", triad_session_token)
	viper.WriteConfig()
}

func AddBlueMagicSpell(spell_name string, spell_id int) bool {
	client := &http.Client{}
	data := url.Values{}
	data.Set("authenticity_token", viper.GetString("ffxiv_collect_authenticity_token"))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://ffxivcollect.com/spells/%d/add", spell_id), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_collect_session=%s", viper.GetString("ffxiv_collect_session_token")))

	// for now, we don't care about the response, so just make the request
	resp, _ := client.Do(req)
	if resp.StatusCode == 422 {
		getSessionToken()
		data.Set("authenticity_token", viper.GetString("ffxiv_collect_authenticity_token"))
		req, err = http.NewRequest("POST", fmt.Sprintf("https://ffxivcollect.com/spells/%d/add", spell_id), strings.NewReader(data.Encode()))
	}
	return resp.StatusCode == 204
}

func AddOrchestrion(orchestrion_name string, orchestrion_id int) bool {
	client := &http.Client{}
	data := url.Values{}
	data.Set("authenticity_token", viper.GetString("ffxiv_collect_authenticity_token"))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://ffxivcollect.com/orchestrions/%d/add", orchestrion_id), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_collect_session=%s", viper.GetString("ffxiv_collect_session_token")))

	// for now, we don't care about the response, so just make the request
	resp, _ := client.Do(req)
	if resp.StatusCode == 422 {
		getSessionToken()
		data.Set("authenticity_token", viper.GetString("ffxiv_collect_authenticity_token"))
		req, err = http.NewRequest("POST", fmt.Sprintf("https://ffxivcollect.com/orchestrions/%d/add", orchestrion_id), strings.NewReader(data.Encode()))
	}
	return resp.StatusCode == 204
}

func AddCard(card_name string, card_id int) bool {
	client := &http.Client{}
	data := url.Values{}
	data.Set("authenticity_token", viper.GetString("ffxiv_triple_triad_authenticity_token"))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://triad.raelys.com/cards/%d/add", card_id), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_triple_triad_session=%s", viper.GetString("ffxiv_triple_triad_session_token")))

	// for now, we don't care about the response, so just make the request
	resp, _ := client.Do(req)
	if resp.StatusCode == 422 {
		getSessionToken()
		data.Set("authenticity_token", viper.GetString("ffxiv_triple_triad_authenticity_token"))
		req, err = http.NewRequest("POST", fmt.Sprintf("https://triad.raelys.com/cards/%d/add", card_id), strings.NewReader(data.Encode()))
	}
	return resp.StatusCode == 204
}

func GetBlueMagicSpells() map[string]BlueMagicSpell {
	if viper.GetString("ffxiv_collect_session_token") == "" {
		getSessionToken()
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ffxivcollect.com/spells", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_collect_session=%s", viper.GetString("ffxiv_collect_session_token")))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// update authenticity token
	ffxiv_collect_authenticity_token := strings.Split(doc.Find("script").Text(), "'")[9]
	viper.Set("ffxiv_collect_authenticity_token", ffxiv_collect_authenticity_token)
	viper.WriteConfig()

	blueMagicSpellMap := make(map[string]BlueMagicSpell)
	blueMagicElements := doc.Find(".collectable")
	blueMagicElements.Each(func(_ int, blueMagicElement *goquery.Selection) {
		name := blueMagicElement.Find(".name").Text()
		id, _ := blueMagicElement.Find(".name").Attr("href")
		id = strings.Split(id, "/")[2]
		converted_id, _ := strconv.Atoi(id)
		obtained := blueMagicElement.HasClass("owned")
		blueMagicSpellMap[name] = BlueMagicSpell{Name: name, Id: converted_id, Obtained: obtained}
	})
	return blueMagicSpellMap
}

func GetOrchestrions() map[string]Orchestrion {
	if viper.GetString("ffxiv_collect_session_token") == "" {
		getSessionToken()
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ffxivcollect.com/orchestrions", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_collect_session=%s", viper.GetString("ffxiv_collect_session_token")))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// update authenticity token
	ffxiv_collect_authenticity_token := strings.Split(doc.Find("script").Text(), "'")[9]
	viper.Set("ffxiv_collect_authenticity_token", ffxiv_collect_authenticity_token)
	viper.WriteConfig()

	orchestrionMap := make(map[string]Orchestrion)
	orchestrionElements := doc.Find(".collectable")
	orchestrionElements.Each(func(_ int, orchestrionElement *goquery.Selection) {
		name := orchestrionElement.Find(".name").Text()
		id, _ := orchestrionElement.Find(".name").Attr("href")
		id = strings.Split(id, "/")[2]
		converted_id, _ := strconv.Atoi(id)
		obtained := orchestrionElement.HasClass("owned")
		orchestrionMap[name] = Orchestrion{Name: name, Id: converted_id, Obtained: obtained}
	})
	return orchestrionMap
}

func GetCards() map[string]TripleTriadCard {
	if viper.GetString("ffxiv_triple_triad_session_token") == "" {
		getTripleTriadSessionToken()
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://triad.raelys.com/cards/mine", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Cookie", fmt.Sprintf("_ffxiv_triple_triad_session=%s", viper.GetString("ffxiv_triple_triad_session_token")))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// validate and update authenticity token
	if !validateAndUpdateTriadSessionToken(strings.Split(doc.Find("script").Text(), "'")) {
		fmt.Println("The triple triad session token appears to be wrong, make sure you have logged in first and try again")
		return nil
	}

	cardMap := make(map[string]TripleTriadCard)
	cardElements := doc.Find(".card-row")
	cardElements.Each(func(_ int, cardElement *goquery.Selection) {
		name := cardElement.Find(".name").Text()
		id, _ := cardElement.Find(".name").Attr("href")
		id = strings.Split(id, "/")[2]
		converted_id, _ := strconv.Atoi(id)
		obtained := cardElement.HasClass("has-card")
		cardMap[name] = TripleTriadCard{Name: name, Id: converted_id, Obtained: obtained}
	})
	return cardMap
}
