package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DeathbladeBeetle ...
func DeathbladeBeetle(c *match.Card) {

	c.Name = "Deathblade Beetle"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker4000)

}

// ForestHornet ...
func ForestHornet(c *match.Card) {

	c.Name = "Forest Hornet"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}

// RedEyeScorpion ...
func RedEyeScorpion(c *match.Card) {

	c.Name = "Red-Eye Scorpion"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ReturnToMana)

}
