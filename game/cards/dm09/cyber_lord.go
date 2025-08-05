package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Hokira ...
func Hokira(c *match.Card) {

	c.Name = "Hokira"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = hokiraTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func hokiraTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(
		card,
		ctx,
		fmt.Sprintf("%s's effect: Choose a race. Whenever one of your creatures of that race would be destroyed this turn, return it to your hand instead.", card.Name),
	)

	if family != "" {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("Whenever one of your '%s' creatures would be destroyed this turn, return it to your hand instead.", family))

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// Because at the end of the turn, creatures might still be destroyed
			// by some effects
			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				ctx2.ScheduleAfter(func() {
					exit()
				})
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family)
				},
			).Map(func(x *match.Card) {
				if fx.WouldBeDestroyed(x, ctx2) {
					x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
					ctx2.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was returned to hand instead of being destroyed due to %s's effect.", x.Name, card.Name))
				}
			})
		})
	}

}
