package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BaleskBajTheTimeburner ...
func BaleskBajTheTimeburner(c *match.Card) {

	c.Name = "Balesk Baj, the Timeburner"
	c.Power = 8000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Fire}

	attacked := false

	c.Use(fx.Creature, fx.Doublebreaker, fx.Evolution,
		func(card *match.Card, ctx *match.Context) {
			if _, ok := ctx.Event.(*match.UntapStep); ok {
				attacked = false
			}
		},
		fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, func(card *match.Card, ctx *match.Context) {
			if card.Zone != match.BATTLEZONE {
				attacked = false
				return
			}

			attacked = true
		}),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
			if card.Zone != match.BATTLEZONE {
				attacked = false
				return
			}

			ctx.ScheduleAfter(func() {
				_, err := card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)

				if err == nil {
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the %s's hand", card.Name, card.Player.Username()))
				} else {
					attacked = false
				}
			})
		}),
		func(card *match.Card, ctx *match.Context) {
			if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
				ctx.ScheduleAfter(func() {
					if attacked {
						attacked = false
						ctx.InterruptFlow()
						ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: You will take an extra turn after this one.", card.Name))
						ctx.Match.BeginNewTurn(true)
					}
				})
			}
		},
	)

}
