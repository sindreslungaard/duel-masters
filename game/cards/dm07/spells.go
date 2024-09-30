package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ApocalypseVise ...
func ApocalypseVise(c *match.Card) {

	c.Name = "Apocalypse Vise"
	c.Civ = civ.Fire
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		correctSelection := false

		for !correctSelection {

			creatures := fx.Select(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Destroy any number of your opponent's creatures that have total power 8000 or less.",
				0,
				16,
				false,
			)

			totalPower := 0

			for _, creature := range creatures {
				totalPower += ctx.Match.GetPower(creature, false)
			}

			if totalPower <= 8000 {

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
				}

				correctSelection = true

			}

		}

	}))

}
