package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MaskedPomegranate ...
func MaskedPomegranate(c *match.Card) {

	c.Name = "Masked Pomegranate"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	// Get +1000 power for each other nature card in your battle zone
	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return (getNatureCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	c.Use(fx.Creature, fx.CantBeBlockedByPowerUpTo4000)

}

// Return the number of nature creatures in your battle zone
func getNatureCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Nature {
			count++
		}
	}

	return count
}
