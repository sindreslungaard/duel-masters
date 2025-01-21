package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MissileSoldierUltimo ...
func MissileSoldierUltimo(c *match.Card) {

	c.Name = "Missile Soldier Ultimo"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, func(card *match.Card, ctx *match.Context) {
			c.AddUniqueSourceCondition(cnd.PowerAttacker, 4000, card.ID)
			c.AddUniqueSourceCondition(cnd.AttackUntapped, nil, card.ID)
		}),
	)

}

// SlaphappySoldierGalback ...
func SlaphappySoldierGalback(c *match.Card) {

	c.Name = "Slaphappy Soldier Galback"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, fx.WheneverThisAttacks(fx.DestroyOpCreatureXPowerOrLess(4000, match.DestroyedByMiscAbility))),
	)

}
