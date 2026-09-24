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
	soulPhoenixUID            = "cfa7c730-2818-4f78-aa4b-1395060cd687"
	soulPhoenixFireBirdUID    = "479b477f-f535-4564-ac0c-5f0aeaafe914" // Baby Zoppe (Fire Bird)
	soulPhoenixEarthDragonUID = "ea9dbf9c-d049-4213-a658-e47ab25867e6" // Terradragon Regarion (Earth Dragon)
	soulPhoenixSetupSrc       = "soul_phoenix_avatar_of_unity_test_setup"
)

// summonSoulPhoenix pays for Soul Phoenix with copies of itself and plays it,
// then answers the stacking-order prompt that appears once its controller has
// exactly one Fire Bird and one Earth Dragon: with no real choice for either
// requirement, both selections resolve on their own.
func summonSoulPhoenix(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, topID string, bottomID string) *match.Card {
	t.Helper()

	phoenix, err := player.Player.SpawnCard(soulPhoenixUID, match.HAND)
	require.NoError(t, err)

	for range phoenix.ManaCost {
		_, err := player.Player.SpawnCard(soulPhoenixUID, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, phoenix.ID))

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected the stacking-order prompt")
	require.Len(t, action.Cards, 2)

	require.NoError(t, scn.SubmitAction(player, topID, bottomID))
	require.NoError(t, scn.WaitForEventLoop())

	return phoenix
}

func TestSoulPhoenixAvatarOfUnity(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phoenix := putCardInBattlezone(t, scn, player.Player, soulPhoenixUID, soulPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, phoenix, "Soul Phoenix, Avatar of Unity", 13000, 4, []string{civ.Fire, civ.Nature})
		assert.True(t, phoenix.HasFamily(family.Phoenix))
		assert.True(t, phoenix.IsMulticolored())
		assert.True(t, phoenix.HasCondition(cnd.TripleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(soulPhoenixUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, soulPhoenixSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("cannot be summoned without a Fire Bird and an Earth Dragon in the battle zone", func(t *testing.T) {
		cases := []struct {
			name        string
			fireBird    bool
			earthDragon bool
		}{
			{"neither", false, false},
			{"only a Fire Bird", true, false},
			{"only an Earth Dragon", false, true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				scn, player, _ := setupDuel(t)

				if c.fireBird {
					putCardInBattlezone(t, scn, player.Player, soulPhoenixFireBirdUID, soulPhoenixSetupSrc)
				}
				if c.earthDragon {
					putCardInBattlezone(t, scn, player.Player, soulPhoenixEarthDragonUID, soulPhoenixSetupSrc)
				}

				phoenix, err := player.Player.SpawnCard(soulPhoenixUID, match.HAND)
				require.NoError(t, err)
				for range phoenix.ManaCost {
					_, err := player.Player.SpawnCard(soulPhoenixUID, match.MANAZONE)
					require.NoError(t, err)
				}

				// Rejected before mana is even asked for, exactly like a
				// regular evolution with no legal base: no prompt is ever
				// opened for the caller to answer.
				require.Error(t, scn.ActionPlayCard(player, phoenix.ID), "there is nothing for it to evolve from")
				assert.Equal(t, match.HAND, phoenix.Zone, "it stays in hand")
			})
		}
	})

	t.Run("vortex evolution consumes one of each and stacks on top of the chosen base", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		fireBird := putCardInBattlezone(t, scn, player.Player, soulPhoenixFireBirdUID, soulPhoenixSetupSrc)
		earthDragon := putCardInBattlezone(t, scn, player.Player, soulPhoenixEarthDragonUID, soulPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Tapped, and chosen as the top of the stack: Soul Phoenix should
		// take this tap state, matching a regular evolution landing on a
		// tapped creature.
		earthDragon.Tapped = true

		phoenix := summonSoulPhoenix(t, scn, player, earthDragon.ID, fireBird.ID)

		require.Equal(t, match.BATTLEZONE, phoenix.Zone)
		assert.Equal(t, 13000, phoenix.Power)
		assert.True(t, phoenix.Tapped, "it took the tap state of the base it was put on top of")

		assert.Equal(t, match.HIDDENZONE, fireBird.Zone)
		assert.Equal(t, match.HIDDENZONE, earthDragon.Zone)

		attached := phoenix.Attachments()
		require.Len(t, attached, 2)
		assert.Equal(t, earthDragon.ID, attached[0].ID, "the chosen top base is attached first")
		assert.Equal(t, fireBird.ID, attached[1].ID)

		assert.False(t, phoenix.HasCondition(cnd.SummoningSickness), "vortex evolution has no summoning sickness, like a regular evolution")
	})

	t.Run("leaving the battle zone separates its bases back into the battle zone instead of taking them with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		fireBird := putCardInBattlezone(t, scn, player.Player, soulPhoenixFireBirdUID, soulPhoenixSetupSrc)
		earthDragon := putCardInBattlezone(t, scn, player.Player, soulPhoenixEarthDragonUID, soulPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		phoenix := summonSoulPhoenix(t, scn, player, fireBird.ID, earthDragon.ID)

		_, err := player.Player.MoveCard(phoenix.ID, match.BATTLEZONE, match.GRAVEYARD, soulPhoenixSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, phoenix.Zone, "only the top card actually leaves")
		assert.Equal(t, match.BATTLEZONE, fireBird.Zone, "its bases separate into the battle zone instead of following it")
		assert.Equal(t, match.BATTLEZONE, earthDragon.Zone)

		assert.False(t, fireBird.HasCondition(cnd.SummoningSickness), "a separated base is not newly summoned")
		assert.False(t, earthDragon.HasCondition(cnd.SummoningSickness))

		assert.Empty(t, phoenix.Attachments(), "nothing is left attached once they separate")
	})

	t.Run("leaving the battle zone through any destination still separates the bases", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		fireBird := putCardInBattlezone(t, scn, player.Player, soulPhoenixFireBirdUID, soulPhoenixSetupSrc)
		earthDragon := putCardInBattlezone(t, scn, player.Player, soulPhoenixEarthDragonUID, soulPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		phoenix := summonSoulPhoenix(t, scn, player, fireBird.ID, earthDragon.ID)

		_, err := player.Player.MoveCard(phoenix.ID, match.BATTLEZONE, match.HAND, soulPhoenixSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, phoenix.Zone)
		assert.Equal(t, match.BATTLEZONE, fireBird.Zone)
		assert.Equal(t, match.BATTLEZONE, earthDragon.Zone)
	})
}
