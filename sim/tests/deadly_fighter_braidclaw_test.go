package tests

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeadlyFighterBraidClaw(t *testing.T) {
	scn := scenario.New()

	t.Run("Player can end turn without Deadly Fighter Braid Claw in the battle zone", func(t *testing.T) {
		scn.Match.Player1.Player.SpawnCard("c5a869f4-a959-4667-a352-92df5369e0b9", match.BATTLEZONE)

		ctx := match.NewContext(scn.Match, &match.EndTurnEvent{})
		scn.Match.HandleFx(ctx)

		assert.Equal(t, ctx.Cancelled(), false)
	})

	t.Run("Player cannot end turn without attacking with Deadly Fighter Braid Claw", func(t *testing.T) {
		scn.Match.Player1.Player.SpawnCard("c5a869f4-a959-4667-a352-92df5369e0b9", match.BATTLEZONE)

		ctx := match.NewContext(scn.Match, &match.EndTurnEvent{})
		scn.Match.HandleFx(ctx)

		assert.Equal(t, ctx.Cancelled(), true)
	})

}
