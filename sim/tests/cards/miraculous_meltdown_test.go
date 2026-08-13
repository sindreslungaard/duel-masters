package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	miraculousMeltdownUID = "3713abe0-27bb-475d-b94e-35d2c2fd3f79"
	meltdownSetupSrc      = "miraculous_meltdown_test_setup"
)

func TestMiraculousMeltdown(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(miraculousMeltdownUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Miraculous Meltdown", spell.Name)
		assert.Equal(t, 6, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, spell.Civs)
	})

	t.Run("it cannot be cast while the shields are level", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		mine, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		theirs, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.Equal(t, len(mine), len(theirs), "both start with the same number")

		spell := seedMeltdown(t, player)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.Error(t, scn.ActionPlayCard(player, spell.ID), "the cast should be refused")

		warnings, err := scn.Warnings(player, warningStart)
		require.NoError(t, err)
		assert.NotEmpty(t, warnings)
		assert.Equal(t, match.HAND, spell.Zone, "it stays in hand")
	})

	t.Run("the mana is not spent on a refused cast", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		spell := seedMeltdown(t, player)

		require.Error(t, scn.ActionPlayCard(player, spell.ID))
		require.NoError(t, scn.WaitForEventLoop())

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		for _, card := range mana {
			assert.False(t, card.Tapped, "no mana should have been paid")
		}
	})

	t.Run("it levels the shields when the opponent is ahead", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		trimShields(t, player, 2, meltdownSetupSrc)

		theirs, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.Len(t, theirs, 5)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		spell := seedMeltdown(t, player)
		require.NoError(t, scn.ActionPlayCard(player, spell.ID))

		// The opponent picks which two survive.
		answerInTurn(t, scn, opponent, theirs[0].ID, theirs[1].ID)
		require.NoError(t, scn.WaitForEventLoop())

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, shields, 2, "they keep as many as the caster has")
		assert.Len(t, hand, len(handBefore)+3, "the rest go to their hand")
		assert.Equal(t, match.SHIELDZONE, theirs[0].Zone)
		assert.Equal(t, match.HAND, theirs[4].Zone)
	})
}

// seedMeltdown puts the spell in hand with the mana to pay for it.
func seedMeltdown(t *testing.T, player *match.PlayerReference) *match.Card {
	t.Helper()

	spell, err := player.Player.SpawnCard(miraculousMeltdownUID, match.HAND)
	require.NoError(t, err)

	for range 6 {
		_, err := player.Player.SpawnCard(miraculousMeltdownUID, match.MANAZONE)
		require.NoError(t, err)
	}

	return spell
}

// trimShields leaves a player with exactly keep shields.
func trimShields(t *testing.T, player *match.PlayerReference, keep int, source string) {
	t.Helper()

	shields, err := player.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)

	for _, shield := range shields[keep:] {
		_, err := player.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, source)
		require.NoError(t, err)
	}
}
