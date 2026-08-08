package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/server"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sirenConcertoUID = "49b6747c-6bdf-4cf0-9a5e-8978c9af15c1"

func TestSirenConcerto(t *testing.T) {
	t.Run("has its printed metadata and shield trigger", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		player.Player.SpawnCard(sirenConcertoUID, match.HAND)
		spell, err := scn.FindCard(player.Player, match.HAND, sirenConcertoUID)
		require.NoError(t, err)

		assert.Equal(t, "Siren Concerto", spell.Name)
		assert.Equal(t, 1, spell.ManaCost)
		assert.Equal(t, civ.Water, spell.Civ)

		require.NoError(t, scn.ActionEndTurn(player))
		assert.True(t, spell.HasCondition(cnd.Spell))
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("excludes only the resolving copy from the hand selection", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		clearSirenConcertoTestHand(t, player.Player)

		player.Player.SpawnCard(sirenConcertoUID, match.HAND)
		player.Player.SpawnCard(sirenConcertoUID, match.HAND)
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		require.Len(t, hand, 2)
		castingSpell := hand[0]
		otherCopy := hand[1]

		player.Player.SpawnCard(sirenConcertoUID, match.MANAZONE)
		retrievedCard, err := scn.FindCard(player.Player, match.MANAZONE, sirenConcertoUID)
		require.NoError(t, err)

		actionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, castingSpell.ID))
		action, err := scn.LatestAction(player, actionStart)
		require.NoError(t, err)
		offered := sirenConcertoActionCardIDs(action)
		assert.NotContains(t, offered, castingSpell.ID)
		assert.Contains(t, offered, otherCopy.ID)
		assert.Contains(t, offered, retrievedCard.ID)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, castingSpell.ID))
		require.NoError(t, scn.WaitForMessage(player, warningStart, "action_error"))

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, retrievedCard.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, castingSpell.Zone)
		assert.Equal(t, match.HAND, otherCopy.Zone)
		assert.Equal(t, match.MANAZONE, retrievedCard.Zone)
	})

	t.Run("automatically returns the only eligible card instead of selecting itself", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		clearSirenConcertoTestHand(t, player.Player)

		player.Player.SpawnCard(sirenConcertoUID, match.HAND)
		castingSpell, err := scn.FindCard(player.Player, match.HAND, sirenConcertoUID)
		require.NoError(t, err)
		player.Player.SpawnCard(sirenConcertoUID, match.MANAZONE)
		retrievedCard, err := scn.FindCard(player.Player, match.MANAZONE, sirenConcertoUID)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, castingSpell.ID))

		assert.Equal(t, match.GRAVEYARD, castingSpell.Zone)
		assert.Equal(t, match.MANAZONE, retrievedCard.Zone)
	})
}

func clearSirenConcertoTestHand(t *testing.T, player *match.Player) {
	t.Helper()

	hand, err := player.Container(match.HAND)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), hand...) {
		moved, err := player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, "siren_concerto_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}

func sirenConcertoActionCardIDs(action *server.ActionMessage) []string {
	ids := make([]string, len(action.Cards))
	for i, card := range action.Cards {
		ids[i] = card.CardID
	}
	return ids
}
