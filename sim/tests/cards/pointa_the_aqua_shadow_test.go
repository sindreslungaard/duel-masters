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
	pointaTheAquaShadowUID      = "bc128af9-0fc2-4a1b-b10e-2f695f05c24e"
	pointaTheAquaShadowSetupSrc = "pointa_the_aqua_shadow_test_setup"
)

func TestPointaTheAquaShadow(t *testing.T) {
	t.Run("looks at one opposing shield, then discards a random card", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		pointa := spawnMulticolorCardInHand(t, scn, player, pointaTheAquaShadowUID)
		spawnCivMana(t, player, civ.Water, 3)
		spawnCivMana(t, player, civ.Darkness, 3)

		opponentHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(opponentHand)
		require.NotEmpty(t, opponentHand)

		assert.Equal(t, "Pointa, the Aqua Shadow", pointa.Name)
		assert.Equal(t, 2000, pointa.Power)
		assert.Equal(t, 5, pointa.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, pointa.Civs)
		assert.True(t, pointa.HasFamily(family.LiquidPeople))
		assert.True(t, pointa.HasFamily(family.Ghost))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, pointa.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		chosen, err := opponent.Player.GetCard(action.Cards[0].CardID, match.SHIELDZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, pointa.Zone)
		assert.Equal(t, match.SHIELDZONE, chosen.Zone, "looking never moves the shield")

		opponentHandAfter, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, opponentHandAfter, handCount-1, "the discard is not optional")
		opponentGraveyard, err := opponent.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		assert.Len(t, opponentGraveyard, 1)
	})

	t.Run("resolves when the opponent has neither shields nor cards in hand", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		pointa := spawnMulticolorCardInHand(t, scn, player, pointaTheAquaShadowUID)
		spawnCivMana(t, player, civ.Water, 3)
		spawnCivMana(t, player, civ.Darkness, 3)
		clearZone(t, opponent.Player, match.SHIELDZONE, pointaTheAquaShadowSetupSrc)
		clearZone(t, opponent.Player, match.HAND, pointaTheAquaShadowSetupSrc)

		require.NoError(t, scn.ActionPlayCard(player, pointa.ID))
		assert.Equal(t, match.BATTLEZONE, pointa.Zone)
	})
}
