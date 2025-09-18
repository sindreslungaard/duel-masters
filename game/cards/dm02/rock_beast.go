package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Galsaur ...
func Galsaur(c *match.Card) {

	c.Name = "Galsaur"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if (len(fx.Find(c.Player, match.BATTLEZONE)) == 1) && attacking {
			return 4000
		}

		return 0

	}

	c.Use(fx.Creature, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		if len(fx.Find(card.Player, match.BATTLEZONE)) == 1 {
			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
		}

	}))

}

// Bombersaur ...
func Bombersaur(c *match.Card) {

	c.Name = "Bombersaur"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Destroyed, fx.EachPlayerDestroys2Mana))

}
