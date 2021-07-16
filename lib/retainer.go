package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/electr0sheep/lodestone-cli/lodestone"
	"github.com/spf13/cobra"
)

type Retainer struct {
	Name    string
	OwnerId string
	Id      string
	Items   []Item
}

func (r *Retainer) GetItems() {
	client := &http.Client{}
	req := lodestone.SetupRequest(fmt.Sprintf("retainer/%s/baggage", r.Id), r.OwnerId)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		cobra.CheckErr("There was an error retrieiving data from Lodestone. Is your session token correct?")
	}

	if resp.StatusCode == 404 {
		lodestone.GetSessionToken()
		req := lodestone.SetupRequest(fmt.Sprintf("retainer/%s/baggage", r.Id), r.OwnerId)
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
