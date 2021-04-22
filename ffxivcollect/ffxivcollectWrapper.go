package ffxivcollectWrapper

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

func getSessionToken() {
	tokenPrompt := promptui.Prompt{
		Label: "FFXIV Collect Session Token",
	}
	ffxiv_collect_session_token, err := tokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	authenticityTokenPrompt := promptui.Prompt{
		Label: "FFXIV Collect Authenticity Token",
	}
	ffxiv_collect_authenticity_token, err := authenticityTokenPrompt.Run()
	if err != nil {
		panic(err)
	}
	viper.Set("ffxiv_collect_session_token", ffxiv_collect_session_token)
	viper.Set("ffxiv_collect_authenticity_token", ffxiv_collect_authenticity_token)
	viper.WriteConfig()
}

// blueMagicSpellMap is a map containing the names and ffxivcollect.com ids for
// the spells. For example, "Water Cannon": 3
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

func GetBlueMagicSpells() map[string]BlueMagicSpell {
	// TODO: We need to figure out how to make sure the user is logged in for this
	// If the user is not logged in, we can't tell which spells we need to add
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
