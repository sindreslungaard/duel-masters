package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonCusdalf ...
func TerradragonCusdalf(c *match.Card) {

	c.Name = "Terradragon Cusdalf"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker4000, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.UntapManaEvent); ok && card.Zone == match.BATTLEZONE {
			if event.CurrentPlayer == card.Player {
				ctx.InterruptFlow()
			}
		}
	})

}
