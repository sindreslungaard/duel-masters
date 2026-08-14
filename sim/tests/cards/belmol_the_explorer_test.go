package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	belmolTheExplorerUID      = "d0b9d738-bd33-4dd1-b6dc-234897f52266"
	belmolTheExplorerSetupSrc = "belmolTheExplorer_test_setup"
	belmolAttackerUID         = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
)

func TestBelmolTheExplorer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, belmolTheExplorerUID, belmolTheExplorerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Belmol, the Explorer", 3500, 2, []string{civ.Light})
		assert.True(t, card.HasFamily(family.Gladiator))
		assert.True(t, card.HasCondition(cnd.Blocker))
		assert.True(t, card.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("it blocks whenever it is able", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, belmolTheExplorerUID, belmolTheExplorerSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, belmolAttackerUID, belmolTheExplorerSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// It blocks without being asked, and 3500 beats 2000.
		assert.Equal(t, match.GRAVEYARD, attacker.Zone)
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone, "the attack never reached the shields")
	})
}
