package lib

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/electr0sheep/lodestone-cli/lodestone"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

type Character struct {
	CityState        string
	Clan             string
	FreeCompany      string
	Gender           string
	GrandCompany     string
	GrandCompanyRank string
	Guardian         string
	Id               string
	Name             string
	Nameday          string
	Race             string
	Title            string
	World            string
	Linkshells       []string
	Achievements     []Achievement
	Cards            []Card
	Jobs             []Job
	Minions          []Minion
	Mounts           []Mount
	Orchestrions     []Orchestrion
	Retainers        []*Retainer
	Spells           []Spell
}

func (c *Character) GetProfile() {
	client := &http.Client{}
	req := lodestone.SetupRequest("", c.Id)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	c.Name = doc.Find(".frame__chara__name").Text()
	c.Title = doc.Find(".frame__chara__title").Text()
	c.World = strings.ReplaceAll(doc.Find(".frame__chara__world").Text(), "\u00a0", " ")
	c.Race = doc.Find(".character-block__name").Nodes[0].FirstChild.Data
	c.Clan = strings.Split(doc.Find(".character-block__name").Nodes[0].FirstChild.NextSibling.NextSibling.Data, " / ")[0]
	c.Gender = strings.Split(doc.Find(".character-block__name").Nodes[0].FirstChild.NextSibling.NextSibling.Data, " / ")[1]
	c.Nameday = doc.Find(".character-block__birth").Text()
	c.Guardian = doc.Find(".character-block__name").Nodes[1].FirstChild.Data
	c.CityState = doc.Find(".character-block__name").Nodes[2].FirstChild.Data
	c.GrandCompany = strings.Split(doc.Find(".character-block__name").Nodes[3].FirstChild.Data, " / ")[0]
	c.GrandCompanyRank = strings.Split(doc.Find(".character-block__name").Nodes[3].FirstChild.Data, " / ")[1]
	c.FreeCompany = doc.Find(".character__freecompany__name a").Text()
	linkshellElements := doc.Find(".character__linkshell__name li")
	linkshellElements.Each(func(_ int, linkshellElement *goquery.Selection) {
		name := linkshellElement.Text()
		c.Linkshells = append(c.Linkshells, name)
	})
}

func (c *Character) GetAchievements() {
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

	c.Achievements = achievements
}

func (c *Character) GetCards() {
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

	c.Cards = cards
}

func (c *Character) GetJobs() {
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
	roleElements := doc.Find(".character__job__role .heading--lead")

	roleElements.Each(func(_ int, roleElement *goquery.Selection) {
		role := roleElement.Text()
		jobElements := roleElement.Next().Find(".character__job li")
		jobElements.Each(func(_ int, jobElement *goquery.Selection) {
			name := jobElement.Find(".character__job__name").Text()
			level := jobElement.Find(".character__job__level").Text()
			xp := jobElement.Find(".character__job__exp").Text()
			jobs = append(jobs, Job{Name: name, Level: level, Role: role, Xp: xp})
		})
	})

	c.Jobs = jobs
}

// getting all minion data at once causes issues
// this allows to lazy load
func (c *Character) GetMinionDetails(m *Minion) {
	client := &http.Client{}
	req := lodestone.SetupRequest(fmt.Sprintf("minion/tooltip/%s", m.Id), c.Id)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	m.AcquistionDate = doc.Find(".minion__header__data").Text()
	if strings.Contains(m.AcquistionDate, "ldst_strftime(") {
		// need to convert epoch to date
		epoch := strings.Split(m.AcquistionDate, "ldst_strftime(")[1]
		epoch = strings.Split(epoch, ",")[0]
		timestamp, _ := strconv.ParseInt(epoch, 10, 64)
		myDate := time.Unix(timestamp, 0)
		m.AcquistionDate = fmt.Sprintf("%d/%d/%d", myDate.Month(), myDate.Day(), myDate.Year())
	}
	m.Behavior = doc.Find(".minion__type span").Text()
	m.Description = doc.Find(".minion__text").Text()
}

func (c *Character) GetMinions() {
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
	minionElements := doc.Find(".minion__list__item")

	minionElements.Each(func(_ int, minionElement *goquery.Selection) {
		tooltipHref, _ := minionElement.Attr("data-tooltip_href")
		id := strings.Split(tooltipHref, "/")[6]
		name := minionElement.Find(".minion__name").Text()
		minions = append(minions, Minion{Id: id, Name: name})
	})

	c.Minions = minions
}

// getting all minion data at once causes issues
// this allows to lazy load
func (c *Character) GetMountDetails(m *Mount) {
	client := &http.Client{}
	req := lodestone.SetupRequest(fmt.Sprintf("mount/tooltip/%s", m.Id), c.Id)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	m.AcquistionDate = doc.Find(".mount__header__data").Text()
	if strings.Contains(m.AcquistionDate, "ldst_strftime(") {
		// need to convert epoch to date
		epoch := strings.Split(m.AcquistionDate, "ldst_strftime(")[1]
		epoch = strings.Split(epoch, ",")[0]
		timestamp, _ := strconv.ParseInt(epoch, 10, 64)
		myDate := time.Unix(timestamp, 0)
		m.AcquistionDate = fmt.Sprintf("%d/%d/%d", myDate.Month(), myDate.Day(), myDate.Year())
	}
	fmt.Println(doc.Find(".mount__text").Html())
	fmt.Println(doc.Find(".mount__text").Text())
	m.Movement = doc.Find(".mount__type span").Text()
	m.Description = doc.Find(".mount__text").Text()
}

func (c *Character) GetMounts() {
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

	mountElements := doc.Find(".mount__list__item")

	mountElements.Each(func(_ int, mountElement *goquery.Selection) {
		tooltipHref, _ := mountElement.Attr("data-tooltip_href")
		id := strings.Split(tooltipHref, "/")[6]
		name := mountElement.Find(".mount__name").Text()
		mounts = append(mounts, Mount{Id: id, Name: name})
	})

	c.Mounts = mounts
}

func (c *Character) GetOrchestrions() {
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

	c.Orchestrions = orchestrions
}

func (c *Character) GetRetainers() {
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

	c.Retainers = retainers
}

func (c *Character) GetSpells() {
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
	spellElements := doc.Find(".sys-reward")

	spellElements.Each(func(_ int, spellElement *goquery.Selection) {
		name := spellElement.Find(".bluemage-action__name").Text()
		name = strings.ReplaceAll(name, "\t", "")
		name = strings.ReplaceAll(name, "\n", "")
		// TODO: Move this crap into a separate massager class?
		name = strings.ReplaceAll(name, " of", " Of")
		name = strings.ReplaceAll(name, " the", " The")
		detail := spellElement.Find(".bluemage-detail__action__type").Text()
		spellType := strings.Split(strings.Split(detail, "\n")[1], ": ")[1]
		aspect := strings.Split(strings.Split(detail, "\n")[2], ": ")[1]
		rank := strings.Split(strings.Split(detail, "\n")[3], ": ")[1]
		description := strings.TrimSpace(spellElement.Find(".bluemage-detail__text").Text())
		description = strings.ReplaceAll(description, "\n", "")
		hint := strings.TrimSpace(spellElement.Find(".bluemage-detail__hint__text").Text())
		spells = append(spells, Spell{Name: name, Type: spellType, Aspect: aspect, Rank: rank, Description: description, WhereToLearn: hint})
	})

	c.Spells = spells
}