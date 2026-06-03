package tests

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deadlyFighterBraidClawUID = "c5a869f4-a959-4667-a352-92df5369e0b9"

func TestDeadlyFighterBraidClaw(t *testing.T) {
	t.Run("Player can end turn without Deadly Fighter Braid Claw in the battle zone", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		turn := scn.Match.Turn

		require.NoError(t, scn.ActionEndTurn(player))

		assert.NotEqual(t, turn, scn.Match.Turn)
	})

	t.Run("Player cannot end turn without attacking with Deadly Fighter Braid Claw", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		player.Player.SpawnCard(deadlyFighterBraidClawUID, match.HAND)
		player.Player.SpawnCard(deadlyFighterBraidClawUID, match.MANAZONE)

		braidClaw, err := scn.FindCard(player.Player, match.HAND, deadlyFighterBraidClawUID)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, braidClaw.ID))

		_, err = scn.FindCard(player.Player, match.BATTLEZONE, deadlyFighterBraidClawUID)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		turn := scn.Match.Turn
		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, turn, scn.Match.Turn)
	})
}
