package item

import "strings"

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
