package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// OpponentChoosesOwnShieldToGraveyard makes the opponent pick one of their own
// shields, face down, and put it into their graveyard. The shield is not
// broken, so its shield trigger is never offered.
func OpponentChoosesOwnShieldToGraveyard(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)

	SelectBackside(
		opponent,
		ctx.Match,
		opponent,
		match.SHIELDZONE,
		fmt.Sprintf("%s: Choose one of your shields to put into your graveyard.", card.Name),
		1,
		1,
		false,
	).Map(func(shield *match.Card) {
		moved, err := opponent.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)

		if err != nil || moved.Zone != match.GRAVEYARD {
			return
		}

		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s put one of their shields into their graveyard because of %s", opponent.Username(), card.Name))
	})
}

// OpponentKeepsOneOfTheseAndLosesTheOther offers the opponent a pair of their
// own cards. The one they choose goes to their hand; whatever is left over is
// handed to the loser callback.
//
// A single card is simply kept: there is nothing to lose alongside it.
func OpponentKeepsOneOfTheseAndLosesTheOther(card *match.Card, ctx *match.Context, chosen []*match.Card, zone string, text string, lose func(*match.Card)) {
	if len(chosen) < 1 {
		return
	}

	opponent := ctx.Match.Opponent(card.Player)

	offered := make(map[string]bool, len(chosen))
	for _, c := range chosen {
		offered[c.ID] = true
	}

	kept := ""

	SelectFilter(
		opponent,
		ctx.Match,
		opponent,
		zone,
		text,
		1,
		1,
		false,
		func(x *match.Card) bool { return offered[x.ID] },
		false,
	).Map(func(keep *match.Card) {
		kept = keep.ID

		moved, err := opponent.MoveCard(keep.ID, zone, match.HAND, card.ID)

		if err != nil || moved.Zone != match.HAND {
			return
		}

		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s returned %s to their hand because of %s", opponent.Username(), keep.Name, card.Name))
	})

	for _, c := range chosen {
		if c.ID == kept || c.Zone != zone {
			continue
		}

		lose(c)
	}
}

// OpponentChoosesOwnCreatureToHand makes the opponent choose one of their own
// creatures and return it to their hand.
func OpponentChoosesOwnCreatureToHand(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)

	Select(
		opponent,
		ctx.Match,
		opponent,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Choose one of your creatures and return it to your hand.", card.Name),
		1,
		1,
		false,
	).Map(func(creature *match.Card) {
		moved, err := opponent.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)

		if err != nil || moved.Zone != match.HAND {
			return
		}

		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, opponent.Username(), card.Name))
	})
}
