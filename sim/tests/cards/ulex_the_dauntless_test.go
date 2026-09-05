package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ulexTheDauntlessUID = "48799b4f-c7cd-4e06-9430-8a3782fd5084"
	ulexTechnoTotemUID  = "a1ef4e4e-d8c9-4b33-bfb0-3f8f7d021ff5" // Techno Totem, tap ability taps one opposing creature
	ulexAllyUID         = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000), a normal tap target
	ulexSetupSrc        = "ulex_the_dauntless_test_setup"
)

func TestUlexTheDauntless(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		ulex := putCardInBattlezone(t, scn, player.Player, ulexTheDauntlessUID, ulexSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Ulex, the Dauntless", ulex.Name)
		assert.Equal(t, 3000, ulex.Power)
		assert.Equal(t, 3, ulex.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, ulex.Civs)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, ulex.ManaRequirement)
		assert.True(t, ulex.IsMulticolored())
		assert.True(t, ulex.HasFamily(family.SpiritQuartz))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(ulexTheDauntlessUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, ulexSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("a targeted tap effect from the opponent fizzles when it is the only target", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		totem := putCardInBattlezone(t, scn, player.Player, ulexTechnoTotemUID, ulexSetupSrc)
		ulex := putCardInBattlezone(t, scn, opponent.Player, ulexTheDauntlessUID, ulexSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionUseTapAbility(player, totem.ID))

		assert.False(t, ulex.Tapped, "its owner's opponent can't tap it")
		assert.True(t, totem.Tapped, "using the tap ability still taps the totem itself")
	})

	t.Run("a targeted tap effect from the opponent still taps a different creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		totem := putCardInBattlezone(t, scn, player.Player, ulexTechnoTotemUID, ulexSetupSrc)
		ulex := putCardInBattlezone(t, scn, opponent.Player, ulexTheDauntlessUID, ulexSetupSrc)
		ally := putCardInBattlezone(t, scn, opponent.Player, ulexAllyUID, ulexSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		// Ulex is the only other creature the totem's tap ability could pick,
		// so filtering it out leaves a single unambiguous legal target and the
		// prompt is bypassed, just like Techno Totem's own single-target case.
		require.NoError(t, scn.ActionUseTapAbility(player, totem.ID))

		assert.False(t, ulex.Tapped, "its owner's opponent can't tap it")
		assert.True(t, ally.Tapped, "the effect still reaches a creature without the protection")
	})
}
