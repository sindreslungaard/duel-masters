package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Milporo's silent skill is "Draw a card", which makes it a convenient stand-in
// for the mechanic itself: whether the ability resolved is a hand count.
const milporoUID = "460fc2eb-c7cd-42d5-9bed-a98de4f59026"

func TestSilentSkill(t *testing.T) {
	t.Run("a tapped creature may keep itself tapped to use its ability", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)
		creature.Tapped = true

		handBefore := len(handOf(t, player))

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, true))

		assert.True(t, creature.Tapped, "using the ability keeps the creature tapped")
		assert.Len(t, handOf(t, player), handBefore+2, "the drawn card plus the draw step")
	})

	t.Run("the ability is optional", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)
		creature.Tapped = true

		handBefore := len(handOf(t, player))

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, false))

		assert.False(t, creature.Tapped, "declining untaps the creature as normal")
		assert.Len(t, handOf(t, player), handBefore+1, "only the draw step")
	})

	t.Run("an untapped creature is not offered its ability", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		startOwnTurn(t, scn, player, opponent)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "silent skill only triggers on a tapped creature")
		assert.False(t, creature.Tapped)
	})

	t.Run("it does not trigger at the start of the opponent's turn", func(t *testing.T) {
		scn, player, _ := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)
		creature.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "the ability is only offered on its controller's turn")
		assert.True(t, creature.Tapped, "the opponent's untap step leaves it tapped")
	})

	t.Run("a creature that left the battle zone does not resolve its ability", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)
		creature.Tapped = true

		require.NoError(t, scn.ActionEndTurn(player))

		// Removed during the opponent's turn, after the ability would have been
		// eligible on the previous turn but before it could resolve on the next.
		_, err := player.Player.MoveCard(creature.ID, match.BATTLEZONE, match.GRAVEYARD, creature.ID)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, creature.Zone)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(opponent))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})

	t.Run("each copy is offered separately", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		first := silentSkillCreature(t, scn, player)
		second := silentSkillCreature(t, scn, player)
		first.Tapped = true
		second.Tapped = true

		handBefore := len(handOf(t, player))

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, true))
		require.NoError(t, answerSilentSkill(t, scn, player, true))

		assert.True(t, first.Tapped)
		assert.True(t, second.Tapped)
		assert.Len(t, handOf(t, player), handBefore+3, "both abilities plus the draw step")
	})

	t.Run("one copy may be used while the other is declined", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		first := silentSkillCreature(t, scn, player)
		second := silentSkillCreature(t, scn, player)
		first.Tapped = true
		second.Tapped = true

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, true))
		require.NoError(t, answerSilentSkill(t, scn, player, false))

		// The prompts arrive in battle zone order, so the accepted one is the
		// creature that is still tapped afterwards.
		assert.True(t, first.Tapped)
		assert.False(t, second.Tapped)
	})

	t.Run("the ability is offered again on the following turn", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)
		creature.Tapped = true

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, true))
		require.True(t, creature.Tapped)

		startOwnTurn(t, scn, player, opponent)
		require.NoError(t, answerSilentSkill(t, scn, player, true))

		assert.True(t, creature.Tapped, "staying tapped keeps it eligible every turn")
	})

	t.Run("the creature carries the silent skill condition in every zone", func(t *testing.T) {
		scn, player, opponent := silentSkillScenario(t)
		creature := silentSkillCreature(t, scn, player)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.True(t, creature.HasCondition(cnd.SilentSkill))

		// Cards other effects ask about are not always in play, so the keyword
		// has to be rebuilt outside the battle zone too.
		_, err := player.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, creature.ID)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, match.HAND, creature.Zone)
		assert.True(t, creature.HasCondition(cnd.SilentSkill))
	})
}

func silentSkillScenario(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

func silentSkillCreature(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(milporoUID, match.HAND)
	require.NoError(t, err)

	moved, err := player.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "silent_skill_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	return moved
}

// startOwnTurn hands the turn to the opponent and back, which is what puts the
// player at the start of a turn of their own.
func startOwnTurn(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference) {
	t.Helper()

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))
}

// answerSilentSkill answers one pending silent skill prompt and waits for the
// engine to finish with it, so the next assertion sees settled state.
func answerSilentSkill(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, use bool) error {
	t.Helper()

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err)
	require.NotEmpty(t, action.Text)

	if use {
		if err := scn.SubmitAction(player); err != nil {
			return err
		}
	} else if err := scn.CancelAction(player); err != nil {
		return err
	}

	if err := scn.WaitForMessage(player, start, "action", "state_update"); err != nil {
		return err
	}

	// A second silent skill creature opens its own prompt as soon as this one
	// is answered, and the event loop stays blocked until that one is answered
	// too, so there is nothing to settle yet.
	headers, err := scn.MessageHeaders(player, start)
	if err != nil {
		return err
	}
	if slices.Contains(headers, "action") {
		return nil
	}

	return scn.WaitForEventLoop()
}

func handOf(t *testing.T, player *match.PlayerReference) []*match.Card {
	t.Helper()

	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)

	return hand
}
