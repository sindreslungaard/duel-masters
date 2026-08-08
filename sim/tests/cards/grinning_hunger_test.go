package cards

import (
	"duel-masters/game/match"
	"duel-masters/server"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	grinningHungerUID              = "5a81f94e-aff8-4c15-a0bb-a787eb7362f4"
	grinningHungerGamilUID         = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	grinningHungerPippieKuppieUID  = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
	grinningHungerMightyShouterUID = "0e26fe1a-a9d1-4c78-80e9-7f4cc0e4c1c8"
)

func TestGrinningHunger(t *testing.T) {
	t.Run("opponent chooses between nonempty zones", func(t *testing.T) {
		scn, caster, opponent, spell := setupGrinningHunger(t)
		chosen := putGrinningHungerTestCreatureInBattlezone(t, scn, opponent.Player, grinningHungerGamilUID)
		other := putGrinningHungerTestCreatureInBattlezone(t, scn, opponent.Player, grinningHungerPippieKuppieUID)

		questionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		question, err := scn.WaitForAction(opponent, questionStart)
		require.NoError(t, err)
		assert.Equal(t, []string{"Battle zone", "Shields"}, question.Choices)

		selectionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitChoice(opponent, 0))

		selection, err := scn.WaitForAction(opponent, selectionStart)
		require.NoError(t, err)
		assert.Empty(t, selection.Choices)
		assert.ElementsMatch(t, []string{chosen.ID, other.ID}, grinningHungerActionCardIDs(selection))

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, chosen.Zone)
		assert.Equal(t, match.BATTLEZONE, other.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("empty battle zone forces a hidden shield choice", func(t *testing.T) {
		scn, caster, opponent, spell := setupGrinningHunger(t)

		selectionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		selection, err := scn.WaitForAction(opponent, selectionStart)
		require.NoError(t, err)
		assert.Empty(t, selection.Choices)
		require.NotEmpty(t, selection.Cards)
		for _, shield := range selection.Cards {
			assert.Equal(t, "backside", shield.ImageID)
			assert.Empty(t, shield.Name)
		}

		selectedShield, err := opponent.Player.GetCard(selection.Cards[0].CardID, match.SHIELDZONE)
		require.NoError(t, err)
		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, selectedShield.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, selectedShield.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("empty shield zone forces a creature choice and honors destruction replacement", func(t *testing.T) {
		scn, caster, opponent, spell := setupGrinningHunger(t)
		clearGrinningHungerTestZone(t, opponent.Player, match.SHIELDZONE)
		replacementCreature := putGrinningHungerTestCreatureInBattlezone(t, scn, opponent.Player, grinningHungerMightyShouterUID)
		other := putGrinningHungerTestCreatureInBattlezone(t, scn, opponent.Player, grinningHungerGamilUID)

		selectionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		selection, err := scn.WaitForAction(opponent, selectionStart)
		require.NoError(t, err)
		assert.Empty(t, selection.Choices)
		assert.ElementsMatch(t, []string{replacementCreature.ID, other.ID}, grinningHungerActionCardIDs(selection))

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, replacementCreature.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.MANAZONE, replacementCreature.Zone)
		assert.Equal(t, match.BATTLEZONE, other.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("no legal cards resolves without prompting", func(t *testing.T) {
		scn, caster, opponent, spell := setupGrinningHunger(t)
		clearGrinningHungerTestZone(t, opponent.Player, match.SHIELDZONE)

		messageStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		headers, err := scn.MessageHeaders(opponent, messageStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}

func setupGrinningHunger(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	caster := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(caster.Player))

	caster.Player.SpawnCard(grinningHungerUID, match.HAND)
	for range 4 {
		caster.Player.SpawnCard(grinningHungerUID, match.MANAZONE)
	}

	spell, err := scn.FindCard(caster.Player, match.HAND, grinningHungerUID)
	require.NoError(t, err)
	return scn, caster, opponent, spell
}

func putGrinningHungerTestCreatureInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "grinning_hunger_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

func clearGrinningHungerTestZone(t *testing.T, player *match.Player, zone string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, "grinning_hunger_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}

func grinningHungerActionCardIDs(action *server.ActionMessage) []string {
	ids := make([]string, len(action.Cards))
	for i, card := range action.Cards {
		ids[i] = card.CardID
	}
	return ids
}
