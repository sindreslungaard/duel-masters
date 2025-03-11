package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AmnisHolyElemental ...
func AmnisHolyElemental(c *match.Card) {

	c.Name = "Amnis, Holy Elemental"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature,
		fx.DarknessBlocker,
		func(card *match.Card, ctx *match.Context) {

			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card == card {
				if event.Context == match.DestroyedInBattle && event.Source.Civ == civ.Darkness {
					ctx.InterruptFlow()
				}
			}

		})

}
