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
	kingOquanosUID     = "296b783c-9200-466e-9eb3-a1a82933c19d"
	kingOquanosManaUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
)

func TestKingOquanos(t *testing.T) {
	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	owner.Player.SpawnCard(kingOquanosUID, match.HAND)
	king, err := scn.FindCard(owner.Player, match.HAND, kingOquanosUID)
	require.NoError(t, err)
	opponent.Player.SpawnCard(kingOquanosManaUID, match.MANAZONE)
	opponentMana, err := opponent.Player.Container(match.MANAZONE)
	require.NoError(t, err)
	opponentMana[0].Tapped = true
	assert.Equal(t, 2000, scn.Match.GetPower(king, false), "the ability is inactive outside the battle zone")
	clearKingOquanosTestMana(t, opponent.Player)

	moved, err := owner.Player.MoveCard(king.ID, match.HAND, match.BATTLEZONE, "king_oquanos_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	assert.Equal(t, "King Oquanos", king.Name)
	assert.Equal(t, 2000, king.Power)
	assert.Equal(t, 8, king.ManaCost)
	assert.Equal(t, []string{civ.Water}, king.Civs)
	assert.True(t, king.HasFamily(family.Leviathan))
	assert.Equal(t, 2000, scn.Match.GetPower(king, false))
	assert.False(t, king.HasCondition(cnd.DoubleBreaker))

	opponent.Player.SpawnCard(kingOquanosManaUID, match.HAND)
	for range 2 {
		opponent.Player.SpawnCard(kingOquanosManaUID, match.MANAZONE)
	}

	require.NoError(t, scn.ActionEndTurn(owner))
	playable, err := scn.FindCard(opponent.Player, match.HAND, kingOquanosManaUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(opponent, playable.ID))

	assert.Equal(t, 6000, scn.Match.GetPower(king, false))
	assert.Equal(t, 2000, king.Power, "dynamic power must not mutate printed power")
	assert.True(t, king.HasCondition(cnd.DoubleBreaker))

	require.NoError(t, scn.ActionEndTurn(opponent))
	assert.Equal(t, 6000, scn.Match.GetPower(king, false))
	require.NoError(t, scn.ActionEndTurn(owner))
	assert.Equal(t, 2000, scn.Match.GetPower(king, false))
	assert.False(t, king.HasCondition(cnd.DoubleBreaker))
}

func clearKingOquanosTestMana(t *testing.T, player *match.Player) {
	t.Helper()

	mana, err := player.Container(match.MANAZONE)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), mana...) {
		moved, err := player.MoveCard(card.ID, match.MANAZONE, match.GRAVEYARD, "king_oquanos_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
