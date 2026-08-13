package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sanfistTheSavageVizierUID       = "ad428b55-4b0b-4a63-b48b-dfc1541b8b81"
	sanfistTheSavageVizierMillerUID = "ba955ab0-5bb3-4aaf-82f3-293522e65a9c" // Locomotiver, discards on summon
	sanfistTheSavageVizierManaUID   = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul (darkness mana)
	sanfistTheSavageVizierSetupSrc  = "sanfist_the_savage_vizier_test_setup"
)

func TestSanfistTheSavageVizier(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		sanfist := putCardInBattlezone(t, scn, player.Player, sanfistTheSavageVizierUID, sanfistTheSavageVizierSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Sanfist, the Savage Vizier", sanfist.Name)
		assert.Equal(t, 3000, sanfist.Power)
		assert.Equal(t, 3, sanfist.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, sanfist.Civs)
		assert.Equal(t, []string{civ.Light, civ.Nature}, sanfist.ManaRequirement)
		assert.True(t, sanfist.IsMulticolored())
		assert.True(t, sanfist.HasFamily(family.BeastFolk))
		assert.True(t, sanfist.HasFamily(family.Initiate))
		assert.True(t, sanfist.HasCondition(cnd.Blocker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("a discard on the opponent's turn may put it into the battle zone instead", func(t *testing.T) {
		scn, player, opponent, sanfist := setupSanfistDiscardTest(t)

		require.NoError(t, scn.ActionPlayCard(opponent, sanfistDiscarder(t, scn, opponent)))

		// The prompt belongs to Sanfist's controller, who is not the player
		// taking the turn.
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, sanfist.Zone, "it is put into play instead of discarded")
	})

	t.Run("declining discards it as normal", func(t *testing.T) {
		scn, player, opponent, sanfist := setupSanfistDiscardTest(t)

		require.NoError(t, scn.ActionPlayCard(opponent, sanfistDiscarder(t, scn, opponent)))
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, sanfist.Zone)
	})

	t.Run("it arrives with summoning sickness", func(t *testing.T) {
		scn, player, opponent, sanfist := setupSanfistDiscardTest(t)

		require.NoError(t, scn.ActionPlayCard(opponent, sanfistDiscarder(t, scn, opponent)))
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		require.Equal(t, match.BATTLEZONE, sanfist.Zone)
		assert.True(t, sanfist.HasCondition(cnd.SummoningSickness), "it was put into play, not summoned early")
	})

	t.Run("a discard on its controller's own turn is not replaced", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		sanfist, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Discarding during its controller's own turn: the printed trigger
		// only covers the opponent's turn, so nothing should be offered.
		_, err = player.Player.MoveCard(sanfist.ID, match.HAND, match.GRAVEYARD, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.GRAVEYARD, sanfist.Zone)
	})

	t.Run("moving it out of hand for another reason is not a discard", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		sanfist, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Charging it as mana on the opponent's turn is still a move out of
		// hand, but it is not a discard.
		_, err = player.Player.MoveCard(sanfist.ID, match.HAND, match.MANAZONE, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.MANAZONE, sanfist.Zone)
	})
}

// setupSanfistDiscardTest leaves the opponent on turn with a creature ready to
// summon that discards a random card, and Sanfist as the only card in the other
// player's hand so the discard can only hit Sanfist.
func setupSanfistDiscardTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	for range 4 {
		_, err := opponent.Player.SpawnCard(sanfistTheSavageVizierManaUID, match.MANAZONE)
		require.NoError(t, err)
	}

	_, err := opponent.Player.SpawnCard(sanfistTheSavageVizierMillerUID, match.HAND)
	require.NoError(t, err)

	require.NoError(t, scn.ActionEndTurn(player))
	require.True(t, scn.Match.IsPlayerTurn(opponent.Player))

	// Emptied last so the drawn opening hand cannot be discarded instead.
	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)
	for _, card := range hand {
		_, err := player.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)
	}

	sanfist, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
	require.NoError(t, err)

	return scn, player, opponent, sanfist
}

func sanfistDiscarder(t *testing.T, scn *scenario.TestScenario, opponent *match.PlayerReference) string {
	t.Helper()

	card, err := scn.FindCard(opponent.Player, match.HAND, sanfistTheSavageVizierMillerUID)
	require.NoError(t, err)

	return card.ID
}
