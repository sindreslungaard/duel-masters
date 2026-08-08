package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	coreCrashLizardUID              = "e46de1fa-e307-4119-ac71-69c95dc5e443"
	coreCrashLizardManaUID          = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	coreCrashLizardShieldTriggerUID = "a9469dd3-73fa-4a8c-a9dd-7b852487308a" // Aqua Jolter (shield trigger)
	coreCrashLizardSetupSrc         = "core_crash_lizard_test_setup"
)

func TestCoreCrashLizard(t *testing.T) {
	t.Run("puts a chosen shield into the opponent's graveyard face down", func(t *testing.T) {
		scn, player, opponent, lizard := setupCoreCrashLizardTest(t)

		assert.Equal(t, "Core-Crash Lizard", lizard.Name)
		assert.Equal(t, 6000, lizard.Power)
		assert.Equal(t, 7, lizard.ManaCost)
		assert.Equal(t, civ.Fire, lizard.Civ)
		assert.True(t, lizard.HasFamily(family.MeltWarrior))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))
		assert.Equal(t, match.BATTLEZONE, lizard.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)
		assert.False(t, action.Cancellable, "the effect is mandatory")
		for _, shield := range action.Cards {
			assert.Equal(t, "backside", shield.ImageID, "shields stay hidden while being chosen")
			assert.Empty(t, shield.Name)
		}

		chosen, err := opponent.Player.GetCard(action.Cards[0].CardID, match.SHIELDZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, chosen.Zone)
		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-1)
	})

	t.Run("does not break the shield, so its shield trigger is never offered", func(t *testing.T) {
		scn, player, opponent, lizard := setupCoreCrashLizardTest(t)
		triggers := replaceCoreCrashLizardTestShields(t, scn, opponent.Player)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, action.Cards[0].CardID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		destroyed := 0
		for _, trigger := range triggers {
			if trigger.Zone == match.GRAVEYARD {
				destroyed++
			}
			assert.NotEqual(t, match.HAND, trigger.Zone, "a shield put into the graveyard never reaches the hand")
			assert.NotEqual(t, match.BATTLEZONE, trigger.Zone)
		}
		assert.Equal(t, 1, destroyed)

		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "action", "no shield trigger prompt is opened")
	})

	t.Run("resolves without a prompt when the opponent has no shields", func(t *testing.T) {
		scn, player, opponent, lizard := setupCoreCrashLizardTest(t)
		clearCoreCrashLizardTestZone(t, opponent.Player, match.SHIELDZONE)

		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))

		assert.Equal(t, match.BATTLEZONE, lizard.Zone)
		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "action")
	})
}

func setupCoreCrashLizardTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(coreCrashLizardUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(coreCrashLizardManaUID, match.MANAZONE)
	}

	lizard, err := scn.FindCard(player.Player, match.HAND, coreCrashLizardUID)
	require.NoError(t, err)

	return scn, player, opponent, lizard
}

// replaceCoreCrashLizardTestShields swaps the opponent's shields for cards that
// do have a shield trigger, so a broken shield would be observable.
func replaceCoreCrashLizardTestShields(t *testing.T, scn *scenario.TestScenario, player *match.Player) []*match.Card {
	t.Helper()

	clearCoreCrashLizardTestZone(t, player, match.SHIELDZONE)

	shields := make([]*match.Card, 0, 3)
	for range 3 {
		player.SpawnCard(coreCrashLizardShieldTriggerUID, match.HAND)
		card, err := scn.FindCard(player, match.HAND, coreCrashLizardShieldTriggerUID)
		require.NoError(t, err)
		moved, err := player.MoveCard(card.ID, match.HAND, match.SHIELDZONE, coreCrashLizardSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.SHIELDZONE, moved.Zone)
		shields = append(shields, moved)
	}

	return shields
}

func clearCoreCrashLizardTestZone(t *testing.T, player *match.Player, zone string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, coreCrashLizardSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
