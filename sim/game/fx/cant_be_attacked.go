package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantBeAttacked
func CantBeAttacked(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeAttacked, nil, card.ID)

	}

}

// CantBeAttackedIf
func CantBeAttackedIf(test func(attacker *match.Card) bool) func(*match.Card, *match.Context) {

	return func(card *match.Card, ctx *match.Context) {

		switch event := ctx.Event.(type) {

		case *match.AttackCreature:

			// am i in the list of attackable cards?
			me := false
			for _, creature := range event.AttackableCreatures {
				if creature.ID == card.ID {
					me = true
				}
			}

			if !me {
				return
			}

			// find attacking creature ref
			attackingCreature, err := ctx.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			if !test(attackingCreature) {
				return
			}

			// i need to be removed from list of attackable creatures
			filtered := []*match.Card{}
			for _, creature := range event.AttackableCreatures {
				if creature.ID != card.ID {
					filtered = append(filtered, creature)
				}
			}

			event.AttackableCreatures = filtered

		default:
			return
		}

	}

}
