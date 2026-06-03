package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const bronzeArmTribeUID = "015fd6bb-37a9-45cf-bb6b-a5497412b880"
const bronzeArmTribeDeckSeedUID = "c5a869f4-a959-4667-a352-92df5369e0b9"

func TestBronzeArmTribe(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.DestroyDeck()
	player.Player.SpawnCard(bronzeArmTribeDeckSeedUID, match.DECK)
	player.Player.SpawnCard(bronzeArmTribeUID, match.DECK)

	player.Player.SpawnCard(bronzeArmTribeUID, match.HAND)
	player.Player.SpawnCard(bronzeArmTribeUID, match.MANAZONE)
	player.Player.SpawnCard(bronzeArmTribeUID, match.MANAZONE)
	player.Player.SpawnCard(bronzeArmTribeUID, match.MANAZONE)

	topDeck := player.Player.PeekDeck(1)
	require.Len(t, topDeck, 1)
	require.Equal(t, bronzeArmTribeDeckSeedUID, topDeck[0].ImageID)

	manaBefore, err := player.Player.Container(match.MANAZONE)
	require.NoError(t, err)
	deckBefore, err := player.Player.Container(match.DECK)
	require.NoError(t, err)

	bronzeArmTribe, err := scn.FindCard(player.Player, match.HAND, bronzeArmTribeUID)
	require.NoError(t, err)

	require.NoError(t, scn.ActionPlayCard(player, bronzeArmTribe.ID))

	_, err = scn.FindCard(player.Player, match.BATTLEZONE, bronzeArmTribeUID)
	require.NoError(t, err)

	manaAfter, err := player.Player.Container(match.MANAZONE)
	require.NoError(t, err)
	deckAfter, err := player.Player.Container(match.DECK)
	require.NoError(t, err)

	movedCard, err := scn.FindCard(player.Player, match.MANAZONE, bronzeArmTribeDeckSeedUID)
	require.NoError(t, err)

	assert.Len(t, manaAfter, len(manaBefore)+1)
	assert.Len(t, deckAfter, len(deckBefore)-1)
	assert.Equal(t, bronzeArmTribeDeckSeedUID, movedCard.ImageID)
	assert.False(t, movedCard.Tapped)
}
