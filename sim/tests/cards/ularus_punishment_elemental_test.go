package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ularusUID      = "05c5496d-e5fa-4691-8542-2d6c6919f402"
	ularusOtherUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	ularusSrc      = "ularus_punishment_elemental_test_setup"
)

// faceUpShields counts how many of a player's shields are showing.
func faceUpShields(t *testing.T, player *match.Player) int {
	t.Helper()

	shields, err := player.Container(match.SHIELDZONE)
	require.NoError(t, err)

	showing := 0
	for _, shield := range shields {
		if shield.ShieldFaceUp {
			showing++
		}
	}

	return showing
}

func TestUlarusPunishmentElemental(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Kept in hand, because arriving opens a prompt.
		card := spawnForLater(t, player, ularusUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Ularus, Punishment Elemental", 4500, 5, []string{civ.Light})
		assert.True(t, card.HasFamily(family.AngelCommand))
	})

	t.Run("alone it turns up a single shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, ularusUID)

		action, err := scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)
		assert.Contains(t, action.Cards, "Your shields")
		assert.Contains(t, action.Cards, "Your opponent's shields")

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		answerInTurn(t, scn, player, shields[0].ID)
		settleTurn(t, scn)

		assert.Equal(t, 1, faceUpShields(t, opponent.Player))
		assert.Equal(t, 0, faceUpShields(t, player.Player))
	})

	t.Run("every creature it has buys another shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, ularusOtherUID, ularusSrc)
		putCardInBattlezone(t, scn, player.Player, ularusOtherUID, ularusSrc)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, ularusUID)

		action, err := scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 3, action.MaxSelections, "two creatures plus Ularus itself")

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		answerInTurn(t, scn, player, shields[0].ID, shields[1].ID, shields[2].ID)
		settleTurn(t, scn)

		assert.Equal(t, 3, faceUpShields(t, opponent.Player))
	})

	t.Run("either shield zone is fair game", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, ularusOtherUID, ularusSrc)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, ularusUID)

		_, err = scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)

		mine, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		theirs, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		answerInTurn(t, scn, player, mine[0].ID, theirs[0].ID)
		settleTurn(t, scn)

		assert.Equal(t, 1, faceUpShields(t, player.Player))
		assert.Equal(t, 1, faceUpShields(t, opponent.Player))
	})

	t.Run("it may turn up nothing at all", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, ularusUID)

		_, err = scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)

		cancelInTurn(t, scn, player)
		settleTurn(t, scn)

		assert.Equal(t, 0, faceUpShields(t, opponent.Player))
		assert.Equal(t, 0, faceUpShields(t, player.Player))
	})

	t.Run("a shield already showing is not offered twice", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shields[0].ShieldFaceUp = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, ularusUID)

		action, err := scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)

		offered := make([]string, 0)
		for _, group := range action.Cards {
			for _, card := range group {
				offered = append(offered, card.CardID)
			}
		}
		assert.NotContains(t, offered, shields[0].ID)

		cancelInTurn(t, scn, player)
		settleTurn(t, scn)
	})
}
