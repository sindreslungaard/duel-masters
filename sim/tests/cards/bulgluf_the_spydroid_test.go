package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bulglufTheSpydroidUID      = "511d5780-d5ee-41df-a9ed-c266c55f7a84"
	bulglufTheSpydroidSeedUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	bulglufTheSpydroidSetupSrc = "bulgluf_the_spydroid_test_setup"
)

func TestBulglufTheSpydroid(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		bulgluf := putCardInBattlezone(t, scn, player.Player, bulglufTheSpydroidUID, bulglufTheSpydroidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Bulgluf, the Spydroid", bulgluf.Name)
		assert.Equal(t, 4000, bulgluf.Power)
		assert.Equal(t, 6, bulgluf.ManaCost)
		assert.Equal(t, []string{civ.Light}, bulgluf.Civs)
		assert.Equal(t, []string{civ.Light}, bulgluf.ManaRequirement)
		assert.True(t, bulgluf.HasFamily(family.Soltrooper))
		assert.True(t, bulgluf.HasCondition(cnd.SilentSkill))
	})

	t.Run("adds the top card of the deck to the shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		bulgluf := putCardInBattlezone(t, scn, player.Player, bulglufTheSpydroidUID, bulglufTheSpydroidSetupSrc)
		bulgluf.Tapped = true

		player.Player.DestroyDeck()
		for range 4 {
			_, err := player.Player.SpawnCard(bulglufTheSpydroidSeedUID, match.DECK)
			require.NoError(t, err)
		}

		topBefore := player.Player.PeekDeck(1)
		require.Len(t, topBefore, 1)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)
		assert.Equal(t, match.SHIELDZONE, topBefore[0].Zone)
		assert.True(t, bulgluf.Tapped)
	})

	t.Run("the new shield is a real shield the opponent can break", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		bulgluf := putCardInBattlezone(t, scn, player.Player, bulglufTheSpydroidUID, bulglufTheSpydroidSetupSrc)
		bulgluf.Tapped = true

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		added := shields[len(shields)-1]
		assert.Equal(t, match.SHIELDZONE, added.Zone)
		assert.Contains(t, player.Player.ShieldMap, added.ID, "it is numbered like every other shield")
	})
}
