package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AncientGiant ...
func AncientGiant(c *match.Card) {

	c.Name = "Ancient Giant"
	c.Power = 9000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {

				blockers := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {
					if blocker.Civ != civ.Darkness {
						blockers = append(blockers, blocker)
					}
				}

				event.Blockers = blockers

			})

		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {

				blockers := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {
					if blocker.Civ != civ.Darkness {
						blockers = append(blockers, blocker)
					}
				}

				event.Blockers = blockers

			})
		}

	}, fx.Creature, fx.Doublebreaker)
}
