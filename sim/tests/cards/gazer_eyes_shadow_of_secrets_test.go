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
	gazerEyesUID      = "479cfe64-5fed-4afb-9d64-b8de738de8d2"
	gazerEyesSetupSrc = "gazer_eyes_shadow_of_secrets_test_setup"
)

func TestGazerEyesShadowOfSecrets(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gazer := putCardInBattlezone(t, scn, player.Player, gazerEyesUID, gazerEyesSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, gazer, "Gazer Eyes, Shadow of Secrets", 3000, 4, []string{civ.Darkness})
		assert.True(t, gazer.HasFamily(family.Ghost))
		assert.True(t, gazer.HasCondition(cnd.SilentSkill))
	})

	t.Run("its controller picks which card the opponent discards", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gazer := putCardInBattlezone(t, scn, player.Player, gazerEyesUID, gazerEyesSetupSrc)
		gazer.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(hand), 2)
		chosen := hand[1]

		useSilentSkill(t, scn, player)
		answerInTurn(t, scn, player, chosen.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, chosen.Zone)
		assert.Equal(t, match.HAND, hand[0].Zone, "only the chosen card is discarded")
	})

	t.Run("an empty opposing hand opens no prompt", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gazer := putCardInBattlezone(t, scn, player.Player, gazerEyesUID, gazerEyesSetupSrc)
		gazer.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		for _, card := range hand {
			_, err := opponent.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, gazerEyesSetupSrc)
			require.NoError(t, err)
		}

		useSilentSkill(t, scn, player)

		assert.True(t, gazer.Tapped)
	})
}
