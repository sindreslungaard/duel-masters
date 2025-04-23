package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RocketdiveSkyterror ...
func RocketdiveSkyterror(c *match.Card) {

	c.Name = "Rocketdive Skyterror"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.CantBeAttacked, fx.CantAttackPlayers, fx.PowerAttacker1000)
}

// TorpedoSkyterror ...
func TorpedoSkyterror(c *match.Card) {

	c.Name = "Torpedo Skyterror"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if !attacking {
			return 0
		}

		return len(
			fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.Tapped && x.ID != c.ID
				},
			),
		) * 2000
	}

	c.Use(fx.Creature)
}
