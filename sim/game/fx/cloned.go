package fx

import (
	"duel-masters/game/match"
)

// ClonedCopiesInGraveyards counts the copies of a card sitting in either
// player's graveyard.
//
// It exists for the DM-12 "Cloned" cycle, which each choose one mandatory
// target and then "for each <this card> in each graveyard, you may choose
// another". The count is therefore the number of *extra* targets on offer, and
// the total selection is one more than this.
//
// Copies are matched by printing rather than by name, and a card never counts
// itself. That self-exclusion matters more than it looks: a spell is still in
// its owner's hand while it resolves, so it could only reach the graveyard it
// is counting after the fact, but Cloned Spike-Horn asks the same question
// about a creature that may well be lying in a graveyard itself.
func ClonedCopiesInGraveyards(card *match.Card, m *match.Match) int {
	copies := 0

	for _, player := range []*match.Player{m.Player1.Player, m.Player2.Player} {
		copies += len(FindFilter(player, match.GRAVEYARD, func(x *match.Card) bool {
			return x.ImageID == card.ImageID && x.ID != card.ID
		}))
	}

	return copies
}
