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
	warlordAilzoniusUID = "e1bd40c0-3c76-4854-a967-302e6a6706b3"
	warlordBaseUID      = "da51845c-4a6b-4c36-9c7d-fbb654ba2aa2" // Kanesill, the Explorer (Gladiator)
	warlordSetupSrc     = "warlord_ailzonius_test_setup"
)

func TestWarlordAilzonius(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		warlord := putCardInBattlezone(t, scn, player.Player, warlordAilzoniusUID, warlordSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, warlord, "Warlord Ailzonius", 8000, 5, []string{civ.Light})
		assert.True(t, warlord.HasFamily(family.Gladiator))
		assert.True(t, warlord.HasCondition(cnd.Evolution))
		assert.True(t, warlord.HasCondition(cnd.DoubleBreaker))
		assert.True(t, warlord.HasCondition(cnd.CantBeSelectedByOpp))
	})

	t.Run("the opponent cannot choose it with an effect", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		warlord := putCardInBattlezone(t, scn, player.Player, warlordAilzoniusUID, warlordSetupSrc)
		other := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, warlordSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		castSpell(t, scn, opponent, yulianaRemovalSpellUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, warlord.Zone)
		assert.Equal(t, match.GRAVEYARD, other.Zone)
	})

	t.Run("it breaks two shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		warlord := putCardInBattlezone(t, scn, player.Player, warlordAilzoniusUID, warlordSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shieldsBefore), 2)

		_, err = scn.ActionAttackPlayer(player, warlord.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldsBefore[0].ID, shieldsBefore[1].ID))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)-2)
	})

	t.Run("it evolves onto a Gladiator", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, warlordBaseUID, warlordSetupSrc)

		warlord, err := player.Player.SpawnCard(warlordAilzoniusUID, match.HAND)
		require.NoError(t, err)
		for range 5 {
			_, err := player.Player.SpawnCard(warlordAilzoniusUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, warlord.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, warlord.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone)
	})
}
