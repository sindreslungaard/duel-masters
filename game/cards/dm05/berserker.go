package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BallusDogfightEnforcerQ ...
func BallusDogfightEnforcerQ(c *match.Card) {

	c.Name = "Ballus, Dogfight Enforcer Q"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Berserker
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
		).Map(func(x *match.Card) {
			card.Tapped = false
		})

	}))

}
