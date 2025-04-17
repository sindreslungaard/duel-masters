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
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// Because at the end of the turn, creatures might still be destroyed
			// by some effects
			if _, ok := ctx2.Event.(*match.BeginTurnStep); ok {
				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family) && !x.HasCondition(card.ID+"-custom")
				},
			).Map(func(x *match.Card) {
				ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit2 func()) {

					if _, ok := ctx3.Event.(*match.BeginTurnStep); ok {
						x.RemoveConditionBySource(card.ID)
						exit2()
						return
					}

					x.AddUniqueSourceCondition(card.ID+"-custom", true, card.ID)

					if fx.WouldBeDestroyed(x, ctx3) {
						ctx3.InterruptFlow()
						x.RemoveConditionBySource(card.ID)
						x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
						ctx3.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was returned to hand instead of being destroyed due to %s's effect.", x.Name, card.Name))
						exit2()
					}

				})
			})
		})
	}

}
