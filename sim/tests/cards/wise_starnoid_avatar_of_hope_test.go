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
	wiseStarnoidUID              = "ba894f7e-b7d7-409e-8393-cab4285a879c"
	wiseStarnoidLightBringerUID  = "7b58e8c2-0b1e-4ef5-812f-e667c2092c73" // Reusol, the Oracle (vanilla Light Bringer, 2000)
	wiseStarnoidLightBringer2UID = "616c146e-049f-4720-a225-0a189729ca79" // Chilias, the Oracle (another Light Bringer)
	wiseStarnoidCyberLordUID     = "7a6f1c82-a8ac-4646-b3e9-fb8592bdd0a4" // Tropico (Cyber Lord, 3000)
	wiseStarnoidSetupSrc         = "wise_starnoid_avatar_of_hope_test_setup"
)

// summonWiseStarnoid pays for Wise Starnoid with copies of itself and plays
// it, then answers the single stacking-order prompt that appears when its
// controller has exactly one Light Bringer and one Cyber Lord: with no real
// choice for either requirement, both selections resolve on their own and
// only the ordering (topID goes on top, then bottomID) is actually asked.
func summonWiseStarnoid(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, topID string, bottomID string) *match.Card {
	t.Helper()

	starnoid, err := player.Player.SpawnCard(wiseStarnoidUID, match.HAND)
	require.NoError(t, err)

	for range starnoid.ManaCost {
		_, err := player.Player.SpawnCard(wiseStarnoidUID, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, starnoid.ID))

	orderPrompt, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected the stacking-order prompt")
	require.Len(t, orderPrompt.Cards, 2)

	require.NoError(t, scn.SubmitAction(player, topID, bottomID))
	require.NoError(t, scn.WaitForEventLoop())

	return starnoid
}

func TestWiseStarnoidAvatarOfHope(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		starnoid := putCardInBattlezone(t, scn, player.Player, wiseStarnoidUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, starnoid, "Wise Starnoid, Avatar of Hope", 9000, 5, []string{civ.Light, civ.Water})
		assert.True(t, starnoid.HasFamily(family.Starnoid))
		assert.True(t, starnoid.IsMulticolored())
		assert.True(t, starnoid.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(wiseStarnoidUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, wiseStarnoidSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("cannot be summoned without a Light Bringer and a Cyber Lord in the battle zone", func(t *testing.T) {
		cases := []struct {
			name       string
			lightBring bool
			cyberLord  bool
		}{
			{"neither", false, false},
			{"only a Light Bringer", true, false},
			{"only a Cyber Lord", false, true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				scn, player, _ := setupDuel(t)

				if c.lightBring {
					putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringerUID, wiseStarnoidSetupSrc)
				}
				if c.cyberLord {
					putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)
				}

				starnoid, err := player.Player.SpawnCard(wiseStarnoidUID, match.HAND)
				require.NoError(t, err)
				for range starnoid.ManaCost {
					_, err := player.Player.SpawnCard(wiseStarnoidUID, match.MANAZONE)
					require.NoError(t, err)
				}

				// Rejected before mana is even asked for, exactly like a
				// regular evolution with no legal base: no prompt is ever
				// opened for the caller to answer.
				require.Error(t, scn.ActionPlayCard(player, starnoid.ID), "there is nothing for it to evolve from")
				assert.Equal(t, match.HAND, starnoid.Zone, "it stays in hand")
			})
		}
	})

	t.Run("a single creature cannot fill both roles by itself", func(t *testing.T) {
		// Tropico is only a Cyber Lord, not a Light Bringer, but this proves
		// the engine requires two distinct creatures rather than merely
		// checking that both characteristics are present somewhere.
		scn, player, _ := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)

		starnoid, err := player.Player.SpawnCard(wiseStarnoidUID, match.HAND)
		require.NoError(t, err)
		for range starnoid.ManaCost {
			_, err := player.Player.SpawnCard(wiseStarnoidUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.Error(t, scn.ActionPlayCard(player, starnoid.ID))
		assert.Equal(t, match.HAND, starnoid.Zone)
	})

	t.Run("vortex evolution consumes one of each and stacks on top of the chosen base", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lightBringer := putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringerUID, wiseStarnoidSetupSrc)
		cyberLord := putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Tapped, and chosen as the top of the stack: Wise Starnoid should
		// take this tap state, matching a regular evolution landing on a
		// tapped creature.
		cyberLord.Tapped = true

		starnoid := summonWiseStarnoid(t, scn, player, cyberLord.ID, lightBringer.ID)

		require.Equal(t, match.BATTLEZONE, starnoid.Zone)
		assert.Equal(t, 9000, starnoid.Power)
		assert.True(t, starnoid.Tapped, "it took the tap state of the base it was put on top of")

		assert.Equal(t, match.HIDDENZONE, lightBringer.Zone)
		assert.Equal(t, match.HIDDENZONE, cyberLord.Zone)

		attached := starnoid.Attachments()
		require.Len(t, attached, 2)
		assert.Equal(t, cyberLord.ID, attached[0].ID, "the chosen top base is attached first")
		assert.Equal(t, lightBringer.ID, attached[1].ID)

		assert.False(t, starnoid.HasCondition(cnd.SummoningSickness), "vortex evolution has no summoning sickness, like a regular evolution")
	})

	t.Run("the controller chooses which creature to use when more than one qualifies", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		wanted := putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringerUID, wiseStarnoidSetupSrc)
		unwanted := putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringer2UID, wiseStarnoidSetupSrc)
		cyberLord := putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		starnoid, err := player.Player.SpawnCard(wiseStarnoidUID, match.HAND)
		require.NoError(t, err)
		for range starnoid.ManaCost {
			_, err := player.Player.SpawnCard(wiseStarnoidUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionPlayCard(player, starnoid.ID))

		lbPrompt, err := scn.LatestAction(player, 0)
		require.NoError(t, err, "expected a choice between the two Light Bringers")
		offered := make([]string, 0, len(lbPrompt.Cards))
		for _, c := range lbPrompt.Cards {
			offered = append(offered, c.CardID)
		}
		assert.ElementsMatch(t, []string{wanted.ID, unwanted.ID}, offered)

		orderPromptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, wanted.ID))

		// Only one Cyber Lord exists, so that selection resolves on its own
		// and the very next prompt is the stacking order.
		_, err = scn.WaitForAction(player, orderPromptStart)
		require.NoError(t, err, "expected the stacking-order prompt next")

		require.NoError(t, scn.SubmitAction(player, wanted.ID, cyberLord.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HIDDENZONE, wanted.Zone, "the chosen Light Bringer became a base")
		assert.Equal(t, match.BATTLEZONE, unwanted.Zone, "the other one was left alone")
	})

	t.Run("it can attack the same turn it evolves", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lightBringer := putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringerUID, wiseStarnoidSetupSrc)
		cyberLord := putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		starnoid := summonWiseStarnoid(t, scn, player, lightBringer.ID, cyberLord.ID)

		_, err := scn.ActionAttackPlayer(player, starnoid.ID)
		require.NoError(t, err, "the attack should open the shield-selection prompt, not be rejected for summoning sickness")

		require.NoError(t, scn.CancelAction(player))
	})

	t.Run("attacking adds the top card of the deck to its controller's shields face-down", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		starnoid := putCardInBattlezone(t, scn, player.Player, wiseStarnoidUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		require.NotEmpty(t, deckBefore)
		topOfDeck := deckBefore[0]

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		action, err := scn.ActionAttackPlayer(player, starnoid.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))
		settleTurn(t, scn)

		assert.Equal(t, match.SHIELDZONE, topOfDeck.Zone, "the top card of the deck became a shield")

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, len(deckBefore)-1)
	})

	t.Run("leaving the battle zone adds a shield and takes both bases with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lightBringer := putCardInBattlezone(t, scn, player.Player, wiseStarnoidLightBringerUID, wiseStarnoidSetupSrc)
		cyberLord := putCardInBattlezone(t, scn, player.Player, wiseStarnoidCyberLordUID, wiseStarnoidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		starnoid := summonWiseStarnoid(t, scn, player, lightBringer.ID, cyberLord.ID)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		require.NotEmpty(t, deckBefore)
		topOfDeck := deckBefore[0]

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = player.Player.MoveCard(starnoid.ID, match.BATTLEZONE, match.GRAVEYARD, wiseStarnoidSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, starnoid.Zone)
		assert.Equal(t, match.GRAVEYARD, lightBringer.Zone, "its base leaves the hidden zone together with it")
		assert.Equal(t, match.GRAVEYARD, cyberLord.Zone, "so does the other base")

		assert.Equal(t, match.SHIELDZONE, topOfDeck.Zone, "leaving the battle zone also adds the top card of the deck as a shield")

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)
	})
}
