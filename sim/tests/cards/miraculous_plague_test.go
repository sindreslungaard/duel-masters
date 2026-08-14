package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	miraculousPlagueUID = "c912b52a-6d6b-4997-ba2d-49468e18eb65"
	plagueCreatureUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	plagueSetupSrc      = "miraculous_plague_test_setup"
)

func TestMiraculousPlague(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(miraculousPlagueUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Miraculous Plague", spell.Name)
		assert.Equal(t, 7, spell.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, spell.Civs)
	})

	t.Run("the opponent keeps one of each pair and loses the other", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		firstCreature := putCardInBattlezone(t, scn, opponent.Player, plagueCreatureUID, plagueSetupSrc)
		secondCreature := putCardInBattlezone(t, scn, opponent.Player, plagueCreatureUID, plagueSetupSrc)

		firstMana, err := opponent.Player.SpawnCard(plagueCreatureUID, match.MANAZONE)
		require.NoError(t, err)
		secondMana, err := opponent.Player.SpawnCard(plagueCreatureUID, match.MANAZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, miraculousPlagueUID)

		// Two creatures on the board, so the caster's choice is forced; the
		// opponent then says which of them survives.
		answerInTurn(t, scn, opponent, firstCreature.ID)

		// Two cards in mana, likewise.
		answerInTurn(t, scn, opponent, firstMana.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, firstCreature.Zone, "the one they kept")
		assert.Equal(t, match.GRAVEYARD, secondCreature.Zone, "the other is destroyed")
		assert.Equal(t, match.HAND, firstMana.Zone)
		assert.Equal(t, match.GRAVEYARD, secondMana.Zone)
	})

	t.Run("a single creature is simply kept", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		only := putCardInBattlezone(t, scn, opponent.Player, plagueCreatureUID, plagueSetupSrc)

		castSpell(t, scn, player, miraculousPlagueUID)
		require.NoError(t, scn.WaitForEventLoop())

		// One creature and no mana: every selection resolves without a choice
		// and nothing is left over to lose.
		assert.Equal(t, match.HAND, only.Zone)
	})

	t.Run("an empty board costs the opponent nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, miraculousPlagueUID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore))
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
