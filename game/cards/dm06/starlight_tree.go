package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LightningGrass(c *match.Card) {

	c.Name = "Lightning Grass"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)
}

func RazorpineTree(c *match.Card) {

	c.Name = "Lightning Grass"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		shields, err := c.Player.Container(match.SHIELDZONE)

		if err != nil {
			return 0
		}

		return 2000 * len(shields)
	}

	c.Use(fx.Creature, fx.ShieldTrigger)
}
