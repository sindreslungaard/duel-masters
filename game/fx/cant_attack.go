package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// CantAttackPlayers prevents a card from attacking players
func CantAttackPlayers(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
	}

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackPlayers) {
			return
		}

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))

		ctx.InterruptFlow()

	}

}

// CantAttackPlayers prevents a card from attacking players
func CantAttackPlayersIfNoShields(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantAttackPlayersIfNoShields, true, card.ID)
	}

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackPlayersIfNoShields) {
			return
		}

		if len(Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)) == 0 {
		    ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players if they have no shields", card.Name))
			ctx.InterruptFlow()
		}

	}
}


// CantAttackPlayers prevents a card from attacking players
func CantAttackCreaturesIfNoShields(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantAttackCreaturesIfNoShields, true, card.ID)
	}

	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackCreaturesIfNoShields) {
			return
		}

		if len(Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)) == 0 {
		    ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures if they have no shields", card.Name))
			ctx.InterruptFlow()
		}

	}
}

// CantAttackCreatures prevents a card from attacking players
func CantAttackCreatures(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
	}

	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackCreatures) {
			return
		}

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))

		ctx.InterruptFlow()

	}

}
