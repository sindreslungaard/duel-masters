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

func TrenchdiveShark(c *match.Card) {

	c.Name = "Trenchdive Shark"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(c *match.Card, ctx *match.Context) {
		fx.RotateShields(c, ctx, 2)
	}))
}
