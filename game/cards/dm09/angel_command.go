package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MihailCelestialElemental ...
func MihailCelestialElemental(c *match.Card) {

	c.Name = "Mihail, Celestial Elemental"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && card.Zone == match.BATTLEZONE {
			if event.Card.ID != card.ID {
				ctx.InterruptFlow()
				ctx.Match.ReportActionInChat(event.Card.Player, fmt.Sprintf("%s stayed in the battlezone instead of being destroyed due to %s's effect.", event.Card.Name, card.Name))
			}
		}
	})

}
