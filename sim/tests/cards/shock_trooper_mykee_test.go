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
	shockTrooperMykeeUID      = "04dc2f45-6c44-486c-9bdf-9f99694792ff"
	shockTrooperMykeeManaUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	shockTrooperMykeeWeakUID  = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	shockTrooperMykeeToughUID = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (5000)
	shockTrooperMykeeSetupSrc = "shock_trooper_mykee_test_setup"
)

func TestShockTrooperMykee(t *testing.T) {
	t.Run("attacks the turn it is summoned and destroys a weak creature", func(t *testing.T) {
		scn, player, opponent := setupShockTrooperMykeeTest(t)
		weak := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeWeakUID)
		tough := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeToughUID)
		mykee := playShockTrooperMykee(t, scn, player)

		assert.Equal(t, "Shock Trooper Mykee", mykee.Name)
		assert.Equal(t, 1000, mykee.Power)
		assert.Equal(t, 6, mykee.ManaCost)
		assert.Equal(t, []string{civ.Fire}, mykee.Civs)
		assert.True(t, mykee.HasFamily(family.Human))
		assert.True(t, mykee.HasCondition(cnd.SpeedAttacker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		shieldAction, err := scn.ActionAttackPlayer(player, mykee.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldAction.Cards[0].CardID))

		destroyAction, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.True(t, destroyAction.Cancellable, "the effect is optional")
		offered := make([]string, 0, len(destroyAction.Cards))
		for _, card := range destroyAction.Cards {
			offered = append(offered, card.CardID)
		}
		assert.Contains(t, offered, weak.ID)
		assert.NotContains(t, offered, tough.ID, "5000 power is above the limit")

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, weak.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, weak.Zone)
		assert.Equal(t, match.BATTLEZONE, tough.Zone)
		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-1, "the shield is still broken")
	})

	t.Run("may decline to destroy a creature", func(t *testing.T) {
		scn, player, opponent := setupShockTrooperMykeeTest(t)
		weak := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeWeakUID)
		mykee := playShockTrooperMykee(t, scn, player)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		shieldAction, err := scn.ActionAttackPlayer(player, mykee.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldAction.Cards[0].CardID))

		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, weak.Zone)
	})

	t.Run("opens no prompt when no creature is weak enough", func(t *testing.T) {
		scn, player, opponent := setupShockTrooperMykeeTest(t)
		tough := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeToughUID)
		mykee := playShockTrooperMykee(t, scn, player)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		shieldAction, err := scn.ActionAttackPlayer(player, mykee.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldAction.Cards[0].CardID))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		actions := 0
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "only the shield selection is prompted")
		assert.Equal(t, match.BATTLEZONE, tough.Zone)
	})

	t.Run("does not trigger when the attack is blocked", func(t *testing.T) {
		scn, player, opponent := setupShockTrooperMykeeTest(t)
		weak := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeWeakUID)
		blocker := putShockTrooperMykeeTestCardInBattlezone(t, scn, opponent.Player, shockTrooperMykeeBlockerUID)
		mykee := playShockTrooperMykee(t, scn, player)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		shieldAction, err := scn.ActionAttackPlayer(player, mykee.ID)
		require.NoError(t, err)

		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldAction.Cards[0].CardID))

		blockAction, err := scn.WaitForAction(opponent, opponentStart)
		require.NoError(t, err)
		require.NotEmpty(t, blockAction.Cards)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, blocker.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, mykee.Zone, "the 1000 power attacker loses the battle")
		assert.Equal(t, match.BATTLEZONE, weak.Zone, "a blocked attack breaks no shields and destroys nothing")
		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount)
	})
}

const shockTrooperMykeeBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker, 2000)

func setupShockTrooperMykeeTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(shockTrooperMykeeUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(shockTrooperMykeeManaUID, match.MANAZONE)
	}

	return scn, player, opponent
}

// playShockTrooperMykee cycles a turn so the untap step grants Mykee its speed
// attacker condition, then summons it and returns it from the battle zone.
func playShockTrooperMykee(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	mykee, err := scn.FindCard(player.Player, match.HAND, shockTrooperMykeeUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, mykee.ID))
	require.Equal(t, match.BATTLEZONE, mykee.Zone)

	return mykee
}

func putShockTrooperMykeeTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, shockTrooperMykeeSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
