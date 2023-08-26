package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BurningMane ...
func BurningMane(c *match.Card) {

	c.Name = "Burning Mane"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}

// FearFang ...
func FearFang(c *match.Card) {

	c.Name = "Fear Fang"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}

// SteelSmasher ...
func SteelSmasher(c *match.Card) {

	c.Name = "Steel Smasher"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantAttackPlayers)

}

// BronzeArmTribe ...
func BronzeArmTribe(c *match.Card) {

	c.Name = "Bronze-Arm Tribe"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.DrawToMana)

}

// GoldenWingStriker ...
func GoldenWingStriker(c *match.Card) {

	c.Name = "Golden Wing Striker"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// MightyShouter ...
func MightyShouter(c *match.Card) {

	c.Name = "Mighty Shouter"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ReturnToMana)

}
