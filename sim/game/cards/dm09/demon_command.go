package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StallobTheLifequasher ...
func StallobTheLifequasher(c *match.Card) {

	c.Name = "Stallob, the Lifequasher"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})

		fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}
