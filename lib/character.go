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
	Currencies       []Currency
	GoldSaucer       GoldSaucer
	Jobs             []Job
	Minions          []Minion
	Mounts           []Mount
	Orchestrions     []Orchestrion
	Reputations      []Reputation
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
		achievementElements := doc.Find(".entry__achievement")

		if achievementElements.Length() == 0 {
			morePages = false
		} else {
			achievementElements.Each(func(_ int, achievementElement *goquery.Selection) {
				name := achievementElement.Find(".entry__activity__txt").Text()
				name = strings.Split(name, "\"")[1]
				acquistionDate := achievementElement.Find(".entry__activity__time").Text()
				if strings.Contains(acquistionDate, "ldst_strftime(") {
					// need to convert epoch to date
					epoch := strings.Split(acquistionDate, "ldst_strftime(")[1]
					epoch = strings.Split(epoch, ",")[0]
					timestamp, _ := strconv.ParseInt(epoch, 10, 64)
					myDate := time.Unix(timestamp, 0)
					acquistionDate = fmt.Sprintf("%d/%d/%d", myDate.Month(), myDate.Day(), myDate.Year())
				}
				achievements = append(achievements, Achievement{AcquistionDate: acquistionDate, Name: name})
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
		cardElements := doc.Find(".tripletriad-card_list li")

		if cardElements.Length() == 0 {
			morePages = false
		} else {
			cardElements.Each(func(_ int, cardElement *goquery.Selection) {
				cardAttackDown, _ := cardElement.Find(".strength .down").Attr("alt")
				cardAttackLeft, _ := cardElement.Find(".strength .left").Attr("alt")
				cardAttackRight, _ := cardElement.Find(".strength .right").Attr("alt")
				cardAttackUp, _ := cardElement.Find(".strength .up").Attr("alt")
				cardName := cardElement.Find(".name_inner").Text()
				_, cardNotAcquired := cardElement.Attr("class")
				cardRarity := len(cardElement.Find(".tripletriad-tooltip__card .rarity img").Nodes)
				var cardType string
				cardTypeElements := cardElement.Find(".tripletriad-tooltip__text .type span")
				if len(cardTypeElements.Nodes) == 1 {
					cardType = ""
				} else {
					cardType = cardTypeElements.Nodes[1].FirstChild.Data
				}
				cardDescription := cardElement.Find(".flavor").Text()
				cards = append(cards, Card{Acquired: !cardNotAcquired, AttackDown: cardAttackDown, AttackLeft: cardAttackLeft, AttackRight: cardAttackRight, AttackUp: cardAttackUp, Description: cardDescription, Name: cardName, Rarity: cardRarity, Type: cardType})
			})
		}
	}

	c.Cards = cards
}

func (c *Character) GetCurrenciesAndRep() {
	client := &http.Client{}
	var currencies []Currency
	var reputations []Reputation

	req := lodestone.SetupRequest("currency", c.Id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	currencyElements := doc.Find(".character__currency__list li")
	// this is kinda weird, but Lodestone doesn't group the heading--lead classes
	// consistenly
	reputationElements := doc.Find(".character__reputation .heading--lead")

	currencyElements.Each(func(_ int, currencyElement *goquery.Selection) {
		currencyType := currencyElement.Find(".heading--lead").Text()
		currencyElement.Find(".currency__box").Each(func(_ int, individualElement *goquery.Selection) {
			currentAmount := ""
			maximum := ""
			thing := individualElement.Find(".currency__box__text__name")
			name := thing.Text()
			if len(thing.Nodes) > 0 {
				amount := thing.Nodes[0].Parent.LastChild.Data
				splitString := strings.Split(amount, "/")
				if len(splitString) > 1 {
					currentAmount = strings.TrimSpace(strings.Split(amount, "/")[0])
					maximum = strings.TrimSpace(strings.Split(amount, "/")[1])
				} else {
					currentAmount = strings.TrimSpace(strings.Split(amount, "/")[0])
				}
			}
			if currentAmount == "" {
				currentAmount = strings.TrimSpace(currencyElement.Find(".currency__box__text").Text())
			}
			currencies = append(currencies, Currency{CurrentAmount: currentAmount, Maximum: maximum, Name: name, Type: currencyType})
		})
	})

	c.Currencies = currencies

	reputationElements.Each(func(index int, reputationElement *goquery.Selection) {
		reputationType := reputationElement.Text()
		// Player Comms
		if index == 0 {
			commNode := reputationElement.Parent().Find(".character-block__box")
			name := commNode.Find(".character-block__name").Text()
			currentAmount := commNode.Find(".character-block__value").Text()
			maximum := ""
			reputationLevel := ""
			reputations = append(reputations, Reputation{CurrentAmount: currentAmount, Maximum: maximum, Name: name, ReputationLevel: reputationLevel, Type: reputationType})
			// Beast Tribe
		} else if index == 1 {
			repNodeList := goquery.NewDocumentFromNode(reputationElement.Nodes[0].NextSibling.NextSibling)
			repNodes := repNodeList.Find(".character-block__box--beast_tribe")
			repNodes.Each(func(_ int, beastTribe *goquery.Selection) {
				name := beastTribe.Find(".character-block__box--beast_tribe__name").Text()
				splitString := strings.Split(beastTribe.Find(".character-block__point").Text(), "/")
				currentAmount := strings.TrimSpace(splitString[0])
				maximum := strings.TrimSpace(splitString[1])
				reputationLevel := beastTribe.Find(".character-block__friendship").Text()
				reputations = append(reputations, Reputation{CurrentAmount: currentAmount, Maximum: maximum, Name: name, ReputationLevel: reputationLevel, Type: reputationType})
			})
		}
	})

	c.Reputations = reputations
}

func (c *Character) GetGoldSaucer() {
	client := &http.Client{}

	req := lodestone.SetupRequest("goldsaucer", c.Id)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	goldSaucerElement := doc.Find(".character__content")
	c.GoldSaucer.MGP = goldSaucerElement.Find(".character__currency__list p").Text()
	c.GoldSaucer.TripleTriadTournamentResult = goldSaucerElement.Find(".character__goldsaucer__text").Nodes[0].FirstChild.Data
	jumboCactpot := goquery.NewDocumentFromNode(goldSaucerElement.Find(".character__goldsaucer__text").Nodes[1])
	c.GoldSaucer.JumboCactpot = strings.ReplaceAll(strings.ReplaceAll(jumboCactpot.Text(), "\t", ""), "\n", "")
	numbers := jumboCactpot.Next().Text()
	numbers = strings.Split(numbers, ": ")[1]
	c.GoldSaucer.JumboCactpotNumberOne = strings.Split(numbers, ", ")[0]
	c.GoldSaucer.JumboCactpotNumberTwo = strings.Split(numbers, ", ")[1]
	c.GoldSaucer.JumboCactpotNumberThree = strings.Split(numbers, ", ")[2]
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
	m.Movement = doc.Find(".mount__type span").Text()
	m.Description = doc.Find(".mount__text").Text()

	if len(strings.Split(m.Movement, "x")) > 1 {
		m.MaxRiders = strings.Split(m.Movement, "x")[1]
		m.Movement = strings.Split(m.Movement, "x")[0]
	} else {
		m.MaxRiders = "1"
	}
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
	orchestrionCategories := doc.Find(".orchestrion-category")

	orchestrionCategories.Each(func(_ int, orchestrionCategory *goquery.Selection) {
		category := orchestrionCategory.Find(".orchestrion-title").Children().Last().Text()
		orchestrionElements := orchestrionCategory.Find(".orchestrion-list").Find("li")

		orchestrionElements.Each(func(_ int, orchestrionElement *goquery.Selection) {
			_, unacquired := orchestrionElement.Attr("class")
			name := orchestrionElement.Find(".orchestrion-list__name").Text()
			name = strings.ReplaceAll(name, "\t", "")
			name = strings.ReplaceAll(name, "\n", "")
			// We need to massage the data a little bit. The names of the orchestrions
			// in ffxivcollect are titles (i.e. The Maiden's Lament as opposed to The maiden's lament)
			whereToFind := orchestrionElement.Find(".orchestrion-detail__text").Text()
			orchestrions = append(orchestrions, Orchestrion{Acquired: !unacquired, Category: category, Name: name, WhereToFind: whereToFind})
		})
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
