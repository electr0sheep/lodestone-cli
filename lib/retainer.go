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
		itemInfo := retainerItemElement.Find(".db-tooltip__item-info__list")
		itemInfo.Each(func(_ int, item *goquery.Selection) {

		})
		gearInfo := retainerItemElement.Find(".db-tooltip__item-info__list")
		extractable := false
		projectable := false
		desynthesizable := ""
		dyable := false
		if gearInfo.Nodes != nil {
			extractable = gearInfo.Nodes[0].FirstChild.LastChild.FirstChild.Data == "Yes"
			projectable = gearInfo.Nodes[0].FirstChild.NextSibling.LastChild.FirstChild.Data == "Yes"
			desynthesizable = gearInfo.Nodes[0].FirstChild.NextSibling.NextSibling.LastChild.FirstChild.Data
			dyable = gearInfo.Nodes[0].FirstChild.NextSibling.NextSibling.NextSibling.LastChild.FirstChild.Data == "Yes"
		}
		// can leave the && out if you want to know if it's purchasable with any currency
		// however, for decluttering I just wanna know if I can buy it with gil
		purchasable := retainerItemElement.Find(".db-view__sells").Nodes[0].NextSibling.FirstChild != nil && strings.Contains(retainerItemElement.Find(".db-view__sells").Nodes[0].NextSibling.FirstChild.Data, "gil")
		if highQuality {
			name = strings.TrimSuffix(name, "")
		}

		r.Items = append(r.Items, Item{
			Name:                 name,
			Quantity:             quantity,
			HighQuality:          highQuality,
			CanBePlacedInArmoire: canBePlacedInArmoire,
			IsUnique:             isUnique,
			ItemCategory:         itemCategory,
			Extractable:          extractable,
			Projectable:          projectable,
			Desynthesizable:      desynthesizable,
			Dyable:               dyable,
			Purchasable:          purchasable})
	})
}
