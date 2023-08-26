package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AlekSolidityEnforcer ...
func AlekSolidityEnforcer(c *match.Card) {

	c.Name = "Alek, Solidity Enforcer"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	// Get +1000 power for each other light card in your battle zone
	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return (getLightCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	c.Use(fx.Creature, fx.Blocker)
}

// Return the number of water creatures in your battle zone
func getLightCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Light {
			count++
		}
	}

	return count
}
