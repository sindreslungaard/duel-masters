package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantBeBlocked allows the card to attack without being blocked
func CantBeBlocked(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeBlocked, nil, card.ID)

	}

}

// CantBeBlockedIf
func CantBeBlockedIf(test func(blocker *match.Card) bool) func(*match.Card, *match.Context) {

	return func(card *match.Card, ctx *match.Context) {

		filter := func(blockers []*match.Card) []*match.Card {
			filtered := []*match.Card{}

			for _, b := range blockers {
				if !test(b) {
					filtered = append(filtered, b)
				}
			}

			return filtered
		}

		switch event := ctx.Event.(type) {

		case *match.AttackCreature:
			if event.CardID != card.ID {
				return
			}
			event.Blockers = filter(event.Blockers)
		case *match.AttackPlayer:
			if event.CardID != card.ID {
				return
			}
			event.Blockers = filter(event.Blockers)
		default:
			return
		}

	}

}
