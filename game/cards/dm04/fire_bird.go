package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PippieKuppie ...
func PippieKuppie(c *match.Card) {

	c.Name = "Pippie Kuppie"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {
			
			if event.Card.HasFamily(family.ArmoredDragon) {
				event.Power += 1000
			}
		}

	})
}
