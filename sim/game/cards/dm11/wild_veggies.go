package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RoyalDurian ...
func RoyalDurian(c *match.Card) {

	c.Name = "Royal Durian"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.SilentSkill(dragonFromMana()))

}

// MachoMelon ...
func MachoMelon(c *match.Card) {

	c.Name = "Macho Melon"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.WaveStriker, fx.WaveStrikerGrant(cnd.PowerAttacker, 3000))

}

// NinjaPumpkin ...
func NinjaPumpkin(c *match.Card) {

	c.Name = "Ninja Pumpkin"
	c.Power = 2000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = fx.WaveStrikerPower(c, 4000)

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.CantBeBlockedByPowerUpTo5000))

}
