package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GeneralDarkFiend ...
func GeneralDarkFiend(c *match.Card) {

	c.Name = "General Dark Fiend"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		shields := fx.Find(
			card.Player,
			match.SHIELDZONE,
		)

		if len(shields) < 1 {
			return
		}

		ctx.Match.MoveCard(shields[0], match.GRAVEYARD, card)

	}))

}
