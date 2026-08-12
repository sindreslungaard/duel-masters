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
	hustleBerryUID      = "74be313e-b8ff-4741-88ce-b7f0d6adba45"
	hustleBerrySeedUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	hustleBerrySetupSrc = "hustle_berry_test_setup"
)

func TestHustleBerry(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		berry := putCardInBattlezone(t, scn, player.Player, hustleBerryUID, hustleBerrySetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Hustle Berry", berry.Name)
		assert.Equal(t, 1000, berry.Power)
		assert.Equal(t, 2, berry.ManaCost)
		assert.Equal(t, []string{civ.Nature}, berry.Civs)
		assert.Equal(t, []string{civ.Nature}, berry.ManaRequirement)
		assert.True(t, berry.HasFamily(family.WildVeggies))
		assert.True(t, berry.HasCondition(cnd.SilentSkill))
	})

	t.Run("puts the top card of the deck into the mana zone", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		berry := putCardInBattlezone(t, scn, player.Player, hustleBerryUID, hustleBerrySetupSrc)
		berry.Tapped = true

		player.Player.DestroyDeck()
		for range 4 {
			_, err := player.Player.SpawnCard(hustleBerrySeedUID, match.DECK)
			require.NoError(t, err)
		}

		topBefore := player.Player.PeekDeck(1)
		require.Len(t, topBefore, 1)

		manaBefore, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, len(manaBefore)+1)
		assert.Equal(t, match.MANAZONE, topBefore[0].Zone, "the card that was on top is the one that moved")
		assert.False(t, topBefore[0].Tapped, "the mana it makes is usable the turn it arrives")
		assert.True(t, berry.Tapped)
	})

	t.Run("declining leaves the deck and mana zone alone", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		berry := putCardInBattlezone(t, scn, player.Player, hustleBerryUID, hustleBerrySetupSrc)
		berry.Tapped = true

		manaBefore, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, len(manaBefore))
		assert.False(t, berry.Tapped)
	})

	t.Run("takes only the opponent's turn off, not their deck", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		berry := putCardInBattlezone(t, scn, player.Player, hustleBerryUID, hustleBerrySetupSrc)
		berry.Tapped = true

		opponentDeckBefore, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)
		opponentManaBefore, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		opponentDeck, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)
		opponentMana, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		// Only the opponent's own draw step should have touched their deck.
		assert.Len(t, opponentDeck, len(opponentDeckBefore)-1)
		assert.Len(t, opponentMana, len(opponentManaBefore))
	})
}
