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
	auraPegasusUID             = "5ffd2e9f-98bc-4e36-8454-8f822740f8eb"
	auraPegasusHornedBeastUID  = "18e0e199-7827-4a4c-a37d-3acfa4e500d6" // Roaring Great-Horn (Horned Beast)
	auraPegasusAngelCommandUID = "5d3d7052-e5fa-4502-8d31-c72673232317" // Hanusa, Radiance Elemental (Angel Command)
	auraPegasusNonEvolutionUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (plain creature)
	auraPegasusEvolutionUID    = "24353d06-89ef-4867-9513-485750d01e10" // Armored Cannon Balbaro (evolves from Human)
	auraPegasusSetupSrc        = "aura_pegasus_avatar_of_life_test_setup"
)

// summonAuraPegasus pays for Aura Pegasus with copies of itself and plays it,
// then answers the stacking-order prompt that appears once its controller has
// exactly one Horned Beast and one Angel Command: with no real choice for
// either requirement, both selections resolve on their own.
func summonAuraPegasus(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, topID string, bottomID string) *match.Card {
	t.Helper()

	pegasus, err := player.Player.SpawnCard(auraPegasusUID, match.HAND)
	require.NoError(t, err)

	for range pegasus.ManaCost {
		_, err := player.Player.SpawnCard(auraPegasusUID, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, pegasus.ID))

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected the stacking-order prompt")
	require.Len(t, action.Cards, 2)

	require.NoError(t, scn.SubmitAction(player, topID, bottomID))
	require.NoError(t, scn.WaitForEventLoop())

	return pegasus
}

// seedTopOfDeck moves card, already in hand, to the front of player's deck,
// so a reveal-the-top-card effect is deterministic.
//
// It must already be in hand rather than spawned fresh here: cnd.Creature and
// cnd.Evolution are only rebuilt for a card that is actually present -
// whatever its zone - when an untap step runs, so a card conjured straight
// into the deck after the untap step this test already passed through would
// carry neither, unlike every card that has been part of the deck since the
// match began.
func seedTopOfDeck(t *testing.T, player *match.PlayerReference, card *match.Card) *match.Card {
	t.Helper()

	top, err := player.Player.MoveCardToFront(card.ID, match.HAND, match.DECK, auraPegasusSetupSrc)
	require.NoError(t, err)

	deck, err := player.Player.Container(match.DECK)
	require.NoError(t, err)
	require.Equal(t, top.ID, deck[0].ID, "the seeded card must be on top")

	return top
}

func TestAuraPegasusAvatarOfLife(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pegasus := putCardInBattlezone(t, scn, player.Player, auraPegasusUID, auraPegasusSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, pegasus, "Aura Pegasus, Avatar of Life", 12000, 6, []string{civ.Light, civ.Nature})
		assert.True(t, pegasus.HasFamily(family.Pegasus))
		assert.True(t, pegasus.IsMulticolored())
		assert.True(t, pegasus.HasCondition(cnd.TripleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(auraPegasusUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, auraPegasusSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("cannot be summoned without a Horned Beast and an Angel Command in the battle zone", func(t *testing.T) {
		cases := []struct {
			name         string
			hornedBeast  bool
			angelCommand bool
		}{
			{"neither", false, false},
			{"only a Horned Beast", true, false},
			{"only an Angel Command", false, true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				scn, player, _ := setupDuel(t)

				if c.hornedBeast {
					putCardInBattlezone(t, scn, player.Player, auraPegasusHornedBeastUID, auraPegasusSetupSrc)
				}
				if c.angelCommand {
					putCardInBattlezone(t, scn, player.Player, auraPegasusAngelCommandUID, auraPegasusSetupSrc)
				}

				pegasus, err := player.Player.SpawnCard(auraPegasusUID, match.HAND)
				require.NoError(t, err)
				for range pegasus.ManaCost {
					_, err := player.Player.SpawnCard(auraPegasusUID, match.MANAZONE)
					require.NoError(t, err)
				}

				// Rejected before mana is even asked for, exactly like a
				// regular evolution with no legal base: no prompt is ever
				// opened for the caller to answer.
				require.Error(t, scn.ActionPlayCard(player, pegasus.ID), "there is nothing for it to evolve from")
				assert.Equal(t, match.HAND, pegasus.Zone, "it stays in hand")
			})
		}
	})

	t.Run("vortex evolution consumes one of each and stacks on top of the chosen base", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hornedBeast := putCardInBattlezone(t, scn, player.Player, auraPegasusHornedBeastUID, auraPegasusSetupSrc)
		angelCommand := putCardInBattlezone(t, scn, player.Player, auraPegasusAngelCommandUID, auraPegasusSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Tapped, and chosen as the top of the stack: Aura Pegasus should take
		// this tap state, matching a regular evolution landing on a tapped
		// creature.
		angelCommand.Tapped = true

		pegasus := summonAuraPegasus(t, scn, player, angelCommand.ID, hornedBeast.ID)

		require.Equal(t, match.BATTLEZONE, pegasus.Zone)
		assert.Equal(t, 12000, pegasus.Power)
		assert.True(t, pegasus.Tapped, "it took the tap state of the base it was put on top of")

		assert.Equal(t, match.HIDDENZONE, hornedBeast.Zone)
		assert.Equal(t, match.HIDDENZONE, angelCommand.Zone)

		attached := pegasus.Attachments()
		require.Len(t, attached, 2)
		assert.Equal(t, angelCommand.ID, attached[0].ID, "the chosen top base is attached first")
		assert.Equal(t, hornedBeast.ID, attached[1].ID)

		assert.False(t, pegasus.HasCondition(cnd.SummoningSickness), "vortex evolution has no summoning sickness, like a regular evolution")
	})

	t.Run("attacking reveals the top card: a non-evolution creature goes to the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pegasus := putCardInBattlezone(t, scn, player.Player, auraPegasusUID, auraPegasusSetupSrc)
		seed := spawnForLater(t, player, auraPegasusNonEvolutionUID)
		passTurnToSelf(t, scn, player, opponent)
		top := seedTopOfDeck(t, player, seed)

		action, err := scn.ActionAttackPlayer(player, pegasus.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, top.Zone)
	})

	t.Run("attacking reveals the top card: an evolution creature goes to hand instead", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pegasus := putCardInBattlezone(t, scn, player.Player, auraPegasusUID, auraPegasusSetupSrc)
		seed := spawnForLater(t, player, auraPegasusEvolutionUID)
		passTurnToSelf(t, scn, player, opponent)
		require.True(t, seed.HasCondition(cnd.Evolution))
		top := seedTopOfDeck(t, player, seed)

		action, err := scn.ActionAttackPlayer(player, pegasus.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, top.Zone, "an evolution creature is not a legal target for this effect")
	})

	t.Run("leaving the battle zone also reveals the top card", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hornedBeast := putCardInBattlezone(t, scn, player.Player, auraPegasusHornedBeastUID, auraPegasusSetupSrc)
		angelCommand := putCardInBattlezone(t, scn, player.Player, auraPegasusAngelCommandUID, auraPegasusSetupSrc)
		seed := spawnForLater(t, player, auraPegasusNonEvolutionUID)
		passTurnToSelf(t, scn, player, opponent)

		pegasus := summonAuraPegasus(t, scn, player, hornedBeast.ID, angelCommand.ID)
		top := seedTopOfDeck(t, player, seed)

		_, err := player.Player.MoveCard(pegasus.ID, match.BATTLEZONE, match.GRAVEYARD, auraPegasusSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, pegasus.Zone)
		assert.Equal(t, match.GRAVEYARD, hornedBeast.Zone, "its base leaves the hidden zone together with it")
		assert.Equal(t, match.GRAVEYARD, angelCommand.Zone, "so does the other base")
		assert.Equal(t, match.BATTLEZONE, top.Zone, "leaving the battle zone also reveals the top card")
	})
}
