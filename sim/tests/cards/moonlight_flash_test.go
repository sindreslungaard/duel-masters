package cards

import (
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moonlightFlashUID     = "b12f1d66-46ee-49b9-878d-59cc3d515633"
	moonlightFlashManaUID = "0cc5279e-0a26-41a8-a2a5-f7711120b772" // Lah, Purification Enforcer (light)
	moonlightFlashSrc     = "moonlight_flash_test_setup"
)

// Moonlight Flash selects its targets straight from the opponent's battle
// zone with no filter excluding cnd.CantBeTappedByOpp, unlike TapOpCreature.
// It still has to honor the condition: every tap, filtered or not, goes
// through player.TapCard now, so this covers a path the filter itself never
// protected.
func TestMoonlightFlashRespectsCantBeTappedByOpp(t *testing.T) {
	scn, player, opponent := setupDuel(t)

	ulex := putCardInBattlezone(t, scn, opponent.Player, ulexTheDauntlessUID, moonlightFlashSrc)
	ally := putCardInBattlezone(t, scn, opponent.Player, ulexAllyUID, moonlightFlashSrc)
	passTurnToSelf(t, scn, player, opponent)

	player.Player.SpawnCard(moonlightFlashUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(moonlightFlashManaUID, match.MANAZONE)
	}

	spell, err := scn.FindCard(player.Player, match.HAND, moonlightFlashUID)
	require.NoError(t, err)

	promptStart, err := scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, spell.ID))

	action, err := scn.LatestAction(player, promptStart)
	require.NoError(t, err)
	require.NotEmpty(t, action.Cards)

	ids := make([]string, 0, len(action.Cards))
	for _, c := range action.Cards {
		ids = append(ids, c.CardID)
	}

	completionStart, err := scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(player, ids...))
	require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

	assert.False(t, ulex.Tapped, "its owner's opponent can't tap it, even through an unfiltered select")
	assert.True(t, ally.Tapped, "the effect still reaches a creature without the protection")
}
