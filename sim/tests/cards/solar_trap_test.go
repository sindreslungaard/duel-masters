package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	solarTrapUID      = "e67f1a65-476f-4bf9-8ea7-35afe19a877d"
	solarTrapSetupSrc = "solar_trap_test_setup"
)

func TestSolarTrap(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(solarTrapUID, match.HAND)
		require.NoError(t, err)
		_ = scn

		assert.Equal(t, "Solar Trap", spell.Name)
		assert.Equal(t, 1, spell.ManaCost)
		assert.Equal(t, []string{civ.Light}, spell.Civs)
		assert.Equal(t, []string{civ.Light}, spell.ManaRequirement)
	})

	t.Run("it taps one of the opponent's creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, solarTrapSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, solarTrapSetupSrc)

		spell := castSpell(t, scn, player, solarTrapUID)

		answerInTurn(t, scn, player, second.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, second.Tapped)
		assert.False(t, first.Tapped, "only the chosen creature is tapped")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("it cannot tap its caster's own creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		own := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, solarTrapSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, solarTrapSetupSrc)

		castSpell(t, scn, player, solarTrapUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, theirs.Tapped, "the only legal target is taken without asking")
		assert.False(t, own.Tapped)
	})

	t.Run("an empty opposing battle zone opens no prompt", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		theirs, err := opponent.Player.Container(match.BATTLEZONE)
		require.NoError(t, err)
		require.Empty(t, theirs)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, solarTrapUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
