package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GachackMechanicalDoll ...
func GachackMechanicalDoll(c *match.Card) {

	c.Name = "Gachack, Mechanical Doll"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.DeathPuppet}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush },
			fx.When(fx.WheneverThisAttacksAndIsntBlocked, fx.DestroyOpponentCreature(true, match.DestroyedByMiscAbility))),
	)

}
