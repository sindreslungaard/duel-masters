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
	mysticMagicianUID         = "45b61aee-5393-4268-ad4b-a678c18b1584"
	mysticMagicianSilentUID   = "460fc2eb-c7cd-42d5-9bed-a98de4f59026" // Milporo (silent skill)
	mysticMagicianOrdinaryUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	mysticMagicianSetupSrc    = "mystic_magician_test_setup"
)

func TestMysticMagician(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Mystic Magician", magician.Name)
		assert.Equal(t, 3000, magician.Power)
		assert.Equal(t, 5, magician.ManaCost)
		assert.Equal(t, []string{civ.Water}, magician.Civs)
		assert.Equal(t, []string{civ.Water}, magician.ManaRequirement)
		assert.True(t, magician.HasFamily(family.Merfolk))
	})

	t.Run("its controller's silent skill creatures enter the battle zone tapped", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		silent, ordinary := mysticMagicianPreparedPair(t, player)
		passTurnToSelf(t, scn, player, opponent)

		arrivedSilent := mysticMagicianPutIntoPlay(t, scn, player, silent)
		arrivedOrdinary := mysticMagicianPutIntoPlay(t, scn, player, ordinary)

		assert.True(t, arrivedSilent.Tapped, "a silent skill creature arrives tapped, ready to use its ability")
		assert.False(t, arrivedOrdinary.Tapped, "everything else arrives as normal")
	})

	t.Run("it does not tap the opponent's silent skill creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		theirs, err := opponent.Player.SpawnCard(mysticMagicianSilentUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, theirs.HasCondition(cnd.SilentSkill))

		arrived := mysticMagicianPutIntoPlay(t, scn, opponent, theirs)

		assert.False(t, arrived.Tapped, "the printed wording is \"your creatures\"")
	})

	t.Run("a silent skill creature goes to its owner's hand instead of being destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		silent, ordinary := mysticMagicianPreparedPair(t, player)
		passTurnToSelf(t, scn, player, opponent)

		inPlaySilent := mysticMagicianPutIntoPlay(t, scn, player, silent)
		inPlayOrdinary := mysticMagicianPutIntoPlay(t, scn, player, ordinary)

		scn.Match.Destroy(inPlaySilent, magician, match.DestroyedByMiscAbility)
		scn.Match.Destroy(inPlayOrdinary, magician, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, inPlaySilent.Zone)
		assert.Equal(t, match.GRAVEYARD, inPlayOrdinary.Zone, "the replacement is only for silent skill creatures")
	})

	t.Run("it does not save the opponent's silent skill creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		theirs, err := opponent.Player.SpawnCard(mysticMagicianSilentUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		inPlay := mysticMagicianPutIntoPlay(t, scn, opponent, theirs)

		scn.Match.Destroy(inPlay, magician, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, inPlay.Zone)
	})

	t.Run("it does not save itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.False(t, magician.HasCondition(cnd.SilentSkill))

		scn.Match.Destroy(magician, magician, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, magician.Zone)
	})

	t.Run("neither half works once it has left the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		silent, err := player.Player.SpawnCard(mysticMagicianSilentUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		_, err = player.Player.MoveCard(magician.ID, match.BATTLEZONE, match.GRAVEYARD, mysticMagicianSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		inPlay := mysticMagicianPutIntoPlay(t, scn, player, silent)
		assert.False(t, inPlay.Tapped)

		scn.Match.Destroy(inPlay, magician, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, inPlay.Zone)
	})

	t.Run("two copies still return the creature exactly once", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		magician := putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)
		putCardInBattlezone(t, scn, player.Player, mysticMagicianUID, mysticMagicianSetupSrc)

		silent, err := player.Player.SpawnCard(mysticMagicianSilentUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		inPlay := mysticMagicianPutIntoPlay(t, scn, player, silent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		scn.Match.Destroy(inPlay, magician, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Equal(t, match.HAND, inPlay.Zone)
		assert.Len(t, hand, len(handBefore)+1, "the second copy must not move it again")
	})
}

// mysticMagicianPreparedPair puts one silent skill creature and one ordinary
// creature into the player's hand. They are spawned before the turn handover so
// an untap step gives them their conditions, the way a card that was in the
// deck all game would have them.
func mysticMagicianPreparedPair(t *testing.T, player *match.PlayerReference) (*match.Card, *match.Card) {
	t.Helper()

	silent, err := player.Player.SpawnCard(mysticMagicianSilentUID, match.HAND)
	require.NoError(t, err)

	ordinary, err := player.Player.SpawnCard(mysticMagicianOrdinaryUID, match.HAND)
	require.NoError(t, err)

	return silent, ordinary
}

func mysticMagicianPutIntoPlay(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, card *match.Card) *match.Card {
	t.Helper()

	moved, err := player.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, mysticMagicianSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	require.NoError(t, scn.WaitForEventLoop())

	return moved
}
