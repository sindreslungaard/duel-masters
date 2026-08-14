package fx

import (
	"duel-masters/game/cnd"
	gamefx "duel-masters/game/fx"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const forcePutCreatureIntoBZTestUID = "2e10b4fb-3f85-4144-8762-51c04fe609d5"

func TestForcePutCreatureIntoBZAppliesSuccessStateAfterMove(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player
	player.SpawnCard(forcePutCreatureIntoBZTestUID, match.MANAZONE)
	creature, err := scn.FindCard(player, match.MANAZONE, forcePutCreatureIntoBZTestUID)
	require.NoError(t, err)

	gamefx.ForcePutCreatureIntoBZ(
		match.NewContext(scn.Match, struct{}{}),
		creature,
		match.MANAZONE,
		creature,
	)

	assert.Equal(t, match.BATTLEZONE, creature.Zone)
	assert.True(t, creature.HasCondition(cnd.SummoningSickness))
}

func TestForcePutCreatureIntoBZDoesNotApplySuccessStateWhenMoveIsPrevented(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player
	player.SpawnCard(forcePutCreatureIntoBZTestUID, match.MANAZONE)
	creature, err := scn.FindCard(player, match.MANAZONE, forcePutCreatureIntoBZTestUID)
	require.NoError(t, err)
	creature.Use(func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.MoveCard); ok &&
			event.CardID == card.ID &&
			event.To == match.BATTLEZONE {
			ctx.InterruptFlow()
		}
	})

	gamefx.ForcePutCreatureIntoBZ(
		match.NewContext(scn.Match, struct{}{}),
		creature,
		match.MANAZONE,
		creature,
	)

	assert.Equal(t, match.MANAZONE, creature.Zone)
	assert.False(t, creature.HasCondition(cnd.SummoningSickness))
}
