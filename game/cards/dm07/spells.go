package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func HopelessVortex(c *match.Card) {
	c.Name = "Hopeless Vortex"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Destroy one of your opponent's creatures", 1, 1, false)

		for _, creature := range creatures {

			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

		}
	}))
}
