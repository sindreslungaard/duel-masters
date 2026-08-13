package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	iceVaporShadowOfAnguishUID      = "ab6c7559-1714-4238-a063-393cfe8adc08"
	iceVaporShadowOfAnguishSpellUID = "b7f236fd-e7eb-41cc-912a-5239c134f265" // Energy Stream
	iceVaporShadowOfAnguishManaUID  = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	iceVaporShadowOfAnguishSetupSrc = "ice_vapor_shadow_of_anguish_test_setup"
)

func TestIceVaporShadowOfAnguish(t *testing.T) {
	t.Run("makes the opponent discard and burn mana when they cast a spell", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		spell := prepareIceVaporShadowOfAnguishTestSpell(t, scn, opponent)

		assert.Equal(t, "Ice Vapor, Shadow of Anguish", iceVapor.Name)
		assert.Equal(t, 1000, iceVapor.Power)
		assert.Equal(t, 5, iceVapor.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, iceVapor.Civs)
		assert.True(t, iceVapor.HasFamily(family.Ghost))

		require.NoError(t, scn.ActionEndTurn(owner))

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(opponent, spell.ID))

		discardAction, err := scn.LatestAction(opponent, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, discardAction.MinSelections, "the discard is mandatory")
		assert.Equal(t, 1, discardAction.MaxSelections)
		assert.False(t, discardAction.Cancellable)

		discarded, err := opponent.Player.GetCard(discardAction.Cards[0].CardID, match.HAND)
		require.NoError(t, err)

		manaPromptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, discarded.ID))

		manaAction, err := scn.WaitForAction(opponent, manaPromptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, manaAction.MinSelections)
		assert.Equal(t, 1, manaAction.MaxSelections)
		assert.False(t, manaAction.Cancellable)

		burned, err := opponent.Player.GetCard(manaAction.Cards[0].CardID, match.MANAZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, burned.ID))
		require.NoError(t, scn.WaitForMessage(opponent, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, discarded.Zone)
		assert.Equal(t, match.GRAVEYARD, burned.Zone)
	})

	t.Run("does not trigger on its own controller's spell", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		spell := prepareIceVaporShadowOfAnguishTestSpell(t, scn, owner)

		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		ownerHandBefore, err := owner.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(ownerHandBefore)
		ownerManaBefore, err := owner.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		manaCount := len(ownerManaBefore)

		require.NoError(t, scn.ActionPlayCard(owner, spell.ID))

		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "action")

		ownerHandAfter, err := owner.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, ownerHandAfter, handCount-1+2, "only Energy Stream's own draw changed the hand")
		ownerManaAfter, err := owner.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, ownerManaAfter, manaCount)
	})

	t.Run("does not trigger from outside the battle zone", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		owner.Player.SpawnCard(iceVaporShadowOfAnguishUID, match.HAND)
		spell := prepareIceVaporShadowOfAnguishTestSpell(t, scn, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))

		opponentManaBefore, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		manaCount := len(opponentManaBefore)

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(opponent, spell.ID))

		headers, err := scn.MessageHeaders(opponent, promptStart)
		require.NoError(t, err)
		actions := 0
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "only the mana payment is prompted")

		opponentManaAfter, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, opponentManaAfter, manaCount)
	})

	t.Run("still burns mana when the opponent has no cards left in hand", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		spell := prepareIceVaporShadowOfAnguishTestSpell(t, scn, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))
		clearIceVaporShadowOfAnguishTestHand(t, opponent.Player, spell.ID)

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(opponent, spell.ID))

		// Energy Stream refills the hand before Ice Vapor resolves, so the discard
		// is still offered; the mana burn is what must survive either way.
		discardAction, err := scn.LatestAction(opponent, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, discardAction.Cards)

		manaPromptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, discardAction.Cards[0].CardID))

		manaAction, err := scn.WaitForAction(opponent, manaPromptStart)
		require.NoError(t, err)
		burned, err := opponent.Player.GetCard(manaAction.Cards[0].CardID, match.MANAZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, burned.ID))
		require.NoError(t, scn.WaitForMessage(opponent, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, burned.Zone)
	})
}

func setupIceVaporShadowOfAnguishTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}

func putIceVaporShadowOfAnguishTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player) *match.Card {
	t.Helper()

	player.SpawnCard(iceVaporShadowOfAnguishUID, match.HAND)
	card, err := scn.FindCard(player, match.HAND, iceVaporShadowOfAnguishUID)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, iceVaporShadowOfAnguishSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

func prepareIceVaporShadowOfAnguishTestSpell(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	player.Player.SpawnCard(iceVaporShadowOfAnguishSpellUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(iceVaporShadowOfAnguishManaUID, match.MANAZONE)
	}

	spell, err := scn.FindCard(player.Player, match.HAND, iceVaporShadowOfAnguishSpellUID)
	require.NoError(t, err)
	return spell
}

func clearIceVaporShadowOfAnguishTestHand(t *testing.T, player *match.Player, keepID string) {
	t.Helper()

	hand, err := player.Container(match.HAND)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), hand...) {
		if card.ID == keepID {
			continue
		}

		moved, err := player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, iceVaporShadowOfAnguishSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
