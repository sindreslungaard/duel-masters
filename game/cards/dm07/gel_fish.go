package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SplashZebrafish(c *match.Card) {

	c.Name = "Splash Zebrafish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.CantBeBlocked, fx.When(fx.Summoned, fx.ReturnMyCardFromMZToHand))
}
