package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	cosmicDartsUID         = "0fe53ab5-bef9-4bbd-bb05-c94ccb9b1342"
	cosmicDartsSpellUID    = "b7f236fd-e7eb-41cc-912a-5239c134f265" // Energy Stream (draw 2, no further prompts)
	cosmicDartsCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	cosmicDartsSrc         = "cosmic_darts_test_setup"
)

// replaceShields empties a player's shield zone and refills it with fresh
// copies of uid, so a test controls exactly what the opponent has to choose
// from.
func replaceShields(t *testing.T, player *match.PlayerReference, uid string, count int) []*match.Card {
	t.Helper()

	shields, err := player.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)

	for _, shield := range shields {
		_, err := player.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, cosmicDartsSrc)
		require.NoError(t, err)
	}

	result := make([]*match.Card, 0, count)
	for range count {
		card, err := player.Player.SpawnCard(uid, match.SHIELDZONE)
		require.NoError(t, err)
		result = append(result, card)
	}

	return result
}

func TestCosmicDarts(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(cosmicDartsUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Cosmic Darts", spell.Name)
		assert.Equal(t, 1, spell.ManaCost)
		assert.Equal(t, []string{civ.Light}, spell.Civs)
		assert.Equal(t, []string{civ.Light}, spell.ManaRequirement)
	})

	t.Run("the caster may cast a spell shield the opponent reveals for free", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		shields := replaceShields(t, player, cosmicDartsSpellUID, 1)
		passTurnToSelf(t, scn, player, opponent)
		require.True(t, shields[0].HasCondition(cnd.Spell))

		castSpell(t, scn, player, cosmicDartsUID)

		// The opponent chooses which of the caster's shields gets revealed. The
		// revealed spell is then shown to the caster as a non-dismissible
		// preview before the "may cast it for free" prompt, both on the
		// caster's own connection, not the opponent's.
		playerMsgStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForMessage(player, playerMsgStart, "show_cards_non_dismissible"))
		playerMsgStart, err = scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player)) // dismiss preview
		require.NoError(t, scn.WaitForMessage(player, playerMsgStart, "action"))

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		countBefore := len(deckBefore)

		// The caster may cast the revealed spell immediately for no cost.
		answerInTurn(t, scn, player)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, shields[0].Zone, "Energy Stream resolves straight to the graveyard")

		deckAfter, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deckAfter, countBefore-2, "casting Energy Stream for free still draws 2 cards")
	})

	t.Run("the caster may decline to cast the revealed spell", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		shields := replaceShields(t, player, cosmicDartsSpellUID, 1)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, cosmicDartsUID)

		playerMsgStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForMessage(player, playerMsgStart, "show_cards_non_dismissible"))
		playerMsgStart, err = scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player)) // dismiss preview
		require.NoError(t, scn.WaitForMessage(player, playerMsgStart, "action"))

		// Declining leaves the shield exactly where it was.
		cancelInTurn(t, scn, player)
		settleTurn(t, scn)

		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)
	})

	t.Run("a creature shield cannot be cast and is simply put back", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		shields := replaceShields(t, player, cosmicDartsCreatureUID, 1)
		passTurnToSelf(t, scn, player, opponent)
		require.False(t, shields[0].HasCondition(cnd.Spell))

		castSpell(t, scn, player, cosmicDartsUID)

		messageCountBefore, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		settleTurn(t, scn)

		assert.Equal(t, match.SHIELDZONE, shields[0].Zone, "a creature shield is never cast")

		headers, err := scn.MessageHeaders(player, messageCountBefore)
		require.NoError(t, err)
		for _, header := range headers {
			assert.NotEqual(t, "action", header, "no cast-it-for-free prompt is offered for a non-spell shield")
		}
	})

	t.Run("with no shields the effect quietly does nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		replaceShields(t, player, cosmicDartsSpellUID, 0)

		castSpell(t, scn, player, cosmicDartsUID)
		settleTurn(t, scn)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, 0)
	})
}
