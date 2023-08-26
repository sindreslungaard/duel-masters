package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Tropico ...
func Tropico(c *match.Card) {

	c.Name = "Tropico"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			creatures, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			if len(creatures) < 3 {
				return
			}

			ctx.ScheduleAfter(func() {
				event.Blockers = make([]*match.Card, 0)
			})

		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			creatures, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			if len(creatures) < 3 {
				return
			}

			ctx.ScheduleAfter(func() {
				event.Blockers = make([]*match.Card, 0)
			})

		}

	}, fx.Creature)

}
