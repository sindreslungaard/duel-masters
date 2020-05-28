package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Slayer destroys the source card when the card is destroyed
func Slayer(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Slayer, nil, card.ID)
	}

}

// Suicide destroys the card when it wins a battle
func Suicide(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Source == card {

			creature, err := card.Player.GetCard(event.Source.ID, match.BATTLEZONE)

			if err == nil {

				ctx.Match.Destroy(creature, event.Card)

			}

		}

	}

}
