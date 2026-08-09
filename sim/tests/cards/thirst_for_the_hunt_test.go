package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const thirstForTheHuntUID = "244080a8-c85f-4e05-b403-dfae3fac0618"

func TestThirstForTheHunt(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	firstCreature := putThirstForTheHuntTestCreatureInBattlezone(t, scn, player.Player)
	secondCreature := putThirstForTheHuntTestCreatureInBattlezone(t, scn, player.Player)
	player.Player.SpawnCard(thirstForTheHuntUID, match.HAND)
	player.Player.SpawnCard(thirstForTheHuntUID, match.MANAZONE)
	spell, err := scn.FindCard(player.Player, match.HAND, thirstForTheHuntUID)
	require.NoError(t, err)

	assert.Equal(t, "Thirst for the Hunt", spell.Name)
	assert.Equal(t, 1, spell.ManaCost)
	assert.Equal(t, []string{civ.Nature}, spell.Civs)

	require.NoError(t, scn.ActionPlayCard(player, spell.ID))

	for _, creature := range []*match.Card{firstCreature, secondCreature} {
		assert.Equal(t, 3000, scn.Match.GetPower(creature, true))
		condition, err := creature.GetCondition(cnd.PowerAttacker)
		require.NoError(t, err)
		assert.Equal(t, 1000, condition.Val)
		assert.Equal(t, spell.ID, condition.Src)
	}

	lateCreature := putThirstForTheHuntTestCreatureInBattlezone(t, scn, player.Player)
	assert.Equal(t, 2000, scn.Match.GetPower(lateCreature, true), "creatures arriving later must not receive the resolved effect")

	require.NoError(t, scn.ActionEndTurn(player))
	assert.Equal(t, 2000, scn.Match.GetPower(firstCreature, true))
	assert.Equal(t, 2000, scn.Match.GetPower(secondCreature, true))
}

func putThirstForTheHuntTestCreatureInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player) *match.Card {
	t.Helper()

	player.SpawnCard(scowlingTomatoUID, match.HAND)
	creature, err := scn.FindCard(player, match.HAND, scowlingTomatoUID)
	require.NoError(t, err)
	moved, err := player.MoveCard(creature.ID, match.HAND, match.BATTLEZONE, "thirst_for_the_hunt_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
