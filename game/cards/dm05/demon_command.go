package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DeathCruzerTheAnnihilator ...
func DeathCruzerTheAnnihilator(c *match.Card) {

	c.Name = "Death Cruzer, the Annihilator"
	c.Power = 13000
	c.Civ = civ.Darkness
	c.Family = family.DemonCommand
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {

			if x.ID != card.ID {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			}

		})

	}))

}
