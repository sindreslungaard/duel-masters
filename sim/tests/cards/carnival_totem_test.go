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

const (
	carnivalTotemUID           = "f6d473c1-952f-482a-85da-cb29cfb46b07"
	carnivalTotemFirstHandUID  = "b3975c0b-2978-4b1a-8225-78d420ff941d"
	carnivalTotemSecondHandUID = "1484ec6d-c1b5-4fc4-abaf-a16c08cfc5f7"
)

func TestCarnivalTotem(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	clearCarnivalTotemTestZone(t, player.Player, match.HAND)

	player.Player.SpawnCard(carnivalTotemUID, match.HAND)
	player.Player.SpawnCard(carnivalTotemFirstHandUID, match.HAND)
	player.Player.SpawnCard(carnivalTotemSecondHandUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(carnivalTotemUID, match.MANAZONE)
	}

	carnival, err := scn.FindCard(player.Player, match.HAND, carnivalTotemUID)
	require.NoError(t, err)
	firstHandCard, err := scn.FindCard(player.Player, match.HAND, carnivalTotemFirstHandUID)
	require.NoError(t, err)
	secondHandCard, err := scn.FindCard(player.Player, match.HAND, carnivalTotemSecondHandUID)
	require.NoError(t, err)
	originalMana, err := player.Player.Container(match.MANAZONE)
	require.NoError(t, err)
	originalMana = append([]*match.Card(nil), originalMana...)

	assert.Equal(t, "Carnival Totem", carnival.Name)
	assert.Equal(t, 7000, carnival.Power)
	assert.Equal(t, 6, carnival.ManaCost)
	assert.Equal(t, []string{civ.Nature}, carnival.Civs)
	assert.True(t, carnival.HasFamily(family.MysteryTotem))

	require.NoError(t, scn.ActionPlayCard(player, carnival.ID))

	assert.Equal(t, match.BATTLEZONE, carnival.Zone)
	for _, manaCard := range originalMana {
		assert.Equal(t, match.HAND, manaCard.Zone)
	}
	for _, handCard := range []*match.Card{firstHandCard, secondHandCard} {
		assert.Equal(t, match.MANAZONE, handCard.Zone)
		assert.True(t, handCard.Tapped)
	}
}

func clearCarnivalTotemTestZone(t *testing.T, player *match.Player, zone string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, "carnival_totem_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
