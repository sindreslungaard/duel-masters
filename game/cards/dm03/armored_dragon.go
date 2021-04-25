package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GarkagoDragon ...
func GarkagoDragon(c *match.Card) {

	c.Name = "Garkago Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = family.ArmoredDragon
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	// Get +1000 power for each other fire card in your battle zone
	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return (getFireCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	c.Use(fx.Creature, fx.Doublebreaker, fx.AttackUntapped)
}

// BoltailDragon ...
func BoltailDragon(c *match.Card) {

	c.Name = "Boltail Dragon"
	c.Power = 9000
	c.Civ = civ.Fire
	c.Family = family.ArmoredDragon
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker)
}

// Return the number of water creatures in your battle zone
func getFireCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Fire {
			count++
		}
	}

	return count
}
