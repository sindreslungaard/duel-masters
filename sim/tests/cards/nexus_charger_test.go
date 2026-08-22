package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	nexusChargerUID      = "70cac25c-bf54-4934-8180-cf18867e9d7d"
	nexusChargerSetupSrc = "nexus_charger_test_setup"
)

func TestNexusCharger(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(nexusChargerUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Nexus Charger", spell.Name)
		assert.Equal(t, 6, spell.ManaCost)
		assert.Equal(t, []string{civ.Light}, spell.Civs)
		assert.Equal(t, []string{civ.Light}, spell.ManaRequirement)
	})

	t.Run("cannot choose itself, puts a different hand card into the shields, and goes to the mana zone instead of the graveyard", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, nexusChargerUID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the shield-selection prompt to be open")
		for _, offered := range action.Cards {
			assert.NotEqual(t, spell.ID, offered.CardID, "Nexus Charger must not be an option for its own effect while it is resolving")
		}

		filler, err := scn.FindCard(player.Player, match.HAND, immortalBaronVorgUID)
		require.NoError(t, err)

		answerInTurn(t, scn, player, filler.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.SHIELDZONE, filler.Zone)
		assert.Equal(t, match.MANAZONE, spell.Zone, "a charger goes to the mana zone instead of the graveyard")

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)
	})

	t.Run("with no other card in hand, the shield effect does nothing but it still becomes mana", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		emptyHand(t, player, nexusChargerSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, nexusChargerUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.MANAZONE, spell.Zone)
	})
}
