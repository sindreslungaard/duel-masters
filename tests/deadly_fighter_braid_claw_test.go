package tests

/* func TestDeadlyFighterBraidClaw(t *testing.T) {
	scenario := scenario.New(scenario.Options{})

	t.Run("Player can end turn without Deadly Fighter Braid Claw in the battle zone", func(t *testing.T) {
		scenario.Match.Player1.Player.SpawnCard("c5a869f4-a959-4667-a352-92df5369e0b9", match.BATTLEZONE)

		ctx := match.NewContext(scenario.Match, &match.EndTurnEvent{})
		scenario.Match.HandleFx(ctx)

		assert.Equal(t, ctx.Cancelled(), false)
	})

	t.Run("Player cannot end turn without attacking with Deadly Fighter Braid Claw", func(t *testing.T) {
		_, err := scenario.Match.Player1.Player.SpawnCard("c5a869f4-a959-4667-a352-92df5369e0b9", match.BATTLEZONE)
		assert.Nil(t, err)

		ctx := match.NewContext(scenario.Match, &match.EndTurnEvent{})
		scenario.Match.HandleFx(ctx)

		assert.Equal(t, ctx.Cancelled(), true)
	})

} */
