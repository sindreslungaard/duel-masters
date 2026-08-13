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
	pharziTheOracleUID = "9359a337-0c7e-46a4-acd2-a5f6fc44e2ff"
	pharziSpellUID     = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	pharziCreatureUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	pharziSetupSrc     = "pharzi_the_oracle_test_setup"
)

func TestPharziTheOracle(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, pharziTheOracleUID, pharziSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Pharzi, the Oracle", 1000, 2, []string{civ.Light})
		assert.True(t, card.HasFamily(family.LightBringer))
	})

	t.Run("dying returns a chosen spell to hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pharzi := putCardInBattlezone(t, scn, player.Player, pharziTheOracleUID, pharziSetupSrc)

		spell, err := player.Player.SpawnCard(pharziSpellUID, match.GRAVEYARD)
		require.NoError(t, err)
		creature, err := player.Player.SpawnCard(pharziCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		// Pharzi has to die on the event loop, so it attacks into something
		// bigger rather than being destroyed from the test goroutine.
		wall := putCardInBattlezone(t, scn, opponent.Player, pharziCreatureUID, pharziSetupSrc)
		wall.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, pharzi.ID, wall.ID))

		answerInTurn(t, scn, player, spell.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, pharzi.Zone, "1000 loses to 2000")
		assert.Equal(t, match.HAND, spell.Zone)
		assert.Equal(t, match.GRAVEYARD, creature.Zone, "a creature is not a spell")
	})

	t.Run("the return may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pharzi := putCardInBattlezone(t, scn, player.Player, pharziTheOracleUID, pharziSetupSrc)

		spell, err := player.Player.SpawnCard(pharziSpellUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		wall := putCardInBattlezone(t, scn, opponent.Player, pharziCreatureUID, pharziSetupSrc)
		wall.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, pharzi.ID, wall.ID))

		cancelInTurn(t, scn, player)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("a graveyard without spells asks nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pharzi := putCardInBattlezone(t, scn, player.Player, pharziTheOracleUID, pharziSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		wall := putCardInBattlezone(t, scn, opponent.Player, pharziCreatureUID, pharziSetupSrc)
		wall.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, pharzi.ID, wall.ID))

		// The event loop only goes idle again if nothing is waiting on an answer.
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, pharzi.Zone)
	})
}
