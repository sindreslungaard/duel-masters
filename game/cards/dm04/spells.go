package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FullDefensor ...
func FullDefensor(c *match.Card) {

	c.Name = "Full Defensor"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			ctx.Match.ApplyPersistentEffect(func(_ *match.Card, ctx2 *match.Context, exit func()) {

				// on all events, add blocker to our creatures
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
				})

				// remove persistent effect on start of next turn
				_, ok := ctx2.Event.(*match.StartOfTurnStep)
				if ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
				}

			})

		}
	})
}
