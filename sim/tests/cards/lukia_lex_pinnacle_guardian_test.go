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
	lukiaLexPinnacleGuardianUID      = "4f8f46ae-d907-4e27-ba47-2273b5d9abbb"
	lukiaLexPinnacleGuardianSetupSrc = "lukia_lex_pinnacle_guardian_test_setup"
)

func TestLukiaLexPinnacleGuardian(t *testing.T) {
	t.Run("printed characteristics and power attacker", func(t *testing.T) {
		scn, player, opponent := setupLukiaLexTest(t)
		lukia := putCardInBattlezone(t, scn, player.Player, lukiaLexPinnacleGuardianUID, lukiaLexPinnacleGuardianSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Lukia Lex, Pinnacle Guardian", lukia.Name)
		assert.Equal(t, 2500, lukia.Power)
		assert.Equal(t, 3, lukia.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, lukia.Civs)
		assert.True(t, lukia.HasFamily(family.Guardian))
		assert.True(t, lukia.HasCondition(cnd.PowerAttacker))

		assert.Equal(t, 2500, scn.Match.GetPower(lukia, false))
		assert.Equal(t, 5500, scn.Match.GetPower(lukia, true), "power attacker +3000 while attacking")
	})

	t.Run("may untap itself at the end of its controller's turn", func(t *testing.T) {
		scn, player, _ := setupLukiaLexTest(t)
		lukia := putCardInBattlezone(t, scn, player.Player, lukiaLexPinnacleGuardianUID, lukiaLexPinnacleGuardianSetupSrc)
		lukia.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		action, err := scn.WaitForAction(player, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, action.Text)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, lukia.Tapped, "it untaps during its controller's end step")
		assert.False(t, scn.Match.IsPlayerTurn(player.Player))
	})

	t.Run("may decline to untap", func(t *testing.T) {
		scn, player, _ := setupLukiaLexTest(t)
		lukia := putCardInBattlezone(t, scn, player.Player, lukiaLexPinnacleGuardianUID, lukiaLexPinnacleGuardianSetupSrc)
		lukia.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, lukia.Tapped, "declining leaves it tapped")
	})

	t.Run("opens no prompt while it is untapped", func(t *testing.T) {
		scn, player, _ := setupLukiaLexTest(t)
		lukia := putCardInBattlezone(t, scn, player.Player, lukiaLexPinnacleGuardianUID, lukiaLexPinnacleGuardianSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.False(t, lukia.Tapped)
		assert.False(t, scn.Match.IsPlayerTurn(player.Player))
	})

	t.Run("does not untap on the opponent's turn", func(t *testing.T) {
		scn, player, opponent := setupLukiaLexTest(t)
		lukia := putCardInBattlezone(t, scn, player.Player, lukiaLexPinnacleGuardianUID, lukiaLexPinnacleGuardianSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		lukia.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(opponent))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "the ability only triggers at the end of its controller's turn")
		assert.False(t, lukia.Tapped, "its controller's untap step untapped it instead")
	})
}

func setupLukiaLexTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}
