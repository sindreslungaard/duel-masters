package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EvisceratingWarriorLumez ...
func EvisceratingWarriorLumez(c *match.Card) {

	c.Name = "Eviscerating Warrior Lumez"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Armorloid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	// At 2000 power it is caught by its own sweep, which is the printed cost of
	// the effect rather than an oversight.
	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned,
		fx.DestroyAllCreaturesXPowerOrLess(2000, match.DestroyedByMiscAbility))))

}
