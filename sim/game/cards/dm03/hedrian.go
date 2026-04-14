package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Mudman ...
func Mudman(c *match.Card) {

	c.Name = "Mudman"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Darkness }) {
			return 0
		}

		return 2000
	}

	c.Use(fx.Creature)
}


// Scratchclaw ...
func Scratchclaw(c *match.Card) {

	c.Name = "Scratchclaw"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		return (getDarknessCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	c.Use(fx.Creature, fx.Slayer)
}

// Return the number of darkness creatures in your battle zone
func getDarknessCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Darkness {
			count++
		}
	}

	return count
}