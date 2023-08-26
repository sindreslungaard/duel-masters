package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ChiliasTheOracle ...
func ChiliasTheOracle(c *match.Card) {

	c.Name = "Chilias, the Oracle"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ReturnToHand)

}

// IocantTheOracle ...
func IocantTheOracle(c *match.Card) {

	c.Name = "Iocant, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.HasFamily(family.AngelCommand) }) {
			power += 2000
		}

		return power
	}

}

// ReusolTheOracle ...
func ReusolTheOracle(c *match.Card) {

	c.Name = "Reusol, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}
