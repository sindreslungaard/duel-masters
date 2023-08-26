package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BrawlerZyler ...
func BrawlerZyler(c *match.Card) {

	c.Name = "Brawler Zyler"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// FatalAttackerHorvath ...
func FatalAttackerHorvath(c *match.Card) {

	c.Name = "Fatal Attacker Horvath"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if attacking && match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.HasFamily(family.Armorloid) }) {
			return 2000
		}

		return 0

	}

}

// ImmortalBaronVorg ...
func ImmortalBaronVorg(c *match.Card) {

	c.Name = "Immortal Baron, Vorg"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}
