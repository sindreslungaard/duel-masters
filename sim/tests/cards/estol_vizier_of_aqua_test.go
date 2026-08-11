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
	estolVizierOfAquaUID      = "7b30a97f-ff81-40bc-b7c5-9532a2e1ae85"
	estolVizierOfAquaTopUID   = "7956b4f5-b910-403d-b388-b67c837b7e99" // Scissor Eye
	estolVizierOfAquaSetupSrc = "estol_vizier_of_aqua_test_setup"
)

func TestEstolVizierOfAqua(t *testing.T) {
	t.Run("shields the top card, then looks at one opposing shield", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		estol := spawnMulticolorCardInHand(t, scn, player, estolVizierOfAquaUID)
		spawnCivMana(t, player, civ.Light, 3)
		spawnCivMana(t, player, civ.Water, 3)
		shieldCard := putCardOnTopOfDeck(t, scn, player, estolVizierOfAquaTopUID, estolVizierOfAquaSetupSrc)

		assert.Equal(t, "Estol, Vizier of Aqua", estol.Name)
		assert.Equal(t, 2000, estol.Power)
		assert.Equal(t, 5, estol.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Water}, estol.Civs)
		assert.True(t, estol.HasFamily(family.Initiate))
		assert.True(t, estol.HasFamily(family.LiquidPeople))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, estol.ID))
		assert.Equal(t, match.BATTLEZONE, estol.Zone)
		assert.Equal(t, match.SHIELDZONE, shieldCard.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections, "looking is mandatory and at exactly one shield")
		assert.Equal(t, 1, action.MaxSelections)
		assert.False(t, action.Cancellable)
		for _, shield := range action.Cards {
			assert.Equal(t, "backside", shield.ImageID, "the choice is made blind")
		}

		chosen, err := opponent.Player.GetCard(action.Cards[0].CardID, match.SHIELDZONE)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "show_cards"))

		assert.Equal(t, match.SHIELDZONE, chosen.Zone, "looking never moves the shield")
		opponentShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, opponentShields, 5)
	})

	t.Run("opens no look prompt when the opponent has no shields", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		estol := spawnMulticolorCardInHand(t, scn, player, estolVizierOfAquaUID)
		spawnCivMana(t, player, civ.Light, 3)
		spawnCivMana(t, player, civ.Water, 3)
		clearZone(t, opponent.Player, match.SHIELDZONE, estolVizierOfAquaSetupSrc)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shieldsBefore)

		require.NoError(t, scn.ActionPlayCard(player, estol.ID))
		assert.Equal(t, match.BATTLEZONE, estol.Zone)

		shieldsAfter, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shieldsAfter, shieldCount+1, "its own shield is still added")
	})
}

// clearZone empties a zone by moving every card in it to the graveyard.
func clearZone(t *testing.T, player *match.Player, zone string, source string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, source)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
