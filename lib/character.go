package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
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
		req := lodestone.SetupMobileRequest(fmt.Sprintf("achievement/?page=%d", page), c.Id)

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
		req := lodestone.SetupRequest(fmt.Sprintf("goldsaucer/tripletriad/?hold=1&page=%d", page), c.Id)

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

	req := lodestone.SetupRequest("class_job", c.Id)

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

	req := lodestone.SetupMobileRequest("minion", c.Id)

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
	req := lodestone.SetupMobileRequest("mount", c.Id)
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
	req := lodestone.SetupMobileRequest("orchestrion", c.Id)
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
	req := lodestone.SetupRequest("retainer", c.Id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		cobra.CheckErr("There was an error retrieiving data from Lodestone. Is your session token correct?")
	}

	if resp.StatusCode == 404 {
		lodestone.GetSessionToken()
		req := lodestone.SetupRequest("retainer", c.Id)
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
	req := lodestone.SetupMobileRequest("bluemage", c.Id)
	var spells []Spell

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		lodestone.GetSessionToken()
		req = lodestone.SetupMobileRequest("bluemage", c.Id)
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
