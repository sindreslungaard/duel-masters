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
	necrodragonBryzenagaUID        = "8f0df729-fa1a-439e-a4a1-961650bcd192"
	necrodragonBryzenagaManaUID    = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul
	necrodragonBryzenagaPlainUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	necrodragonBryzenagaTriggerUID = "a9469dd3-73fa-4a8c-a9dd-7b852487308a" // Aqua Jolter (shield trigger)
	necrodragonBryzenagaSetupSrc   = "necrodragon_bryzenaga_test_setup"
)

func TestNecrodragonBryzenaga(t *testing.T) {
	t.Run("puts every shield into its controller's hand", func(t *testing.T) {
		scn, player, opponent := setupNecrodragonBryzenagaTest(t)
		shields := setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaPlainUID, 3)
		opponentShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		opponentShieldCount := len(opponentShields)

		bryzenaga := playNecrodragonBryzenaga(t, scn, player)

		assert.Equal(t, "Necrodragon Bryzenaga", bryzenaga.Name)
		assert.Equal(t, 9000, bryzenaga.Power)
		assert.Equal(t, 6, bryzenaga.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, bryzenaga.Civs)
		assert.True(t, bryzenaga.HasFamily(family.ZombieDragon))

		for _, shield := range shields {
			assert.Equal(t, match.HAND, shield.Zone)
		}
		remaining, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Empty(t, remaining)

		remainingOpponent, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingOpponent, opponentShieldCount, "only its controller's shields move")
	})

	t.Run("is a double breaker after the untap step", func(t *testing.T) {
		scn, player, opponent := setupNecrodragonBryzenagaTest(t)
		setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaPlainUID, 1)
		bryzenaga := playNecrodragonBryzenaga(t, scn, player)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, bryzenaga.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("offers the shield triggers of the shields it moves", func(t *testing.T) {
		scn, player, _ := setupNecrodragonBryzenagaTest(t)
		triggers := setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaTriggerUID, 2)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		bryzenaga := playNecrodragonBryzenagaExpectingPrompt(t, scn, player)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.True(t, action.Cancellable, "using a shield trigger is optional")
		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		for _, trigger := range triggers {
			assert.Contains(t, offered, trigger.ID)
		}

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, triggers[0].ID))

		// One trigger is still in hand, so the prompt reopens for it.
		secondAction, err := scn.WaitForAction(player, completionStart)
		require.NoError(t, err)
		require.NotEmpty(t, secondAction.Cards)

		declineStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, declineStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, triggers[0].Zone, "the played trigger was summoned for free")
		assert.Equal(t, match.HAND, triggers[1].Zone, "the declined trigger stays in hand")
		assert.Equal(t, match.BATTLEZONE, bryzenaga.Zone)
	})

	t.Run("does not offer a prompt for shields without a shield trigger", func(t *testing.T) {
		scn, player, _ := setupNecrodragonBryzenagaTest(t)
		setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaPlainUID, 3)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		playNecrodragonBryzenaga(t, scn, player)

		actions := 0
		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "only the mana payment is prompted")
	})

	t.Run("resolves with no shields at all", func(t *testing.T) {
		scn, player, _ := setupNecrodragonBryzenagaTest(t)
		setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaPlainUID, 0)

		bryzenaga := playNecrodragonBryzenaga(t, scn, player)
		assert.Equal(t, match.BATTLEZONE, bryzenaga.Zone)
	})

	t.Run("does not count as breaking shields for turbo rush", func(t *testing.T) {
		scn, player, _ := setupNecrodragonBryzenagaTest(t)
		setNecrodragonBryzenagaTestShields(t, scn, player, necrodragonBryzenagaPlainUID, 3)

		turboRush := putNecrodragonBryzenagaTestCardInBattlezone(t, scn, player.Player, necrodragonBryzenagaTurboRushUID)
		require.Equal(t, 2000, scn.Match.GetPower(turboRush, true))

		playNecrodragonBryzenaga(t, scn, player)

		assert.Equal(t, 2000, scn.Match.GetPower(turboRush, true), "moving your own shields is not a shield break")
	})
}

// Missile Soldier Ultimo gets +3000 power while attacking if one of your other
// creatures broke a shield this turn.
const necrodragonBryzenagaTurboRushUID = "cfbaa44c-b949-45fc-8fe6-99bbab5dab93"

func setupNecrodragonBryzenagaTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(necrodragonBryzenagaUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(necrodragonBryzenagaManaUID, match.MANAZONE)
	}

	return scn, player, opponent
}

// playNecrodragonBryzenaga cycles a turn first so the untap step has granted the
// shield zone cards their shield trigger conditions, then summons Bryzenaga.
func playNecrodragonBryzenaga(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	bryzenaga := prepareNecrodragonBryzenaga(t, scn, player)
	require.NoError(t, scn.ActionPlayCard(player, bryzenaga.ID))
	require.Equal(t, match.BATTLEZONE, bryzenaga.Zone)

	return bryzenaga
}

// playNecrodragonBryzenagaExpectingPrompt is the same but leaves the shield
// trigger prompt outstanding for the caller to answer.
func playNecrodragonBryzenagaExpectingPrompt(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	bryzenaga := prepareNecrodragonBryzenaga(t, scn, player)
	require.NoError(t, scn.ActionPlayCard(player, bryzenaga.ID))

	return bryzenaga
}

func prepareNecrodragonBryzenaga(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	bryzenaga, err := scn.FindCard(player.Player, match.HAND, necrodragonBryzenagaUID)
	require.NoError(t, err)

	return bryzenaga
}

// setNecrodragonBryzenagaTestShields replaces a player's shields with exactly
// count copies of uid so the shield trigger outcome is known.
func setNecrodragonBryzenagaTestShields(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string, count int) []*match.Card {
	t.Helper()

	existing, err := player.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)
	for _, shield := range append([]*match.Card(nil), existing...) {
		moved, err := player.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, necrodragonBryzenagaSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}

	shields := make([]*match.Card, 0, count)
	for range count {
		player.Player.SpawnCard(uid, match.HAND)
		card, err := scn.FindCard(player.Player, match.HAND, uid)
		require.NoError(t, err)
		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.SHIELDZONE, necrodragonBryzenagaSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.SHIELDZONE, moved.Zone)
		shields = append(shields, moved)
	}

	return shields
}

func putNecrodragonBryzenagaTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, necrodragonBryzenagaSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
