package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SnipStrikerBullraizer ...
func SnipStrikerBullraizer(c *match.Card) {

	c.Name = "Snip Striker Bullraizer"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.UntapStep); ok {

			creatures := fx.Find(card.Player, match.BATTLEZONE)

			oppCreatures := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

			if len(creatures) < len(oppCreatures) {
				card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
				card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
			} else {
				card.RemoveCondition(cnd.CantAttackPlayers)
				card.RemoveCondition(cnd.CantAttackCreatures)
			}

		}

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.From == match.BATTLEZONE || event.To == match.BATTLEZONE {

				creatures := fx.Find(card.Player, match.BATTLEZONE)

				oppCreatures := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

				if len(creatures) < len(oppCreatures) {
					card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
					card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
				} else {
					card.RemoveCondition(cnd.CantAttackPlayers)
					card.RemoveCondition(cnd.CantAttackCreatures)
				}

			}

		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			// Is this event for me or someone else?
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackPlayers) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))

			ctx.InterruptFlow()

		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			// Is this event for me or someone else?
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackCreatures) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))

			ctx.InterruptFlow()

		}
	})

}
