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
	uberdragonZaschackUID        = "ec5bdb88-895a-4ca8-8f83-0cd99d24c765"
	uberdragonZaschackBaseUID    = "31683921-16fc-4d4c-bb77-3225d10f7366" // Ultimate Dragon (Armored Dragon)
	uberdragonZaschackOtherUID   = "e27ac147-3d7d-42d7-a3a4-2e3a1eccdb2c" // Boltail Dragon (Armored Dragon)
	uberdragonZaschackNonBaseUID = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (Rock Beast)
	uberdragonZaschackSrc        = "uberdragon_zaschack_test_setup"
)

func TestUberdragonZaschack(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		zaschack := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackUID, uberdragonZaschackSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, zaschack, "Uberdragon Zaschack", 11000, 9, []string{civ.Fire})
		assert.True(t, zaschack.HasFamily(family.ArmoredDragon))
		assert.True(t, zaschack.HasCondition(cnd.Evolution))
		assert.False(t, zaschack.HasCondition(cnd.ShieldBreakModifier), "no bonus while alone")
		assert.Equal(t, 11000, scn.Match.GetPower(zaschack, false))
	})

	t.Run("it evolves onto an Armored Dragon", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackBaseUID, uberdragonZaschackSrc)

		zaschack, err := player.Player.SpawnCard(uberdragonZaschackUID, match.HAND)
		require.NoError(t, err)
		for range 9 {
			_, err := player.Player.SpawnCard(uberdragonZaschackUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, zaschack.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, zaschack.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone, "the base goes under the evolution")
	})

	t.Run("it cannot evolve onto an unrelated creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		wrongBase := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackNonBaseUID, uberdragonZaschackSrc)

		zaschack, err := player.Player.SpawnCard(uberdragonZaschackUID, match.HAND)
		require.NoError(t, err)
		for range 9 {
			_, err := player.Player.SpawnCard(uberdragonZaschackUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.Error(t, scn.ActionPlayCard(player, zaschack.ID), "there is no legal base to evolve from")
		assert.Equal(t, match.HAND, zaschack.Zone)
		assert.Equal(t, match.BATTLEZONE, wrongBase.Zone)
	})

	t.Run("breaks one more shield for each other own Armored Dragon", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		zaschack := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackUID, uberdragonZaschackSrc)
		putCardInBattlezone(t, scn, player.Player, uberdragonZaschackBaseUID, uberdragonZaschackSrc)
		putCardInBattlezone(t, scn, player.Player, uberdragonZaschackOtherUID, uberdragonZaschackSrc)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, zaschack.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, action.MinSelections, "1 base shield plus 1 per other Armored Dragon")
		assert.Equal(t, 3, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-3)
	})

	t.Run("breaks a single shield with no other Armored Dragons", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		zaschack := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackUID, uberdragonZaschackSrc)

		action, err := scn.ActionAttackPlayer(player, zaschack.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 4)
	})

	t.Run("ignores non-Armored-Dragons and the opponent's Armored Dragons", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		zaschack := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackUID, uberdragonZaschackSrc)

		putCardInBattlezone(t, scn, player.Player, uberdragonZaschackNonBaseUID, uberdragonZaschackSrc)
		assert.False(t, zaschack.HasCondition(cnd.ShieldBreakModifier), "a Rock Beast is not an Armored Dragon")

		putCardInBattlezone(t, scn, opponent.Player, uberdragonZaschackBaseUID, uberdragonZaschackSrc)
		assert.False(t, zaschack.HasCondition(cnd.ShieldBreakModifier), "only your own creatures count")
	})

	t.Run("drops the crew breaker when it leaves the battle zone", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		zaschack := putCardInBattlezone(t, scn, player.Player, uberdragonZaschackUID, uberdragonZaschackSrc)
		putCardInBattlezone(t, scn, player.Player, uberdragonZaschackBaseUID, uberdragonZaschackSrc)
		require.Equal(t, 1, shieldBreakModifierOf(t, zaschack))

		moved, err := player.Player.MoveCard(zaschack.ID, match.BATTLEZONE, match.GRAVEYARD, uberdragonZaschackSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, zaschack.HasCondition(cnd.ShieldBreakModifier))
	})
}
