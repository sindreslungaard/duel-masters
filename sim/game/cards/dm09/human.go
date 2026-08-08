package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ShockTrooperMykee ...
func ShockTrooperMykee(c *match.Card) {

	c.Name = "Shock Trooper Mykee"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(
		fx.Creature,
		fx.SpeedAttacker,
		fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.DestroyOpCreatureXPowerOrLess(3000, true, match.DestroyedByMiscAbility)),
	)

}
