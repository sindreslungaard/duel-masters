package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GrapeGlobbo ...
func GrapeGlobbo(c *match.Card) {

	c.Name = "Grape Globbo"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		opponentHandIds := make([]string, 0)

		fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE).Map(func(x *match.Card) {
			opponentHandIds = append(opponentHandIds, x.ImageID)
		})

		ctx.Match.ShowCards(card.Player, "Your opponent's hand:", opponentHandIds)

	}))
}
