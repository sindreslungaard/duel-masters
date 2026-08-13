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
	hourglassMutantUID          = "06e93dfb-b8d7-4307-a444-a3eba52ad63c"
	hourglassMutantWaterUID     = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle (1000)
	hourglassMutantFireUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	hourglassMutantNatureUID    = "1d72eb3e-5185-449a-a16f-391bd2338343" // Burning Mane
	hourglassMutantLightUID     = "7b58e8c2-0b1e-4ef5-812f-e667c2092c73" // Reusol, the Oracle
	hourglassMutantDefenderUID  = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (5000)
	hourglassMutantTestSetupSrc = "hourglass_mutant_test_setup"
)

func TestHourglassMutant(t *testing.T) {
	t.Run("grants slayer only to its controller's water and fire creatures", func(t *testing.T) {
		scn, owner, opponent := setupHourglassMutantTest(t)
		hourglass := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		ownWater := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		ownFire := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantFireUID)
		ownNature := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantNatureUID)
		ownLight := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantLightUID)
		opponentWater := putHourglassMutantTestCardInBattlezone(t, scn, opponent.Player, hourglassMutantWaterUID)
		refreshHourglassMutantTestTurn(t, scn)

		assert.Equal(t, "Hourglass Mutant", hourglass.Name)
		assert.Equal(t, 2000, hourglass.Power)
		assert.Equal(t, 3, hourglass.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, hourglass.Civs)
		assert.True(t, hourglass.HasFamily(family.Hedrian))

		assert.True(t, ownWater.HasCondition(cnd.Slayer))
		assert.True(t, ownFire.HasCondition(cnd.Slayer))
		assert.False(t, ownNature.HasCondition(cnd.Slayer))
		assert.False(t, ownLight.HasCondition(cnd.Slayer))
		assert.False(t, opponentWater.HasCondition(cnd.Slayer), "only the controller's creatures are affected")
		assert.False(t, hourglass.HasCondition(cnd.Slayer), "Hourglass Mutant is a darkness creature")

		// The grant must survive the end-of-turn condition wipe.
		refreshHourglassMutantTestTurn(t, scn)
		assert.True(t, ownWater.HasCondition(cnd.Slayer))
		assert.True(t, ownFire.HasCondition(cnd.Slayer))
	})

	t.Run("grants slayer to creatures that arrive later", func(t *testing.T) {
		scn, owner, _ := setupHourglassMutantTest(t)
		putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		owner.Player.SpawnCard(hourglassMutantWaterUID, match.HAND)
		refreshHourglassMutantTestTurn(t, scn)

		late, err := scn.FindCard(owner.Player, match.HAND, hourglassMutantWaterUID)
		require.NoError(t, err)
		moved, err := owner.Player.MoveCard(late.ID, match.HAND, match.BATTLEZONE, hourglassMutantTestSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, moved.Zone)

		assert.True(t, late.HasCondition(cnd.Slayer))
	})

	t.Run("removes the grant when it leaves the battle zone", func(t *testing.T) {
		scn, owner, _ := setupHourglassMutantTest(t)
		hourglass := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		ownWater := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		refreshHourglassMutantTestTurn(t, scn)
		require.True(t, ownWater.HasCondition(cnd.Slayer))

		moved, err := owner.Player.MoveCard(hourglass.ID, match.BATTLEZONE, match.GRAVEYARD, hourglassMutantTestSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, ownWater.HasCondition(cnd.Slayer))
	})

	t.Run("removes the grant from a creature that already left the battle zone", func(t *testing.T) {
		scn, owner, _ := setupHourglassMutantTest(t)
		hourglass := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		ownWater := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		refreshHourglassMutantTestTurn(t, scn)
		require.True(t, ownWater.HasCondition(cnd.Slayer))

		bounced, err := owner.Player.MoveCard(ownWater.ID, match.BATTLEZONE, match.HAND, hourglassMutantTestSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.HAND, bounced.Zone)

		moved, err := owner.Player.MoveCard(hourglass.ID, match.BATTLEZONE, match.GRAVEYARD, hourglassMutantTestSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, ownWater.HasCondition(cnd.Slayer), "a stale grant must not follow the creature back into play")
	})

	t.Run("two copies grant independently and only one grant is removed", func(t *testing.T) {
		scn, owner, _ := setupHourglassMutantTest(t)
		firstHourglass := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		ownWater := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		refreshHourglassMutantTestTurn(t, scn)

		grants := 0
		for _, condition := range ownWater.Conditions() {
			if condition.ID == cnd.Slayer {
				grants++
			}
		}
		assert.Equal(t, 2, grants, "each copy contributes its own source")

		moved, err := owner.Player.MoveCard(firstHourglass.ID, match.BATTLEZONE, match.GRAVEYARD, hourglassMutantTestSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.True(t, ownWater.HasCondition(cnd.Slayer), "the remaining copy still grants slayer")
	})

	t.Run("a granted slayer destroys the winner of a battle", func(t *testing.T) {
		scn, owner, opponent := setupHourglassMutantTest(t)
		putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantUID)
		attacker := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		defender := putHourglassMutantTestCardInBattlezone(t, scn, opponent.Player, hourglassMutantDefenderUID)
		refreshHourglassMutantTestTurn(t, scn)
		defender.Tapped = true

		require.True(t, attacker.HasCondition(cnd.Slayer))
		require.NoError(t, scn.ActionAttackCreature(owner, attacker.ID, defender.ID))

		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "the 1000 power attacker loses the battle")
		assert.Equal(t, match.GRAVEYARD, defender.Zone, "slayer destroys the winner after the battle")
	})

	t.Run("does not grant slayer while it is outside the battle zone", func(t *testing.T) {
		scn, owner, _ := setupHourglassMutantTest(t)
		owner.Player.SpawnCard(hourglassMutantUID, match.HAND)
		ownWater := putHourglassMutantTestCardInBattlezone(t, scn, owner.Player, hourglassMutantWaterUID)
		refreshHourglassMutantTestTurn(t, scn)

		assert.False(t, ownWater.HasCondition(cnd.Slayer))
	})
}

func setupHourglassMutantTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}

// refreshHourglassMutantTestTurn cycles a full turn so the untap step rebuilds
// the intrinsic creature conditions the continuous effect filters on and hands
// the turn back to its original owner.
func refreshHourglassMutantTestTurn(t *testing.T, scn *scenario.TestScenario) {
	t.Helper()

	current := scn.Match.CurrentPlayer()
	other := scn.Match.PlayerRef(scn.Match.Opponent(current.Player))
	require.NoError(t, scn.ActionEndTurn(current))
	require.NoError(t, scn.ActionEndTurn(other))
}

func putHourglassMutantTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, hourglassMutantTestSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
