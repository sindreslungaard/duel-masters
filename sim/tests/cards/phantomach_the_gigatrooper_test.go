package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	phantomachUID         = "6a5fc9ec-17c2-4e3f-92ce-5f67344895a0"
	phantomachKinUID      = "5d73062e-acff-47e6-b49a-c0bb1a1762b5" // Gigagiele (Chimera, 3000)
	phantomachStrangerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	phantomachSetupSrc    = "phantomach_the_gigatrooper_test_setup"
)

func TestPhantomachTheGigatrooper(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Phantomach, the Gigatrooper", 6000, 5, []string{civ.Darkness, civ.Fire})
		assert.True(t, card.HasFamily(family.Chimera))
		assert.True(t, card.HasFamily(family.Armorloid))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(phantomachUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, phantomachSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("its own other kin get +2000", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phantomach := putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, phantomachKinUID, phantomachSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, phantomachStrangerUID, phantomachSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, phantomachKinUID, phantomachSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		assert.Equal(t, stranger.Power, scn.Match.GetPower(stranger, false), "only its own races")
		assert.Equal(t, theirKin.Power, scn.Match.GetPower(theirKin, false), "only its controller's creatures")
		assert.Equal(t, phantomach.Power, scn.Match.GetPower(phantomach, false), "\"other\" spares itself")
	})

	t.Run("its kin get double breaker, itself included", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phantomach := putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, phantomachKinUID, phantomachSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, phantomachStrangerUID, phantomachSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, phantomachKinUID, phantomachSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.True(t, kin.HasCondition(cnd.DoubleBreaker))
		assert.True(t, phantomach.HasCondition(cnd.DoubleBreaker), "the clause says each, not other")
		assert.False(t, stranger.HasCondition(cnd.DoubleBreaker))
		assert.False(t, theirKin.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("a granted double breaker really takes two shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, phantomachKinUID, phantomachSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)
		require.GreaterOrEqual(t, shieldCount, 2)

		action, err := scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections, "a granted double breaker still asks for two")

		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))
		settleTurn(t, scn)

		shieldsAfter, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		assert.Len(t, shieldsAfter, shieldCount-2)
	})

	t.Run("both grants leave with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phantomach := putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, phantomachKinUID, phantomachSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		require.True(t, kin.HasCondition(cnd.DoubleBreaker))

		_, err := player.Player.MoveCard(phantomach.ID, match.BATTLEZONE, match.GRAVEYARD, phantomachSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power, scn.Match.GetPower(kin, false))
		assert.False(t, kin.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("two copies grant one bonus each and only one is lost", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		putCardInBattlezone(t, scn, player.Player, phantomachUID, phantomachSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, phantomachKinUID, phantomachSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+4000, scn.Match.GetPower(kin, false))

		_, err := player.Player.MoveCard(first.ID, match.BATTLEZONE, match.GRAVEYARD, phantomachSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false), "only one contribution was withdrawn")
		assert.True(t, kin.HasCondition(cnd.DoubleBreaker), "the survivor still grants it")
	})
}
