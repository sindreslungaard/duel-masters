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
	galekTheShadowWarriorUID        = "bde9c44b-ab35-4063-b827-7816cfc35bad"
	galekTheShadowWarriorBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker)
	galekTheShadowWarriorPlainUID   = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (no blocker)
	galekTheShadowWarriorSetupSrc   = "galek_the_shadow_warrior_test_setup"
)

func TestGalekTheShadowWarrior(t *testing.T) {
	t.Run("only blockers are offered when there is a choice", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		firstBlocker := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorBlockerUID, galekTheShadowWarriorSetupSrc)
		secondBlocker := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorBlockerUID, galekTheShadowWarriorSetupSrc)
		plain := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorPlainUID, galekTheShadowWarriorSetupSrc)

		galek, promptStart := prepareGalek(t, scn, player)

		assert.Equal(t, "Galek, the Shadow Warrior", galek.Name)
		assert.Equal(t, 2000, galek.Power)
		assert.Equal(t, 5, galek.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, galek.Civs)
		assert.True(t, galek.HasFamily(family.Ghost))
		assert.True(t, galek.HasFamily(family.Human))

		opponentHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(opponentHand)

		require.NoError(t, scn.ActionPlayCard(player, galek.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.False(t, action.Cancellable, "the destruction is mandatory")
		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		assert.ElementsMatch(t, []string{firstBlocker.ID, secondBlocker.ID}, offered,
			"only creatures that have blocker may be chosen")

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, firstBlocker.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, firstBlocker.Zone)
		assert.Equal(t, match.BATTLEZONE, secondBlocker.Zone)
		assert.Equal(t, match.BATTLEZONE, plain.Zone)

		opponentHandAfter, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, opponentHandAfter, handCount-1)
	})

	t.Run("destroys the only blocker without prompting", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		blocker := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorBlockerUID, galekTheShadowWarriorSetupSrc)
		plain := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorPlainUID, galekTheShadowWarriorSetupSrc)
		galek, promptStart := prepareGalek(t, scn, player)

		require.NoError(t, scn.ActionPlayCard(player, galek.ID))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		actions := 0
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "a forced single target needs no prompt")

		assert.Equal(t, match.GRAVEYARD, blocker.Zone)
		assert.Equal(t, match.BATTLEZONE, plain.Zone)
	})

	t.Run("still discards when the opponent has no blocker", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		plain := putCardInBattlezone(t, scn, opponent.Player, galekTheShadowWarriorPlainUID, galekTheShadowWarriorSetupSrc)
		galek, promptStart := prepareGalek(t, scn, player)

		opponentHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(opponentHand)

		require.NoError(t, scn.ActionPlayCard(player, galek.ID))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		actions := 0
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "only the mana payment is prompted")

		assert.Equal(t, match.BATTLEZONE, plain.Zone)
		opponentHandAfter, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, opponentHandAfter, handCount-1)
	})
}

// prepareGalek spawns Galek with enough mana and cycles a turn so that blockers
// already in play carry their blocker condition, then returns the message count
// to measure prompts from.
func prepareGalek(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) (*match.Card, int) {
	t.Helper()

	galek := spawnMulticolorCardInHand(t, scn, player, galekTheShadowWarriorUID)
	spawnCivMana(t, player, civ.Darkness, 3)
	spawnCivMana(t, player, civ.Fire, 3)

	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	promptStart, err := scn.MessageCount(player)
	require.NoError(t, err)

	return galek, promptStart
}

// putCardInBattlezone is shared by the multicolored card tests.
func putCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string, source string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, source)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
