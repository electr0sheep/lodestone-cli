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

type Character struct {
	Name         string
	Id           string
	Achievements []Achievement
	Cards        []Card
	Jobs         []Job
	Minions      []Minion
	Mounts       []Mount
	Orchestrions []Orchestrion
	Retainers    []Retainer
	Spells       []Spell
}

func (c Character) GetAchievements() []Achievement {
	client := &http.Client{}
	morePages := true
	var achievements []Achievement

	for page := 1; morePages; page++ {
		progressIndicator := map[int]string{
			0: ".  ",
			1: ".. ",
			2: "...",
		}
		fmt.Printf("\rGetting achievement page %d%s", page, progressIndicator[(page-1)%3])
		req := setupMobileRequest(fmt.Sprintf("achievement/?page=%d", page), c.Id)

		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}
		achievementElements := doc.Find(".entry__achievement").Find(".entry__activity__txt")

		if achievementElements.Length() == 0 {
			fmt.Printf("\r                                 \r")
			morePages = false
		} else {
			achievementElements.Each(func(_ int, achievementElement *goquery.Selection) {
				name := achievementElement.Text()
				name = strings.Split(name, "\"")[1]
				achievements = append(achievements, Achievement{Name: name})
			})
		}
	}

	return achievements
}

func (c Character) GetCards() []Card {
	client := &http.Client{}
	morePages := true
	var cards []Card

	for page := 1; morePages; page++ {
		progressIndicator := map[int]string{
			0: ".  ",
			1: ".. ",
			2: "...",
		}
		fmt.Printf("\rGetting card page %d%s", page, progressIndicator[(page-1)%3])
		req := setupRequest(fmt.Sprintf("goldsaucer/tripletriad/?hold=1&page=%d", page), c.Id)

		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}
		cardElements := doc.Find(".name_inner")

		if cardElements.Length() == 0 {
			fmt.Printf("\r                                 \r")
			morePages = false
		} else {
			cardElements.Each(func(_ int, cardElement *goquery.Selection) {
				name := cardElement.Text()
				// there might be a better way to exclude cards that the character
				// doesn't have, but this will be good enough for now
				if name != "???" {
					cards = append(cards, Card{Name: name})
				}
			})
		}
	}

	return cards
}

func (c Character) GetJobs() []Job {
	client := &http.Client{}
	var jobs []Job

	req := setupRequest("class_job", c.Id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	roleElements := doc.Find(".character__job__role")

	roleElements.Each(func(_ int, roleElement *goquery.Selection) {
		role := roleElement.Find(".heading--lead").Nodes[0].LastChild.Data
		jobElements := roleElement.Find(".character__job li")

		jobElements.Each(func(_ int, jobElement *goquery.Selection) {
			name := jobElement.Find(".character__job__name").Text()
			level := jobElement.Find(".character__job__level").Text()
			xp := jobElement.Find(".character__job__exp").Text()
			jobs = append(jobs, Job{Name: name, Level: level, Role: role, Xp: xp})
		})
	})

	return jobs
}

func (c Character) GetMinions() []Minion {
	client := &http.Client{}
	var minions []Minion

	req := setupMobileRequest("minion", c.Id)

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

	minionElements.Each(func(_ int, minionElement *goquery.Selection) {
		name := minionElement.Text()
		minions = append(minions, Minion{Name: name})
	})

	return minions
}

func (c Character) GetMounts() []Mount {
	client := &http.Client{}
	req := setupMobileRequest("mount", c.Id)
	var mounts []Mount

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

	mountElements.Each(func(_ int, mountElement *goquery.Selection) {
		name := mountElement.Text()
		mounts = append(mounts, Mount{Name: name})
	})

	return mounts
}

func (c Character) GetOrchestrions() []Orchestrion {
	client := &http.Client{}
	req := setupMobileRequest("orchestrion", c.Id)
	var orchestrions []Orchestrion

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

	orchestrionElements.Each(func(_ int, orchestrionElement *goquery.Selection) {
		name := orchestrionElement.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		// We need to massage the data a little bit. The names of the orchestrions
		// in ffxivcollect are titles (i.e. The Maiden's Lament as opposed to The maiden's lament)
		name = strings.Title(name)
		name = strings.ReplaceAll(name, "'S", "'s")
		name = strings.ReplaceAll(name, "Ul'Dah", "Ul'dah")
		orchestrions = append(orchestrions, Orchestrion{Name: name})
	})

	return orchestrions
}

func (c Character) GetRetainers() []*Retainer {
	client := &http.Client{}
	req := setupRequest("retainer", c.Id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		cobra.CheckErr("There was an error retrieiving data from Lodestone. Is your session token correct?")
	}

	if resp.StatusCode == 404 {
		getSessionToken()
		req := setupRequest("retainer", c.Id)
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
	retainerElements := doc.Find(".parts__switch__link")

	var retainers []*Retainer

	retainerElements.Each(func(_ int, retainerElement *goquery.Selection) {
		name := retainerElement.Text()
		id, _ := retainerElement.Attr("href")
		id = strings.Split(id, "/")[5]
		retainers = append(retainers, &Retainer{Name: name, Id: id, OwnerId: c.Id})
	})

	for _, retainer := range retainers {
		retainer.GetItems()
	}

	return retainers
}

func (c Character) GetSpells() []Spell {
	client := &http.Client{}
	req := setupMobileRequest("bluemage", c.Id)
	var spells []Spell

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		getSessionToken()
		req = setupMobileRequest("bluemage", c.Id)
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

	spellElements.Each(func(_ int, spellElement *goquery.Selection) {
		name := spellElement.Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		// TODO: Move this crap into a separate massager class?
		name = strings.ReplaceAll(name, " of", " Of")
		name = strings.ReplaceAll(name, " the", " The")
		spells = append(spells, Spell{Name: name})
	})

	return spells
}

type Achievement struct {
	Name string
}

type Card struct {
	Name string
}

type Job struct {
	Name  string
	Level string
	Role  string
	Xp    string
}

type Minion struct {
	Name string
}

type Mount struct {
	Name string
}

type Orchestrion struct {
	Name string
}

type Spell struct {
	Name string
}

type Retainer struct {
	Name    string
	OwnerId string
	Id      string
	Items   []Item
}

func (r *Retainer) GetItems() {
	client := &http.Client{}
	req := setupRequest(fmt.Sprintf("retainer/%s/baggage", r.Id), r.OwnerId)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		cobra.CheckErr("There was an error retrieiving data from Lodestone. Is your session token correct?")
	}

	if resp.StatusCode == 404 {
		getSessionToken()
		req := setupRequest(fmt.Sprintf("retainer/%s/baggage", r.Id), r.OwnerId)
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
	retainerItemElements := doc.Find(".item-list__list")

	r.Items = nil

	retainerItemElements.Each(func(_ int, retainerItemElement *goquery.Selection) {
		name := retainerItemElement.Find(".item-list__name").Find(".db-tooltip__item__txt").Find(".db-tooltip__item__name").Text()
		quantity := retainerItemElement.Find(".item-list__number").Text()
		highQuality := strings.Contains(name, "")
		canBePlacedInArmoire := strings.Contains(retainerItemElement.Text(), "Cannot be placed in an armoire.")
		isUnique := retainerItemElement.Find(".rare").Nodes != nil
		itemCategory := retainerItemElement.Find(".db-tooltip__item__category").First().Text()
		if highQuality {
			name = strings.TrimSuffix(name, "")
		}
		r.Items = append(r.Items, Item{Name: name, Quantity: quantity, HighQuality: highQuality, CanBePlacedInArmoire: canBePlacedInArmoire, IsUnique: isUnique, ItemCategory: itemCategory})
	})
}

type Item struct {
	Name                 string
	Quantity             string
	HighQuality          bool
	CanBePlacedInArmoire bool
	IsUnique             bool
	ItemCategory         string
}

func (i Item) IsStackable() bool {
	return i.Quantity != "99" && !i.IsUnique && !i.IsMinion() && !i.IsGear() && !i.IsFurnishing() && !i.IsBarding()
}

func (i Item) IsBarding() bool {
	return strings.Contains(i.Name, "Barding")
}

func (i Item) IsFurnishing() bool {
	for _, furnishingCategory := range [2]string{
		"Outdoor Furnishing",
		"Furnishing"} {
		if i.ItemCategory == furnishingCategory {
			return true
		}
	}
	return false
}

func (i Item) IsMinion() bool {
	return i.ItemCategory == "Minion"
}

func (i Item) IsGear() bool {
	for _, gearCategory := range [53]string{
		"Earrings",
		"Necklace",
		"Bracelets",
		"Ring",
		"Shield",
		"Head",
		"Body",
		"Hands",
		"Waist",
		"Legs",
		"Feet",
		"Carpenter's Primary Tool",
		"Blacksmith's Primary Tool",
		"Armorer's Primary Tool",
		"Goldsmith's Primary Tool",
		"Leatherworker's Primary Tool",
		"Weaver's Primary Tool",
		"Alchemist's Primary Tool",
		"Culinarian's Primary Tool",
		"Miner's Primary Tool",
		"Botanist's Primary Tool",
		"Fisher's Primary Tool",
		"Carpenter's Secondary Tool",
		"Blacksmith's Secondary Tool",
		"Armorer's Secondary Tool",
		"Goldsmith's Secondary Tool",
		"Leatherworker's Secondary Tool",
		"Weaver's Secondary Tool",
		"Alchemist's Secondary Tool",
		"Culinarian's Secondary Tool",
		"Miner's Secondary Tool",
		"Botanist's Secondary Tool",
		"Fisher's Secondary Tool",
		"Gladiator's Arm",
		"Marauder's Arm",
		"Dark Knight's Arm",
		"Gunbreaker's Arm",
		"Lancer's Arm",
		"Pugilist's Arm",
		"Samurai's Arm",
		"Rogue's Arm",
		"Archer's Arm",
		"Machinist's Arm",
		"Dancer's Arm",
		"One-handed Thaumaturge's Arm",
		"Two-handed Thaumaturge's Arm",
		"Arcanist's Grimoire",
		"Red Mage's Arm",
		"Blue Mage's Arm",
		"One-handed Conjurer's Arm",
		"Two-handed Conjurer's Arm",
		"Scholar's Arm",
		"Astrologian's Arm"} {
		if i.ItemCategory == gearCategory {
			return true
		}
	}
	return false
}

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
	lodestone_session_token := viper.Get("lodestone_session_token")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s/%s", character_id, endpoint), nil)
	if err != nil {
		panic(err)
	}

	if lodestone_session_token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("ldst_sess=%s;", lodestone_session_token))
	}

	return req
}

// Sets up a request
func setupMobileRequest(endpoint string, character_id string) *http.Request {
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
