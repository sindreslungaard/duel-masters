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

const shamanBroccoliUID = "b22f0d6b-7703-4bd4-b97f-4389f907577e"

func TestShamanBroccoli(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player

	broccoli := putShamanBroccoliTestCardInBattlezone(t, scn, player, shamanBroccoliUID)
	source := putShamanBroccoliTestCardInBattlezone(t, scn, player, scowlingTomatoUID)

	assert.Equal(t, "Shaman Broccoli", broccoli.Name)
	assert.Equal(t, 1000, broccoli.Power)
	assert.Equal(t, 2, broccoli.ManaCost)
	assert.Equal(t, civ.Nature, broccoli.Civ)
	assert.True(t, broccoli.HasFamily(family.WildVeggies))

	scn.Match.Destroy(broccoli, source, match.DestroyedByMiscAbility)

	assert.Equal(t, match.MANAZONE, broccoli.Zone)
	assert.False(t, broccoli.Tapped)
	_, err := scn.FindCard(player, match.GRAVEYARD, shamanBroccoliUID)
	assert.Error(t, err)
}

func putShamanBroccoliTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "shaman_broccoli_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
