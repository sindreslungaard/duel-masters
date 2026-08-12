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
	yulianaChannelerOfSunsUID = "e255f07a-4b76-40a6-8441-faf8c7b6cd41"
	yulianaRemovalSpellUID    = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	yulianaSetupSrc           = "yuliana_channeler_of_suns_test_setup"
)

func TestYulianaChannelerOfSuns(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		yuliana := putCardInBattlezone(t, scn, player.Player, yulianaChannelerOfSunsUID, yulianaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, yuliana, "Yuliana, Channeler of Suns", 3000, 3, []string{civ.Light})
		assert.True(t, yuliana.HasFamily(family.MechaDelSol))
		assert.True(t, yuliana.HasCondition(cnd.Blocker))
		assert.True(t, yuliana.HasCondition(cnd.CantBeSelectedByOpp))
		assert.True(t, yuliana.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("the opponent cannot choose it with an effect", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		yuliana := putCardInBattlezone(t, scn, player.Player, yulianaChannelerOfSunsUID, yulianaSetupSrc)
		other := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, yulianaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		// Terror Pit destroys a creature of the caster's choice. Yuliana must
		// not be offered, leaving the other creature as the only target.
		castSpell(t, scn, opponent, yulianaRemovalSpellUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, yuliana.Zone)
		assert.Equal(t, match.GRAVEYARD, other.Zone)
	})

	t.Run("it can still be attacked", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		yuliana := putCardInBattlezone(t, scn, player.Player, yulianaChannelerOfSunsUID, yulianaSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, beratchaBigCreatureUID, yulianaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		yuliana.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(opponent, attacker.ID, yuliana.ID))

		assert.Equal(t, match.GRAVEYARD, yuliana.Zone, "unchoosable is not untouchable")
	})
}
