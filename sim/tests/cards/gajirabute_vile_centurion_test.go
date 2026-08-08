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

const gajirabuteVileCenturionUID = "e2552e90-61bf-46ec-80a6-28666bfccd1c"

func TestGajirabuteVileCenturion(t *testing.T) {
	t.Run("puts a selected hidden shield into the graveyard", func(t *testing.T) {
		scn, player, opponent, gajirabute := setupGajirabute(t)

		assert.Equal(t, "Gajirabute, Vile Centurion", gajirabute.Name)
		assert.Equal(t, 3000, gajirabute.Power)
		assert.Equal(t, 6, gajirabute.ManaCost)
		assert.Equal(t, civ.Darkness, gajirabute.Civ)
		assert.True(t, gajirabute.HasFamily(family.DemonCommand))

		messageStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, gajirabute.ID))

		action, err := scn.LatestAction(player, messageStart)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		for _, shield := range action.Cards {
			assert.Equal(t, "backside", shield.ImageID)
			assert.Empty(t, shield.Name)
		}

		selectedShield, err := opponent.Player.GetCard(action.Cards[0].CardID, match.SHIELDZONE)
		require.NoError(t, err)
		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, selectedShield.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, selectedShield.Zone)
		assert.Equal(t, match.BATTLEZONE, gajirabute.Zone)
	})

	t.Run("resolves without prompting when the opponent has no shields", func(t *testing.T) {
		scn, player, opponent, gajirabute := setupGajirabute(t)
		clearGajirabuteTestShields(t, opponent.Player)

		require.NoError(t, scn.ActionPlayCard(player, gajirabute.ID))
		assert.Equal(t, match.BATTLEZONE, gajirabute.Zone)
	})
}

func setupGajirabute(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	player.Player.SpawnCard(gajirabuteVileCenturionUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(gajirabuteVileCenturionUID, match.MANAZONE)
	}
	gajirabute, err := scn.FindCard(player.Player, match.HAND, gajirabuteVileCenturionUID)
	require.NoError(t, err)
	return scn, player, opponent, gajirabute
}

func clearGajirabuteTestShields(t *testing.T, player *match.Player) {
	t.Helper()

	shields, err := player.Container(match.SHIELDZONE)
	require.NoError(t, err)
	for _, shield := range append([]*match.Card(nil), shields...) {
		moved, err := player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, "gajirabute_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
