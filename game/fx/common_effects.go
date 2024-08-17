package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

func GiveTapAbilityToAllies(card *match.Card, ctx *match.Context, alliesFilter func(card *match.Card) bool, tapAbility func(card *match.Card, ctx *match.Context)) {
	// This is added for the case where the card is added to the field. There is another creature
	// that doesn't initially have a tap abbility but should receive one. The change doesn't propagate fast
	// enough to the FE and that creature doesn't get tap ability until another action takes places.
	// This is an ugly workaround.
	FindFilter(
		card.Player,
		match.BATTLEZONE,
		alliesFilter,
	).Map(func(x *match.Card) {
		x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
	})

	ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
		if card.Zone != match.BATTLEZONE {
			Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			exit()
			return
		}

		FindFilter(
			card.Player,
			match.BATTLEZONE,
			alliesFilter,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
		})
	})
}

func OpponentDiscardsRandomCard(card *match.Card, ctx *match.Context) {

	hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

	if err != nil || len(hand) < 1 {
		return
	}

	discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD, card.ID)
	if err == nil {
		ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
	}

}
