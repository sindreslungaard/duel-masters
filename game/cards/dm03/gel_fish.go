package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ChaosFish ...
func ChaosFish(c *match.Card) {

	c.Name = "Chaos Fish"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	// Get +1000 power for each other water card in your battle zone
	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return (getWaterCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	//When this creature attacks, draw as many cards as other water creatures you have in the battle zone
	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			nrCardsToDraw := (getWaterCardsInYourBattleZone(card) - 1) //-1 to exclude self
			fx.DrawBetween(card, ctx, 0, nrCardsToDraw)
		})
	}))

}

// Return the number of water creatures in your battle zone
func getWaterCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Water {
			count++
		}
	}

	return count
}
