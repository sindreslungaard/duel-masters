package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ScissorEye ...
func ScissorEye(c *match.Card) {

	c.Name = "Scissor Eye"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)

}

// PlasmaChaser ...
func PlasmaChaser(c *match.Card) {

	c.Name = "Plasma Chaser"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		toDraw := len(fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE))

		if toDraw < 1 {
			return
		}

		fx.MayDrawAmount(card, ctx, toDraw)

	}))

}
