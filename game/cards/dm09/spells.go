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
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("Until the start of your next turn, each of your '%s' creatures in the battlezone gets 'Blocker'", family))

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family)
					},
				).Map(func(x *match.Card) {
					_, ok := ctx2.Event.(*match.StartOfTurnStep)
					if ok && ctx2.Match.IsPlayerTurn(card.Player) {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					fx.ForceBlocker(x, ctx2, card.ID)
				})
			})
		}
	}))

}
