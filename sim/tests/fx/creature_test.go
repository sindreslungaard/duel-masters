package fx

import (
	gamefx "duel-masters/game/fx"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const creatureUntapTestUID = "2e10b4fb-3f85-4144-8762-51c04fe609d5"

func TestCreatureUntapDoesNotBypassManaUntapPrevention(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player
	player.SpawnCard(creatureUntapTestUID, match.HAND)
	creature, err := scn.FindCard(player, match.HAND, creatureUntapTestUID)
	require.NoError(t, err)
	untapContext := match.NewContext(scn.Match, &match.UntapStep{})

	creature.Tapped = true
	gamefx.Creature(creature, untapContext)
	assert.False(t, creature.Tapped, "non-mana cards retain the existing untap behavior")

	moved, err := player.MoveCard(creature.ID, match.HAND, match.BATTLEZONE, "creature_untap_test_setup")
	require.NoError(t, err)
	assert.False(t, moved.Tapped, "a successful move into the battle zone resets tap state")

	moved.Tapped = true
	gamefx.Creature(moved, untapContext)
	assert.False(t, moved.Tapped)

	_, err = player.MoveCard(moved.ID, match.BATTLEZONE, match.MANAZONE, "creature_untap_test_setup")
	require.NoError(t, err)
	moved.Tapped = true
	gamefx.Creature(moved, untapContext)
	assert.True(t, moved.Tapped, "fx.Creature must not bypass effects that prevent mana from untapping")
}
