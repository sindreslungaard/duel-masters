package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ryudmilaChannelerOfSunsUID = "51e58f0d-6c59-4d15-99a1-ce2908756b28"
	ryudmilaUntappedAllyUID    = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	ryudmilaSecondAllyUID      = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
	ryudmilaNonCreatureUID     = "0219aa19-f201-4e11-92c5-59f4f5aaa697"
)

func TestRyudmilaChannelerOfSuns(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player

	ryudmila := putRyudmilaTestCreatureInBattlezone(t, scn, player, ryudmilaChannelerOfSunsUID)
	firstAlly := putRyudmilaTestCreatureInBattlezone(t, scn, player, ryudmilaUntappedAllyUID)
	secondAlly := putRyudmilaTestCreatureInBattlezone(t, scn, player, ryudmilaSecondAllyUID)
	player.SpawnCard(ryudmilaNonCreatureUID, match.HAND)
	nonCreature, err := scn.FindCard(player, match.HAND, ryudmilaNonCreatureUID)
	require.NoError(t, err)
	_, err = player.MoveCard(nonCreature.ID, match.HAND, match.BATTLEZONE, "ryudmila_test_setup")
	require.NoError(t, err)

	assert.Equal(t, "Ryudmila, Channeler of Suns", ryudmila.Name)
	assert.Equal(t, 2000, ryudmila.Power)
	assert.Equal(t, 5, ryudmila.ManaCost)
	assert.Equal(t, civ.Light, ryudmila.Civ)
	assert.True(t, ryudmila.HasFamily(family.MechaDelSol))
	assert.Equal(t, 6000, scn.Match.GetPower(ryudmila, false), "non-creature cards must not contribute")

	secondAlly.Tapped = true
	assert.Equal(t, 4000, scn.Match.GetPower(ryudmila, false))

	scn.Match.Destroy(ryudmila, firstAlly, match.DestroyedByMiscAbility)
	assert.Equal(t, match.DECK, ryudmila.Zone)
	_, err = scn.FindCard(player, match.GRAVEYARD, ryudmilaChannelerOfSunsUID)
	assert.Error(t, err)
	_, err = scn.FindCard(player, match.DECK, ryudmilaChannelerOfSunsUID)
	require.NoError(t, err)
}

func putRyudmilaTestCreatureInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "ryudmila_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	moved.AddCondition(cnd.Creature, nil, nil)
	return moved
}
