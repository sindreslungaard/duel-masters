package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	blizzardOfSpearsUID  = "0219aa19-f201-4e11-92c5-59f4f5aaa697"
	pippieKuppieUID      = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
	gamilUID             = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	bolgashDragonUID     = "82e062b1-bf9d-444e-9cd3-2fdef007ba64"
	immortalBaronVorgUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
	mightyShouterUID     = "0e26fe1a-a9d1-4c78-80e9-7f4cc0e4c1c8"
	zagaanUIDForBlizzard = "07a0115e-797a-49d8-90bf-9ea6de39978d"
)

func TestBlizzardOfSpears(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.Opponent(player.Player)

	player.Player.SpawnCard(blizzardOfSpearsUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(blizzardOfSpearsUID, match.MANAZONE)
	}

	pippie := putBlizzardTestCreatureInBattlezone(t, scn, player.Player, pippieKuppieUID)
	protectedAtResolution := putBlizzardTestCreatureInBattlezone(t, scn, player.Player, bolgashDragonUID)
	ownersLowPowerCreature := putBlizzardTestCreatureInBattlezone(t, scn, player.Player, immortalBaronVorgUID)

	opponentsBoundaryCreature := putBlizzardTestCreatureInBattlezone(t, scn, opponent, gamilUID)
	opponentsReplacementCreature := putBlizzardTestCreatureInBattlezone(t, scn, opponent, mightyShouterUID)
	opponentsHighPowerCreature := putBlizzardTestCreatureInBattlezone(t, scn, opponent, zagaanUIDForBlizzard)

	require.Equal(t, 5000, scn.Match.GetPower(protectedAtResolution, false))

	spell, err := scn.FindCard(player.Player, match.HAND, blizzardOfSpearsUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, spell.ID))

	assert.Equal(t, match.GRAVEYARD, spell.Zone)
	assert.Equal(t, match.GRAVEYARD, pippie.Zone)
	assert.Equal(t, match.GRAVEYARD, ownersLowPowerCreature.Zone)
	assert.Equal(t, match.GRAVEYARD, opponentsBoundaryCreature.Zone)
	assert.Equal(t, match.MANAZONE, opponentsReplacementCreature.Zone)
	assert.Equal(t, match.BATTLEZONE, opponentsHighPowerCreature.Zone)

	// Blizzard determines all affected creatures before destroying any of them.
	// Pippie initially raises this creature above 4000, so it survives even
	// though Pippie's destruction subsequently lowers it back to 4000.
	assert.Equal(t, match.BATTLEZONE, protectedAtResolution.Zone)
	assert.Equal(t, 4000, scn.Match.GetPower(protectedAtResolution, false))

	// Ensure spell resolution left no blocked Player.Action sender or event loop.
	require.NoError(t, scn.ActionEndTurn(player))
}

func putBlizzardTestCreatureInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)

	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	return moved
}
