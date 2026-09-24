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
	cruelNagaUID        = "23475e43-4cbd-4054-a90d-fde230a28db3"
	cruelNagaMerfolkUID = "acc8adb5-63c9-4438-976c-dcdf8fe1dad8" // Tide Patroller (Merfolk)
	cruelNagaChimeraUID = "5d73062e-acff-47e6-b49a-c0bb1a1762b5" // Gigagiele (Chimera)
	cruelNagaBlockerUID = "c7fec5e8-4e56-451b-a7b6-ad08680703a4" // La Byle, Seeker of the Winds (Blocker)
	cruelNagaFillerUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	cruelNagaSetupSrc   = "cruel_naga_avatar_of_fate_test_setup"
)

// summonCruelNaga pays for Cruel Naga with copies of itself and plays it, then
// answers the stacking-order prompt that appears once its controller has
// exactly one Merfolk and one Chimera: with no real choice for either
// requirement, both selections resolve on their own.
func summonCruelNaga(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, topID string, bottomID string) *match.Card {
	t.Helper()

	naga, err := player.Player.SpawnCard(cruelNagaUID, match.HAND)
	require.NoError(t, err)

	for range naga.ManaCost {
		_, err := player.Player.SpawnCard(cruelNagaUID, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, naga.ID))

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected the stacking-order prompt")
	require.Len(t, action.Cards, 2)

	require.NoError(t, scn.SubmitAction(player, topID, bottomID))
	require.NoError(t, scn.WaitForEventLoop())

	return naga
}

func TestCruelNagaAvatarOfFate(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		naga := putCardInBattlezone(t, scn, player.Player, cruelNagaUID, cruelNagaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, naga, "Cruel Naga, Avatar of Fate", 9000, 6, []string{civ.Water, civ.Darkness})
		assert.True(t, naga.HasFamily(family.Naga))
		assert.True(t, naga.IsMulticolored())
		assert.True(t, naga.HasCondition(cnd.DoubleBreaker))
		assert.True(t, naga.HasCondition(cnd.CantBeBlocked))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(cruelNagaUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, cruelNagaSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("cannot be summoned without a Merfolk and a Chimera in the battle zone", func(t *testing.T) {
		cases := []struct {
			name    string
			merfolk bool
			chimera bool
		}{
			{"neither", false, false},
			{"only a Merfolk", true, false},
			{"only a Chimera", false, true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				scn, player, _ := setupDuel(t)

				if c.merfolk {
					putCardInBattlezone(t, scn, player.Player, cruelNagaMerfolkUID, cruelNagaSetupSrc)
				}
				if c.chimera {
					putCardInBattlezone(t, scn, player.Player, cruelNagaChimeraUID, cruelNagaSetupSrc)
				}

				naga, err := player.Player.SpawnCard(cruelNagaUID, match.HAND)
				require.NoError(t, err)
				for range naga.ManaCost {
					_, err := player.Player.SpawnCard(cruelNagaUID, match.MANAZONE)
					require.NoError(t, err)
				}

				// Rejected before mana is even asked for, exactly like a
				// regular evolution with no legal base: no prompt is ever
				// opened for the caller to answer.
				require.Error(t, scn.ActionPlayCard(player, naga.ID), "there is nothing for it to evolve from")
				assert.Equal(t, match.HAND, naga.Zone, "it stays in hand")
			})
		}
	})

	t.Run("vortex evolution consumes one of each and stacks on top of the chosen base", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		merfolk := putCardInBattlezone(t, scn, player.Player, cruelNagaMerfolkUID, cruelNagaSetupSrc)
		chimera := putCardInBattlezone(t, scn, player.Player, cruelNagaChimeraUID, cruelNagaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Tapped, and chosen as the top of the stack: Cruel Naga should take
		// this tap state, matching a regular evolution landing on a tapped
		// creature.
		chimera.Tapped = true

		naga := summonCruelNaga(t, scn, player, chimera.ID, merfolk.ID)

		require.Equal(t, match.BATTLEZONE, naga.Zone)
		assert.Equal(t, 9000, naga.Power)
		assert.True(t, naga.Tapped, "it took the tap state of the base it was put on top of")

		assert.Equal(t, match.HIDDENZONE, merfolk.Zone)
		assert.Equal(t, match.HIDDENZONE, chimera.Zone)

		attached := naga.Attachments()
		require.Len(t, attached, 2)
		assert.Equal(t, chimera.ID, attached[0].ID, "the chosen top base is attached first")
		assert.Equal(t, merfolk.ID, attached[1].ID)

		assert.False(t, naga.HasCondition(cnd.SummoningSickness), "vortex evolution has no summoning sickness, like a regular evolution")
	})

	t.Run("it cannot be blocked", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		naga := putCardInBattlezone(t, scn, player.Player, cruelNagaUID, cruelNagaSetupSrc)
		blocker := putCardInBattlezone(t, scn, opponent.Player, cruelNagaBlockerUID, cruelNagaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)
		require.True(t, blocker.HasCondition(cnd.Blocker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, naga.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))
		settleTurn(t, scn)

		assert.False(t, blocker.Tapped, "it was never offered as a blocker")

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-2, "the attack went through unblocked")
	})

	t.Run("leaving the battle zone destroys every creature in the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		merfolk := putCardInBattlezone(t, scn, player.Player, cruelNagaMerfolkUID, cruelNagaSetupSrc)
		chimera := putCardInBattlezone(t, scn, player.Player, cruelNagaChimeraUID, cruelNagaSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, cruelNagaFillerUID, cruelNagaSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, cruelNagaFillerUID, cruelNagaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		naga := summonCruelNaga(t, scn, player, merfolk.ID, chimera.ID)

		_, err := player.Player.MoveCard(naga.ID, match.BATTLEZONE, match.GRAVEYARD, cruelNagaSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, naga.Zone)
		assert.Equal(t, match.GRAVEYARD, merfolk.Zone, "its base leaves the hidden zone together with it")
		assert.Equal(t, match.GRAVEYARD, chimera.Zone, "so does the other base")
		assert.Equal(t, match.GRAVEYARD, ally.Zone, "destroy all creatures catches its controller's own")
		assert.Equal(t, match.GRAVEYARD, theirs.Zone, "and the opponent's")
	})
}
