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
	skyswordTheSavageVizierUID      = "f30403b3-75b0-4916-9d5e-5d21e0461326"
	skyswordTheSavageVizierManaUID  = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	skyswordTheSavageVizierTopUID   = "7956b4f5-b910-403d-b388-b67c837b7e99" // Scissor Eye
	skyswordTheSavageVizierNextUID  = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd
	skyswordTheSavageVizierSetupSrc = "skysword_the_savage_vizier_test_setup"
)

func TestSkyswordTheSavageVizier(t *testing.T) {
	t.Run("puts the top card into mana and the next one into the shields", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		skysword := spawnMulticolorCardInHand(t, scn, player, skyswordTheSavageVizierUID)
		spawnCivMana(t, player, civ.Light, 3)
		spawnCivMana(t, player, civ.Nature, 3)

		// Pushed onto the deck bottom first, so the mana card ends up on top.
		shieldCard := putCardOnTopOfDeck(t, scn, player, skyswordTheSavageVizierNextUID, skyswordTheSavageVizierSetupSrc)
		manaCard := putCardOnTopOfDeck(t, scn, player, skyswordTheSavageVizierTopUID, skyswordTheSavageVizierSetupSrc)

		assert.Equal(t, "Skysword, the Savage Vizier", skysword.Name)
		assert.Equal(t, 2000, skysword.Power)
		assert.Equal(t, 5, skysword.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, skysword.Civs)
		assert.True(t, skysword.HasFamily(family.BeastFolk))
		assert.True(t, skysword.HasFamily(family.Initiate))

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shieldsBefore)

		require.NoError(t, scn.ActionPlayCard(player, skysword.ID))
		assert.Equal(t, match.BATTLEZONE, skysword.Zone)

		assert.Equal(t, match.MANAZONE, manaCard.Zone, "the first top card goes to the mana zone")
		assert.False(t, manaCard.Tapped, "a mono coloured card arrives untapped")
		assert.Equal(t, match.SHIELDZONE, shieldCard.Zone, "the next top card becomes a shield")

		shieldsAfter, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shieldsAfter, shieldCount+1)
	})

	t.Run("needs one mana of each of its civilizations", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		skysword := spawnMulticolorCardInHand(t, scn, player, skyswordTheSavageVizierUID)

		lightOnly := spawnCivMana(t, player, civ.Light, 5)
		assert.False(t, player.Player.CanPlayCard(skysword, lightOnly))

		nature := spawnCivMana(t, player, civ.Nature, 1)
		assert.True(t, player.Player.CanPlayCard(skysword, append(lightOnly, nature...)))
	})

	t.Run("is put into the mana zone tapped", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		skysword := spawnMulticolorCardInHand(t, scn, player, skyswordTheSavageVizierUID)
		moved, err := player.Player.MoveCard(skysword.ID, match.HAND, match.MANAZONE, skyswordTheSavageVizierSetupSrc)
		require.NoError(t, err)
		assert.True(t, moved.Tapped)
	})
}

// spawnMulticolorCardInHand is shared by the multicolored card tests.
func spawnMulticolorCardInHand(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	player.Player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player.Player, match.HAND, uid)
	require.NoError(t, err)
	return card
}

// spawnCivMana puts count cards of the given civilization into a player's mana
// zone and returns them.
func spawnCivMana(t *testing.T, player *match.PlayerReference, civilization string, count int) []*match.Card {
	t.Helper()

	uid, ok := multicolorTestManaUIDs[civilization]
	require.True(t, ok, "no test mana card registered for %s", civilization)

	spawned := make([]*match.Card, 0, count)
	for range count {
		card, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
		spawned = append(spawned, card)
	}

	return spawned
}

var multicolorTestManaUIDs = map[string]string{
	civ.Light:    "7b58e8c2-0b1e-4ef5-812f-e667c2092c73", // Reusol, the Oracle
	civ.Water:    "9781089f-1aa9-4a75-b106-35e9d431e31d", // Aqua Vehicle
	civ.Darkness: "e2b992ee-91a3-49d3-8228-7be60a0b9ec5", // Writhing Bone Ghoul
	civ.Fire:     "af3bc221-1cc2-4f58-83ea-2673ac2c66c5", // Immortal Baron, Vorg
	civ.Nature:   "1d72eb3e-5185-449a-a16f-391bd2338343", // Burning Mane
}

// putCardOnTopOfDeck spawns a card and moves it to the front of the deck, which
// is the position PeekDeck reads from.
func putCardOnTopOfDeck(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string, source string) *match.Card {
	t.Helper()

	player.Player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player.Player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.Player.MoveCardToFront(card.ID, match.HAND, match.DECK, source)
	require.NoError(t, err)
	require.Equal(t, match.DECK, moved.Zone)
	return moved
}
