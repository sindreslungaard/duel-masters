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
	deathPhoenixUID             = "dcc5ee70-bc60-420d-92e5-ed8bd7ff949f"
	deathPhoenixZombieDragonUID = "d2f91d1a-5e8d-43ce-8512-ec6be9c3f424" // Necrodragon Giland (Zombie Dragon)
	deathPhoenixFireBirdUID     = "479b477f-f535-4564-ac0c-5f0aeaafe914" // Baby Zoppe (Fire Bird)
	deathPhoenixFillerUID       = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	deathPhoenixSetupSrc        = "death_phoenix_avatar_of_doom_test_setup"
)

// summonDeathPhoenix pays for Death Phoenix with copies of itself and plays
// it, then answers the stacking-order prompt that appears once its controller
// has exactly one Zombie Dragon and one Fire Bird: with no real choice for
// either requirement, both selections resolve on their own.
func summonDeathPhoenix(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, topID string, bottomID string) *match.Card {
	t.Helper()

	phoenix, err := player.Player.SpawnCard(deathPhoenixUID, match.HAND)
	require.NoError(t, err)

	for range phoenix.ManaCost {
		_, err := player.Player.SpawnCard(deathPhoenixUID, match.MANAZONE)
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

func TestDeathPhoenixAvatarOfDoom(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phoenix := putCardInBattlezone(t, scn, player.Player, deathPhoenixUID, deathPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, phoenix, "Death Phoenix, Avatar of Doom", 9000, 4, []string{civ.Darkness, civ.Fire})
		assert.True(t, phoenix.HasFamily(family.Phoenix))
		assert.True(t, phoenix.IsMulticolored())
		assert.True(t, phoenix.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(deathPhoenixUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, deathPhoenixSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("cannot be summoned without a Zombie Dragon and a Fire Bird in the battle zone", func(t *testing.T) {
		cases := []struct {
			name         string
			zombieDragon bool
			fireBird     bool
		}{
			{"neither", false, false},
			{"only a Zombie Dragon", true, false},
			{"only a Fire Bird", false, true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				scn, player, _ := setupDuel(t)

				if c.zombieDragon {
					putCardInBattlezone(t, scn, player.Player, deathPhoenixZombieDragonUID, deathPhoenixSetupSrc)
				}
				if c.fireBird {
					putCardInBattlezone(t, scn, player.Player, deathPhoenixFireBirdUID, deathPhoenixSetupSrc)
				}

				phoenix, err := player.Player.SpawnCard(deathPhoenixUID, match.HAND)
				require.NoError(t, err)
				for range phoenix.ManaCost {
					_, err := player.Player.SpawnCard(deathPhoenixUID, match.MANAZONE)
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
		zombieDragon := putCardInBattlezone(t, scn, player.Player, deathPhoenixZombieDragonUID, deathPhoenixSetupSrc)
		fireBird := putCardInBattlezone(t, scn, player.Player, deathPhoenixFireBirdUID, deathPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Tapped, and chosen as the top of the stack: Death Phoenix should
		// take this tap state, matching a regular evolution landing on a
		// tapped creature.
		fireBird.Tapped = true

		phoenix := summonDeathPhoenix(t, scn, player, fireBird.ID, zombieDragon.ID)

		require.Equal(t, match.BATTLEZONE, phoenix.Zone)
		assert.Equal(t, 9000, phoenix.Power)
		assert.True(t, phoenix.Tapped, "it took the tap state of the base it was put on top of")

		assert.Equal(t, match.HIDDENZONE, zombieDragon.Zone)
		assert.Equal(t, match.HIDDENZONE, fireBird.Zone)

		attached := phoenix.Attachments()
		require.Len(t, attached, 2)
		assert.Equal(t, fireBird.ID, attached[0].ID, "the chosen top base is attached first")
		assert.Equal(t, zombieDragon.ID, attached[1].ID)

		assert.False(t, phoenix.HasCondition(cnd.SummoningSickness), "vortex evolution has no summoning sickness, like a regular evolution")
	})

	t.Run("shields it would break go to the opponent's graveyard instead of hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		phoenix := putCardInBattlezone(t, scn, player.Player, deathPhoenixUID, deathPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		action, err := scn.ActionAttackPlayer(player, phoenix.ID)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(action.Cards), 2)
		targeted := []string{action.Cards[0].CardID, action.Cards[1].CardID}

		require.NoError(t, scn.ResolveAttack(player, targeted...))
		settleTurn(t, scn)

		for _, id := range targeted {
			card, err := opponent.Player.GetCard(id, match.GRAVEYARD)
			require.NoError(t, err, "the shield landed in the graveyard, not the hand")
			assert.Equal(t, match.GRAVEYARD, card.Zone)
		}

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount, "the broken shields never reached hand")
	})

	t.Run("leaving the battle zone makes the opponent discard their hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		zombieDragon := putCardInBattlezone(t, scn, player.Player, deathPhoenixZombieDragonUID, deathPhoenixSetupSrc)
		fireBird := putCardInBattlezone(t, scn, player.Player, deathPhoenixFireBirdUID, deathPhoenixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		phoenix := summonDeathPhoenix(t, scn, player, zombieDragon.ID, fireBird.ID)

		emptyHand(t, opponent, deathPhoenixSetupSrc)
		_, err := opponent.Player.SpawnCard(deathPhoenixFillerUID, match.HAND)
		require.NoError(t, err)
		_, err = opponent.Player.SpawnCard(deathPhoenixFillerUID, match.HAND)
		require.NoError(t, err)

		_, err = player.Player.MoveCard(phoenix.ID, match.BATTLEZONE, match.GRAVEYARD, deathPhoenixSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand, "the opponent's whole hand was discarded")

		assert.Equal(t, match.GRAVEYARD, phoenix.Zone)
		assert.Equal(t, match.GRAVEYARD, zombieDragon.Zone, "its base leaves the hidden zone together with it")
		assert.Equal(t, match.GRAVEYARD, fireBird.Zone, "so does the other base")
	})
}
