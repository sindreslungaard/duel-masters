package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HypersprintWariorUzesol ...
func HypersprintWariorUzesol(c *match.Card) {

	c.Name = "Hypersprint Warior Uzesol"
	c.Power = 1000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.PowerAttacker4000)

}

// WhirlingWarriorMalian ...
func WhirlingWarriorMalian(c *match.Card) {

	c.Name = "Whirling Warrior Malian"
	c.Power = 6000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.TapSelfWhenAnotherCreatureArrives)

}

// FlameTrooperGoliac ...
func FlameTrooperGoliac(c *match.Card) {

	c.Name = "Flame Trooper Goliac"
	c.Power = 4000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Armorloid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned,
		fx.DestroyOpCreatureXPowerOrLess(5000, false, match.DestroyedByMiscAbility))))

}
