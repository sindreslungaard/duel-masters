package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tenTonCrunchUID      = "7dbd2f0a-f53f-4c1b-aeaf-6840984c19b6"
	tenTonCrunchEdgeUID  = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	tenTonCrunchBigUID   = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	tenTonCrunchSetupSrc = "ten_ton_crunch_test_setup"
)

func TestTenTonCrunch(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(tenTonCrunchUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Ten-Ton Crunch", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire}, spell.Civs)
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("it destroys a creature with power 3000 or less", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		edge := putCardInBattlezone(t, scn, opponent.Player, tenTonCrunchEdgeUID, tenTonCrunchSetupSrc)
		big := putCardInBattlezone(t, scn, opponent.Player, tenTonCrunchBigUID, tenTonCrunchSetupSrc)

		castSpell(t, scn, player, tenTonCrunchUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Only one creature is within range, so it is taken without asking.
		assert.Equal(t, match.GRAVEYARD, edge.Zone, "exactly 3000 is within range")
		assert.Equal(t, match.BATTLEZONE, big.Zone)
	})

	t.Run("it cannot destroy its caster's own creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		own := putCardInBattlezone(t, scn, player.Player, tenTonCrunchEdgeUID, tenTonCrunchSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, tenTonCrunchEdgeUID, tenTonCrunchSetupSrc)

		castSpell(t, scn, player, tenTonCrunchUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, own.Zone)
	})

	t.Run("nothing in range opens no prompt", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		big := putCardInBattlezone(t, scn, opponent.Player, tenTonCrunchBigUID, tenTonCrunchSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, tenTonCrunchUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.BATTLEZONE, big.Zone)
	})
}
