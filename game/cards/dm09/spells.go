package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// UnifiedResistance ...
func UnifiedResistance(c *match.Card) {

	c.Name = "Unified Resistance"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(
			card,
			ctx,
			fmt.Sprintf("%s's effect: Choose a race. Until the start of your next turn, each of your creatures in the battlezone of that race gets 'Blocker'", card.Name),
		)

		if family != "" {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family) && !x.HasCondition(card.ID+"-custom")
					},
				).Map(func(x *match.Card) {
					ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit2 func()) {
						x.AddUniqueSourceCondition(card.ID+"-custom", true, card.ID)

						_, ok := ctx3.Event.(*match.StartOfTurnStep)
						if ok && ctx3.Match.IsPlayerTurn(card.Player) {
							x.RemoveConditionBySource(card.ID)
							exit2()
							return
						}

						fx.ForceBlocker(x, ctx3, card.ID)
					})
				})

				_, ok := ctx2.Event.(*match.StartOfTurnStep)
				if ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
					return
				}

			})
		}

	}))

}
