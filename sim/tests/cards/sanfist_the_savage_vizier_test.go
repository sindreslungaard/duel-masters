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
	sanfistTheSavageVizierUID       = "ad428b55-4b0b-4a63-b48b-dfc1541b8b81"
	sanfistTheSavageVizierMillerUID = "ba955ab0-5bb3-4aaf-82f3-293522e65a9c" // Locomotiver, discards on summon
	sanfistTheSavageVizierManaUID   = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul (darkness mana)
	sanfistTheSavageVizierSetupSrc  = "sanfist_the_savage_vizier_test_setup"
)

func TestSanfistTheSavageVizier(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		sanfist := putCardInBattlezone(t, scn, player.Player, sanfistTheSavageVizierUID, sanfistTheSavageVizierSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Sanfist, the Savage Vizier", sanfist.Name)
		assert.Equal(t, 3000, sanfist.Power)
		assert.Equal(t, 3, sanfist.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, sanfist.Civs)
		assert.Equal(t, []string{civ.Light, civ.Nature}, sanfist.ManaRequirement)
		assert.True(t, sanfist.IsMulticolored())
		assert.True(t, sanfist.HasFamily(family.BeastFolk))
		assert.True(t, sanfist.HasFamily(family.Initiate))
		assert.True(t, sanfist.HasCondition(cnd.Blocker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("a discard on the opponent's turn may put it into the battle zone instead", func(t *testing.T) {
		scn, player, opponent, sanfist := setupDiscardReplacementTest(t, sanfistTheSavageVizierUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))

		// The prompt belongs to Sanfist's controller, who is not the player
		// taking the turn.
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, sanfist.Zone, "it is put into play instead of discarded")
	})

	t.Run("declining discards it as normal", func(t *testing.T) {
		scn, player, opponent, sanfist := setupDiscardReplacementTest(t, sanfistTheSavageVizierUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, sanfist.Zone)
	})

	t.Run("it arrives with summoning sickness", func(t *testing.T) {
		scn, player, opponent, sanfist := setupDiscardReplacementTest(t, sanfistTheSavageVizierUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		require.Equal(t, match.BATTLEZONE, sanfist.Zone)
		assert.True(t, sanfist.HasCondition(cnd.SummoningSickness), "it was put into play, not summoned early")
	})

	t.Run("a discard on its controller's own turn is not replaced", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		sanfist, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Discarding during its controller's own turn: the printed trigger
		// only covers the opponent's turn, so nothing should be offered.
		_, err = player.Player.MoveCard(sanfist.ID, match.HAND, match.GRAVEYARD, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.GRAVEYARD, sanfist.Zone)
	})

	t.Run("it can block the very attack whose trigger discarded it", func(t *testing.T) {
		// Reported as a bug: Windmill Mutant's "whenever this creature
		// attacks" ability discards Sanfist, and Sanfist's controller is then
		// offered Sanfist as a blocker for that same attack. This is not a
		// timing error: attack-trigger abilities resolve fully during the
		// attack-declaration part of the sequence, before blockers are ever
		// chosen, so Sanfist is already an untapped creature in the battle
		// zone by the time blocking is decided (the same principle that lets
		// a Ninja Strike creature ambush-block the attack that let it in).
		scn, player, opponent := setupDuel(t)

		windmill := putCardInBattlezone(t, scn, player.Player, windmillMutantUID, sanfistTheSavageVizierSetupSrc)

		// One turn transition into the opponent's turn so they draw whatever
		// they're going to draw, then the hand is cleared and reseeded with
		// only Sanfist so the random discard can only ever hit it. Ending
		// that turn (rather than using passTurnToSelf) gives Sanfist its
		// intrinsic Blocker condition via the untap step without the
		// opponent drawing another card on top of it.
		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, scn.Match.IsPlayerTurn(opponent.Player))

		emptyHand(t, opponent, sanfistTheSavageVizierSetupSrc)
		sanfist, err := opponent.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, scn.Match.IsPlayerTurn(player.Player))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		questionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		action, err := scn.ActionAttackPlayer(player, windmill.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		// Windmill Mutant's attack trigger discards Sanfist (the only card in
		// hand), and Sanfist's controller is offered the discard replacement.
		_, err = scn.WaitForAction(opponent, questionStart)
		require.NoError(t, err, "expected Sanfist's discard-replacement question to be open")

		blockPromptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent))

		blockAction, err := scn.WaitForAction(opponent, blockPromptStart)
		require.NoError(t, err, "expected a blocker prompt offering Sanfist for this attack")

		offered := make([]string, 0, len(blockAction.Cards))
		for _, c := range blockAction.Cards {
			offered = append(offered, c.CardID)
		}
		assert.Contains(t, offered, sanfist.ID)

		require.NoError(t, scn.SubmitAction(opponent, sanfist.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, sanfist.Zone, "it survives the battle (3000 beats 2000)")
		assert.True(t, sanfist.Tapped, "blocking taps it")
		assert.Equal(t, match.GRAVEYARD, windmill.Zone, "the attacker loses the battle")

		remainingShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingShields, shieldCount, "a blocked attack breaks no shields")
	})

	t.Run("moving it out of hand for another reason is not a discard", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		sanfist, err := player.Player.SpawnCard(sanfistTheSavageVizierUID, match.HAND)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Charging it as mana on the opponent's turn is still a move out of
		// hand, but it is not a discard.
		_, err = player.Player.MoveCard(sanfist.ID, match.HAND, match.MANAZONE, sanfistTheSavageVizierSetupSrc)
		require.NoError(t, err)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.MANAZONE, sanfist.Zone)
	})
}
