package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BerochikaChannelerOfSuns ...
func BerochikaChannelerOfSuns(c *match.Card) {

	c.Name = "Berochika, Channeler of Suns"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		myShields, _ := card.Player.Container(match.SHIELDZONE)

		if len(myShields) >= 5 {
			fx.TopCardToShield(card, ctx)
		}
	}))

}
