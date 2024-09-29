package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SiriGloryElemental(c *match.Card) {

	c.Name = "Siri, Glory Elemental"
	c.Power = 7000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker,
		fx.WhenAll([]func(*match.Card, *match.Context) bool{fx.EndOfMyTurnCreatureBZ, fx.IDontHaveShields}, fx.MayUntapSelf),
		fx.When(fx.InTheBattlezone, fx.BlockerWhenNoShields),
	)

}
