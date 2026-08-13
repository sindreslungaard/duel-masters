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
	cometEyeUID         = "6da6ae55-baf5-45a6-b5a2-6b36acd4afa4"
	cometEyeKinUID      = "2e10b4fb-3f85-4144-8762-51c04fe609d5" // Scowling Tomato (Wild Veggies, 2000)
	cometEyeStrangerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	cometEyeSetupSrc    = "comet_eye_the_spectral_spud_test_setup"
)

func TestCometEyeTheSpectralSpud(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Comet Eye, The Spectral Spud", 5500, 4, []string{civ.Light, civ.Nature})
		assert.True(t, card.HasFamily(family.WildVeggies))
		assert.True(t, card.HasFamily(family.RainbowPhantom))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(cometEyeUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, cometEyeSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("its own other kin get +2000", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		cometEye := putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, cometEyeKinUID, cometEyeSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, cometEyeStrangerUID, cometEyeSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, cometEyeKinUID, cometEyeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		assert.Equal(t, stranger.Power, scn.Match.GetPower(stranger, false), "only its own races")
		assert.Equal(t, theirKin.Power, scn.Match.GetPower(theirKin, false), "only its controller's creatures")
		assert.Equal(t, cometEye.Power, scn.Match.GetPower(cometEye, false), "\"other\" spares itself")
	})

	t.Run("the bonus leaves with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		cometEye := putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, cometEyeKinUID, cometEyeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))

		_, err := player.Player.MoveCard(cometEye.ID, match.BATTLEZONE, match.GRAVEYARD, cometEyeSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power, scn.Match.GetPower(kin, false))
	})

	t.Run("it may untap any number of its kin at the end of its turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		cometEye := putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, cometEyeKinUID, cometEyeSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, cometEyeStrangerUID, cometEyeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		kin.Tapped = true
		cometEye.Tapped = true
		stranger.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		action, err := scn.WaitForAction(player, promptStart)
		require.NoError(t, err, "the untap should have been offered")

		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		assert.Contains(t, offered, kin.ID)
		assert.Contains(t, offered, cometEye.ID, "the clause covers its own races, itself included")
		assert.NotContains(t, offered, stranger.ID)

		require.NoError(t, scn.SubmitAction(player, kin.ID, cometEye.ID))
		settleTurn(t, scn)

		assert.False(t, kin.Tapped)
		assert.False(t, cometEye.Tapped)
		assert.True(t, stranger.Tapped, "a stranger stays down")
	})

	t.Run("the untap may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, cometEyeKinUID, cometEyeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		kin.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		assert.True(t, kin.Tapped, "declined, and the opponent's turn does not untap it")
	})

	t.Run("nothing tapped asks nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, cometEyeUID, cometEyeSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, cometEyeStrangerUID, cometEyeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		stranger.Tapped = true

		// A tapped stranger is not a tapped kin, so the turn ends without a
		// prompt and the event loop settles on its own.
		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assert.True(t, scn.Match.IsPlayerTurn(opponent.Player))
	})
}
