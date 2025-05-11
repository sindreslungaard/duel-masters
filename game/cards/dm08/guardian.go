package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SolGallaHaloGuardian ...
func SolGallaHaloGuardian(c *match.Card) {

	c.Name = "Sol Galla, Halo Guardian"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	powerBoost := 0

	c.PowerModifier = func(m *match.Match, attacking bool) int { return powerBoost }

	c.Use(fx.Creature, fx.Blocker(),
		fx.When(fx.AnySpellCast, func(card *match.Card, ctx *match.Context) { powerBoost += 3000 }),
		fx.When(fx.EndOfTurn, func(card *match.Card, ctx *match.Context) { powerBoost = 0 }),
	)
}
