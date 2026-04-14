package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MiniTitanGett ...
func MiniTitanGett(c *match.Card) {

	c.Name = "Mini Titan Gett"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack, fx.PowerAttacker1000)

}

// ArmoredCannonBalbaro ...
func ArmoredCannonBalbaro(c *match.Card) {

	c.Name = "Armored Cannon Balbaro"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking {

			fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasFamily(family.Human) && x != c },
			).Map(func(x *match.Card) {
				power += 2000
			})

		}

		return power

	}

	c.Use(fx.Creature, fx.Evolution)

}

// ArmoredBlasterValdios ...
func ArmoredBlasterValdios(c *match.Card) {

	c.Name = "Armored Blaster Valdios"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.ModifyPowers(func(event *match.GetPowerEvent) {

		if c.Zone != match.BATTLEZONE {
			return
		}

		if event.Card.Player != c.Player {
			return
		}

		if event.Card.ID == c.ID {
			return
		}

		if !event.Card.HasFamily(family.Human) {
			return
		}

		event.Power += 1000

	}))

}
