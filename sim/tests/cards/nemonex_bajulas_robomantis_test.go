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
	nemonexUID         = "8b9ed1f7-70d7-4932-9381-25040dc69d6e"
	nemonexKinUID      = "27c58f10-82ef-47ca-8a69-1cfa2057743d" // Picora's Wrench (Xenoparts, 2000)
	nemonexStrangerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	nemonexWallUID     = "85852774-dd96-4395-8980-eb5b85bf5bfc" // Ferrosaturn, Spectral Knight (blocker, 2000)
	nemonexSetupSrc    = "nemonex_bajulas_robomantis_test_setup"
)

func TestNemonexBajulasRobomantis(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Nemonex, Bajula's Robomantis", 5000, 6, []string{civ.Fire, civ.Nature})
		assert.True(t, card.HasFamily(family.Xenoparts))
		assert.True(t, card.HasFamily(family.GiantInsect))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(nemonexUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, nemonexSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("its own other kin get +2000", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		nemonex := putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, nemonexKinUID, nemonexSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, nemonexStrangerUID, nemonexSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, nemonexKinUID, nemonexSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		assert.Equal(t, stranger.Power, scn.Match.GetPower(stranger, false), "only its own races")
		assert.Equal(t, theirKin.Power, scn.Match.GetPower(theirKin, false), "only its controller's creatures")
		assert.Equal(t, nemonex.Power, scn.Match.GetPower(nemonex, false), "\"other\" spares itself")
	})

	t.Run("the bonus leaves with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		nemonex := putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, nemonexKinUID, nemonexSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))

		_, err := player.Player.MoveCard(nemonex.ID, match.BATTLEZONE, match.GRAVEYARD, nemonexSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power, scn.Match.GetPower(kin, false))
	})

	t.Run("an unblocked kin burns a mana the opponent picks", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, nemonexKinUID, nemonexSetupSrc)

		keep, err := opponent.Player.SpawnCard(nemonexStrangerUID, match.MANAZONE)
		require.NoError(t, err)
		burn, err := opponent.Player.SpawnCard(nemonexStrangerUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		choiceStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		_, err = scn.WaitForAction(opponent, choiceStart)
		require.NoError(t, err, "the opponent chooses which mana is lost")
		require.NoError(t, scn.SubmitAction(opponent, burn.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, burn.Zone)
		assert.Equal(t, match.MANAZONE, keep.Zone)
	})

	t.Run("an attack by a stranger burns nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, nemonexStrangerUID, nemonexSetupSrc)

		mana, err := opponent.Player.SpawnCard(nemonexStrangerUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, stranger.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		settleTurn(t, scn)

		assert.Equal(t, match.MANAZONE, mana.Zone, "only its own races trigger it")
	})

	t.Run("a blocked kin burns nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, nemonexUID, nemonexSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, nemonexKinUID, nemonexSetupSrc)
		wall := putCardInBattlezone(t, scn, opponent.Player, nemonexWallUID, nemonexSetupSrc)

		mana, err := opponent.Player.SpawnCard(nemonexStrangerUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, wall.HasCondition(cnd.Blocker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		blockStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		_, err = scn.WaitForAction(opponent, blockStart)
		require.NoError(t, err, "the blocker should have been offered")
		require.NoError(t, scn.SubmitAction(opponent, wall.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.MANAZONE, mana.Zone, "a blocked attack is not an unblocked one")
		assert.Equal(t, match.GRAVEYARD, wall.Zone, "4000 beats 2000")
	})
}
