package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RagingDashHorn ...
func RagingDashHorn(c *match.Card) {

	c.Name = "Raging Dash-Horn"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Nature }) {
			return 0
		}

		return 3000
	}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		if match.ContainerHas(card.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Nature }) {
			card.RemoveCondition(cnd.DoubleBreaker)
		} else {
			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
		}

	}))

}
