package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/server"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	iceVaporShadowOfAnguishUID        = "ab6c7559-1714-4238-a063-393cfe8adc08"
	iceVaporShadowOfAnguishSpellUID   = "b7f236fd-e7eb-41cc-912a-5239c134f265" // Energy Stream
	iceVaporShadowOfAnguishManaUID    = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	iceVaporShadowOfAnguishRemovalUID = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit (shield trigger)
	iceVaporShadowOfAnguishSetupSrc   = "ice_vapor_shadow_of_anguish_test_setup"
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

	t.Run("still resolves when that same spell destroys it", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		// A second creature so the removal spell has to ask which one to destroy.
		putCardInBattlezone(t, scn, owner.Player, iceVaporShadowOfAnguishManaUID, iceVaporShadowOfAnguishSetupSrc)
		removal := prepareIceVaporShadowOfAnguishTestRemoval(t, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(opponent, removal.ID))

		targetAction, err := scn.LatestAction(opponent, promptStart)
		require.NoError(t, err)
		assert.True(t, iceVaporShadowOfAnguishTestOffers(targetAction, iceVapor.ID), "the removal spell can pick Ice Vapor")

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, iceVapor.ID))

		// The ability triggered while Ice Vapor was still in the battle zone, so it
		// resolves independently of its source being destroyed by that spell.
		discarded, burned, _ := answerIceVaporShadowOfAnguishTrigger(t, scn, opponent, triggerStart)

		assert.Equal(t, match.GRAVEYARD, iceVapor.Zone)
		assert.Equal(t, match.GRAVEYARD, discarded.Zone)
		assert.Equal(t, match.GRAVEYARD, burned.Zone)
	})

	t.Run("burns mana without a discard when the destroyed source leaves an empty hand", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		putCardInBattlezone(t, scn, owner.Player, iceVaporShadowOfAnguishManaUID, iceVaporShadowOfAnguishSetupSrc)
		removal := prepareIceVaporShadowOfAnguishTestRemoval(t, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))
		clearIceVaporShadowOfAnguishTestHand(t, opponent.Player, removal.ID)

		require.NoError(t, scn.ActionPlayCard(opponent, removal.ID))

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, iceVapor.ID))

		// Terror Pit hands back no cards, so there is nothing to discard and only
		// the mana burn is prompted.
		manaAction, err := scn.WaitForAction(opponent, triggerStart)
		require.NoError(t, err)
		require.NotEmpty(t, manaAction.Cards)
		burned, err := opponent.Player.GetCard(manaAction.Cards[0].CardID, match.MANAZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, burned.ID))
		require.NoError(t, scn.WaitForMessage(opponent, completionStart, "state_update"))

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand, "there was never a card to discard")
		assert.Equal(t, match.GRAVEYARD, iceVapor.Zone)
		assert.Equal(t, match.GRAVEYARD, burned.Zone)
	})

	t.Run("resolves once per copy when the spell destroys one of two", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		first := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		second := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		removal := prepareIceVaporShadowOfAnguishTestRemoval(t, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))

		require.NoError(t, scn.ActionPlayCard(opponent, removal.ID))

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, first.ID))

		firstDiscard, firstBurn, nextStart := answerIceVaporShadowOfAnguishTrigger(t, scn, opponent, triggerStart)
		secondDiscard, secondBurn, _ := answerIceVaporShadowOfAnguishTrigger(t, scn, opponent, nextStart)

		assert.Equal(t, match.GRAVEYARD, first.Zone, "the destroyed copy triggered too")
		assert.Equal(t, match.BATTLEZONE, second.Zone)
		assert.NotEqual(t, firstDiscard.ID, secondDiscard.ID)
		assert.NotEqual(t, firstBurn.ID, secondBurn.ID)
		assert.Equal(t, match.GRAVEYARD, firstDiscard.Zone)
		assert.Equal(t, match.GRAVEYARD, secondDiscard.Zone)
		assert.Equal(t, match.GRAVEYARD, firstBurn.Zone)
		assert.Equal(t, match.GRAVEYARD, secondBurn.Zone)
	})

	t.Run("does not trigger for a spell cast after it left the battle zone", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		spell := prepareIceVaporShadowOfAnguishTestSpell(t, scn, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))

		moved, err := owner.Player.MoveCard(iceVapor.ID, match.BATTLEZONE, match.GRAVEYARD, iceVaporShadowOfAnguishSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		manaBefore, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(opponent, spell.ID))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(opponent, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment is prompted")

		manaAfter, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, manaAfter, len(manaBefore))
	})

	t.Run("triggers again after returning to the battle zone", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		bystander := putCardInBattlezone(t, scn, owner.Player, iceVaporShadowOfAnguishManaUID, iceVaporShadowOfAnguishSetupSrc)
		removal := prepareIceVaporShadowOfAnguishTestRemoval(t, opponent)

		require.NoError(t, scn.ActionEndTurn(owner))

		moved, err := owner.Player.MoveCard(iceVapor.ID, match.BATTLEZONE, match.GRAVEYARD, iceVaporShadowOfAnguishSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		back, err := owner.Player.MoveCard(iceVapor.ID, match.GRAVEYARD, match.BATTLEZONE, iceVaporShadowOfAnguishSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, back.Zone)

		require.NoError(t, scn.ActionPlayCard(opponent, removal.ID))

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, bystander.ID))

		discarded, burned, _ := answerIceVaporShadowOfAnguishTrigger(t, scn, opponent, triggerStart)

		assert.Equal(t, match.GRAVEYARD, bystander.Zone)
		assert.Equal(t, match.BATTLEZONE, iceVapor.Zone, "the reinstalled trigger fired from the battle zone")
		assert.Equal(t, match.GRAVEYARD, discarded.Zone)
		assert.Equal(t, match.GRAVEYARD, burned.Zone)
	})

	t.Run("resolves for a shield trigger spell that destroys it during its controller's turn", func(t *testing.T) {
		scn, owner, opponent := setupIceVaporShadowOfAnguishTest(t)
		iceVapor := putIceVaporShadowOfAnguishTestCardInBattlezone(t, scn, owner.Player)
		shield := putIceVaporShadowOfAnguishTestRemovalInShields(t, opponent.Player)
		for range 3 {
			opponent.Player.SpawnCard(iceVaporShadowOfAnguishManaUID, match.MANAZONE)
		}

		// A turn each way so the untap steps give the shield its "shield trigger".
		passTurnToSelf(t, scn, owner, opponent)
		require.True(t, scn.Match.IsPlayerTurn(owner.Player))

		_, err := scn.ActionAttackPlayer(owner, iceVapor.ID)
		require.NoError(t, err)

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(owner, shield.ID))

		shieldTriggerAction, err := scn.WaitForAction(opponent, triggerStart)
		require.NoError(t, err)
		require.True(t, iceVaporShadowOfAnguishTestOffers(shieldTriggerAction, shield.ID), "the broken shield is offered")

		castStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shield.ID))

		// Ice Vapor is the only creature the cast spell can pick, so it is
		// destroyed without a prompt, and the trigger still resolves afterwards.
		discarded, burned, _ := answerIceVaporShadowOfAnguishTrigger(t, scn, opponent, castStart)

		assert.Equal(t, match.GRAVEYARD, iceVapor.Zone)
		assert.Equal(t, match.GRAVEYARD, discarded.Zone)
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

// prepareIceVaporShadowOfAnguishTestRemoval puts Terror Pit in the player's hand
// with enough copies of itself in the mana zone to pay for it. One copy is left
// untapped so the mana burn always has more than one card to choose from and is
// therefore prompted rather than forced.
func prepareIceVaporShadowOfAnguishTestRemoval(t *testing.T, player *match.PlayerReference) *match.Card {
	t.Helper()

	spell, err := player.Player.SpawnCard(iceVaporShadowOfAnguishRemovalUID, match.HAND)
	require.NoError(t, err)

	for range spell.ManaCost + 1 {
		_, err := player.Player.SpawnCard(iceVaporShadowOfAnguishRemovalUID, match.MANAZONE)
		require.NoError(t, err)
	}

	return spell
}

// putIceVaporShadowOfAnguishTestRemovalInShields makes Terror Pit one of the
// player's shields so breaking it offers the shield trigger.
func putIceVaporShadowOfAnguishTestRemovalInShields(t *testing.T, player *match.Player) *match.Card {
	t.Helper()

	card, err := player.SpawnCard(iceVaporShadowOfAnguishRemovalUID, match.HAND)
	require.NoError(t, err)

	moved, err := player.MoveCard(card.ID, match.HAND, match.SHIELDZONE, iceVaporShadowOfAnguishSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.SHIELDZONE, moved.Zone)

	return moved
}

// answerIceVaporShadowOfAnguishTrigger answers one resolution of the ability:
// the opponent's discard followed by their mana burn. It returns the two cards
// they gave up and the message index to hand to the next call, so a second copy
// of Ice Vapor can be answered without racing the prompts.
func answerIceVaporShadowOfAnguishTrigger(t *testing.T, scn *scenario.TestScenario, opponent *match.PlayerReference, since int) (*match.Card, *match.Card, int) {
	t.Helper()

	discardAction, err := scn.WaitForAction(opponent, since)
	require.NoError(t, err)
	require.NotEmpty(t, discardAction.Cards)

	discarded, err := opponent.Player.GetCard(discardAction.Cards[0].CardID, match.HAND)
	require.NoError(t, err)

	manaStart, err := scn.MessageCount(opponent)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(opponent, discarded.ID))

	manaAction, err := scn.WaitForAction(opponent, manaStart)
	require.NoError(t, err)
	require.NotEmpty(t, manaAction.Cards)

	burned, err := opponent.Player.GetCard(manaAction.Cards[0].CardID, match.MANAZONE)
	require.NoError(t, err)

	completionStart, err := scn.MessageCount(opponent)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(opponent, burned.ID))
	require.NoError(t, scn.WaitForMessage(opponent, completionStart, "action", "wait", "state_update"))

	return discarded, burned, completionStart
}

// iceVaporShadowOfAnguishTestOffers reports whether a prompt listed a card.
func iceVaporShadowOfAnguishTestOffers(action *server.ActionMessage, cardID string) bool {
	for _, card := range action.Cards {
		if card.CardID == cardID {
			return true
		}
	}

	return false
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
