package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupDuel starts a match and returns the player whose turn it is plus their
// opponent.
func setupDuel(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

// assertPrinted checks the numbers and civilizations on the card face. Keyword
// conditions are only rebuilt during an untap step, so callers that assert on
// those have to pass a turn first.
func assertPrinted(t *testing.T, card *match.Card, name string, power int, cost int, civs []string) {
	t.Helper()

	assert.Equal(t, name, card.Name)
	assert.Equal(t, power, card.Power)
	assert.Equal(t, cost, card.ManaCost)
	assert.Equal(t, civs, card.Civs)
	assert.Equal(t, civs, card.ManaRequirement)
}

// answerDrawUpTo answers the prompts fx.DrawUpto opens. It asks once per card
// rather than once in total, and each accepted draw opens a preview of the card
// that has to be acknowledged before the effect continues.
func answerDrawUpTo(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, max int, take bool) {
	t.Helper()

	if !take {
		// Declining the first question ends the whole draw.
		cancelInTurn(t, scn, player)
		return
	}

	for range max {
		// Yes to this card, then close the preview it opens. The preview is
		// dismissed with a cancel, which is what its Close button sends.
		answerInTurn(t, scn, player)
		cancelInTurn(t, scn, player)
	}
}

// castSpell plays a spell from hand, paying for it with copies of itself so the
// civilization always matches.
func castSpell(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	spell, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	card, err := scn.FindCard(player.Player, match.HAND, uid)
	require.NoError(t, err)

	for range card.ManaCost {
		_, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, spell.ID))

	return spell
}
