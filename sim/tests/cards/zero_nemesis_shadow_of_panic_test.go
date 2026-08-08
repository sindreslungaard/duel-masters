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
	zeroNemesisShadowOfPanicUID = "ddf61f9a-752d-4072-a497-7ed8b14ca8ce"
	zeroNemesisAttackTestUID    = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
)

func TestZeroNemesisShadowOfPanic(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()

		owner.Player.SpawnCard(zeroNemesisShadowOfPanicUID, match.HAND)
		zeroNemesis, err := scn.FindCard(owner.Player, match.HAND, zeroNemesisShadowOfPanicUID)
		require.NoError(t, err)

		assert.Equal(t, "Zero Nemesis, Shadow of Panic", zeroNemesis.Name)
		assert.Equal(t, 6000, zeroNemesis.Power)
		assert.Equal(t, 6, zeroNemesis.ManaCost)
		assert.Equal(t, civ.Darkness, zeroNemesis.Civ)
		assert.Equal(t, []string{civ.Darkness}, zeroNemesis.ManaRequirement)
		assert.True(t, zeroNemesis.HasFamily(family.Ghost))

		scn.Match.HandleFx(match.NewContext(scn.Match, &match.UntapStep{}))
		assert.True(t, zeroNemesis.HasCondition(cnd.Creature))
		assert.True(t, zeroNemesis.HasCondition(cnd.Evolution))
		assert.True(t, zeroNemesis.HasCondition(cnd.DoubleBreaker))
	})

	for _, sourceZone := range []string{
		match.HAND,
		match.SHIELDZONE,
		match.HIDDENZONE,
		match.MANAZONE,
		match.GRAVEYARD,
		match.DECK,
	} {
		t.Run("discard is inactive in "+sourceZone, func(t *testing.T) {
			scn, owner, opponent, attacker, defender := setupZeroNemesisAttack(t, sourceZone, 1, 1)

			require.NoError(t, scn.ActionAttackCreature(owner, attacker.ID, defender.ID))

			hand, err := opponent.Player.Container(match.HAND)
			require.NoError(t, err)
			assert.Len(t, hand, 1)
		})
	}

	t.Run("discard is active in the battle zone", func(t *testing.T) {
		scn, owner, opponent, attacker, defender := setupZeroNemesisAttack(t, match.BATTLEZONE, 1, 1)

		require.NoError(t, scn.ActionAttackCreature(owner, attacker.ID, defender.ID))

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand)
	})

	t.Run("two copies in the battle zone trigger independently", func(t *testing.T) {
		scn, owner, opponent, attacker, defender := setupZeroNemesisAttack(t, match.BATTLEZONE, 2, 2)

		require.NoError(t, scn.ActionAttackCreature(owner, attacker.ID, defender.ID))

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand)
	})
}

func setupZeroNemesisAttack(t *testing.T, sourceZone string, sourceCount int, opponentHandCount int) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card, *match.Card) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	for range sourceCount {
		owner.Player.SpawnCard(zeroNemesisShadowOfPanicUID, match.HAND)
		zeroNemesis, err := scn.FindCard(owner.Player, match.HAND, zeroNemesisShadowOfPanicUID)
		require.NoError(t, err)

		if sourceZone != match.HAND {
			moved, err := owner.Player.MoveCard(zeroNemesis.ID, match.HAND, sourceZone, "zero_nemesis_test_setup")
			require.NoError(t, err)
			require.Equal(t, sourceZone, moved.Zone)
		}
	}

	attacker, err := scn.FindCard(owner.Player, match.HAND, zeroNemesisAttackTestUID)
	require.NoError(t, err)
	attacker, err = owner.Player.MoveCard(attacker.ID, match.HAND, match.BATTLEZONE, "zero_nemesis_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, attacker.Zone)

	defender, err := scn.FindCard(opponent.Player, match.HAND, zeroNemesisAttackTestUID)
	require.NoError(t, err)
	defender, err = opponent.Player.MoveCard(defender.ID, match.HAND, match.BATTLEZONE, "zero_nemesis_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, defender.Zone)
	defender.Tapped = true

	clearZeroNemesisTestHand(t, opponent.Player)
	for range opponentHandCount {
		opponent.Player.SpawnCard(zeroNemesisAttackTestUID, match.HAND)
	}

	return scn, owner, opponent, attacker, defender
}

func clearZeroNemesisTestHand(t *testing.T, player *match.Player) {
	t.Helper()

	hand, err := player.Container(match.HAND)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), hand...) {
		moved, err := player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, "zero_nemesis_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
