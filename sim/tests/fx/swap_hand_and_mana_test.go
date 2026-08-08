package fx

import (
	gamefx "duel-masters/game/fx"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	swapHandAndManaSourceUID = "f6d473c1-952f-482a-85da-cb29cfb46b07"
	swapHandAndManaHandUID   = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	swapHandAndManaManaUID   = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
)

func TestSwapHandAndMana(t *testing.T) {
	t.Run("swaps the original zone snapshots", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer().Player
		clearSwapHandAndManaTestZone(t, player, match.HAND)

		player.SpawnCard(swapHandAndManaSourceUID, match.HAND)
		source, err := scn.FindCard(player, match.HAND, swapHandAndManaSourceUID)
		require.NoError(t, err)
		movedSource, err := player.MoveCard(source.ID, match.HAND, match.BATTLEZONE, "swap_hand_and_mana_test_setup")
		require.NoError(t, err)

		player.SpawnCard(swapHandAndManaHandUID, match.HAND)
		handCard, err := scn.FindCard(player, match.HAND, swapHandAndManaHandUID)
		require.NoError(t, err)
		player.SpawnCard(swapHandAndManaManaUID, match.MANAZONE)
		manaCard, err := scn.FindCard(player, match.MANAZONE, swapHandAndManaManaUID)
		require.NoError(t, err)

		gamefx.SwapHandAndMana(movedSource, player)

		assert.Equal(t, match.BATTLEZONE, movedSource.Zone)
		assert.Equal(t, match.HAND, manaCard.Zone)
		assert.Equal(t, match.MANAZONE, handCard.Zone)
		assert.True(t, handCard.Tapped)
	})

	t.Run("does not move a resolving source that is still in hand", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer().Player
		clearSwapHandAndManaTestZone(t, player, match.HAND)

		player.SpawnCard(swapHandAndManaSourceUID, match.HAND)
		source, err := scn.FindCard(player, match.HAND, swapHandAndManaSourceUID)
		require.NoError(t, err)
		player.SpawnCard(swapHandAndManaHandUID, match.HAND)
		handCard, err := scn.FindCard(player, match.HAND, swapHandAndManaHandUID)
		require.NoError(t, err)
		player.SpawnCard(swapHandAndManaManaUID, match.MANAZONE)
		manaCard, err := scn.FindCard(player, match.MANAZONE, swapHandAndManaManaUID)
		require.NoError(t, err)

		gamefx.SwapHandAndMana(source, player)

		assert.Equal(t, match.HAND, source.Zone)
		assert.Equal(t, match.HAND, manaCard.Zone)
		assert.Equal(t, match.MANAZONE, handCard.Zone)
		assert.True(t, handCard.Tapped)
	})
}

func clearSwapHandAndManaTestZone(t *testing.T, player *match.Player, zone string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, "swap_hand_and_mana_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
