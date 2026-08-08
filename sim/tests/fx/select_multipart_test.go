package fx

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	selectMultipartCrisisBoulderUID = "5c424a0f-5bbd-41cd-9279-2b408f7e5935"
	selectMultipartGamilUID         = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	selectMultipartPippieKuppieUID  = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
)

func TestSelectMultipart(t *testing.T) {
	t.Run("mandatory single outcome bypasses the prompt", func(t *testing.T) {
		scn, caster, opponent, spell := setupSelectMultipartTest(t)
		target := putSelectMultipartTestCreatureInBattlezone(t, scn, opponent.Player, selectMultipartGamilUID)

		messageStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		headers, err := scn.MessageHeaders(opponent, messageStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.GRAVEYARD, target.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("empty groups are hidden and invalid selections are rejected", func(t *testing.T) {
		scn, caster, opponent, spell := setupSelectMultipartTest(t)
		chosen := putSelectMultipartTestCreatureInBattlezone(t, scn, opponent.Player, selectMultipartGamilUID)
		other := putSelectMultipartTestCreatureInBattlezone(t, scn, opponent.Player, selectMultipartPippieKuppieUID)

		actionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		action, err := scn.WaitForMultipartAction(opponent, actionStart)
		require.NoError(t, err)
		require.Contains(t, action.Cards, "Your creatures")
		assert.NotContains(t, action.Cards, "Your mana")
		assert.Len(t, action.Cards["Your creatures"], 2)

		warningStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, "not-an-offered-card"))
		require.NoError(t, scn.WaitForMessage(opponent, warningStart, "action_error"))
		assert.Equal(t, match.BATTLEZONE, chosen.Zone)
		assert.Equal(t, match.BATTLEZONE, other.Zone)

		warningStart, err = scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(opponent))
		require.NoError(t, scn.WaitForMessage(opponent, warningStart, "action_error"))

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, chosen.Zone)
		assert.Equal(t, match.BATTLEZONE, other.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}

func setupSelectMultipartTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	caster := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(caster.Player))

	caster.Player.SpawnCard(selectMultipartCrisisBoulderUID, match.HAND)
	for range 4 {
		caster.Player.SpawnCard(selectMultipartCrisisBoulderUID, match.MANAZONE)
	}

	spell, err := scn.FindCard(caster.Player, match.HAND, selectMultipartCrisisBoulderUID)
	require.NoError(t, err)
	return scn, caster, opponent, spell
}

func putSelectMultipartTestCreatureInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "select_multipart_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
