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
	jabahasAutomatonUID = "d2936cbe-a710-45ff-bc6a-6601fd40f91e"
	jabahasBaseUID      = "1b786e62-ffbc-4694-a8cd-8dd48f8e18fd" // Vorg's Engine (Xenoparts)
	jabahasNonBaseUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human)
	jabahasAutomatonSrc = "jabahas_automaton_test_setup"
	jabahasManaUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
)

func TestJabahasAutomaton(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		automaton := putCardInBattlezone(t, scn, player.Player, jabahasAutomatonUID, jabahasAutomatonSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, automaton, "Jabaha's Automaton", 6000, 5, []string{civ.Fire})
		assert.True(t, automaton.HasFamily(family.Xenoparts))
		assert.True(t, automaton.HasCondition(cnd.Evolution))
		assert.True(t, automaton.HasCondition(cnd.DoubleBreaker))
		assert.True(t, automaton.HasCondition(cnd.PowerAttacker))

		assert.Equal(t, 6000, scn.Match.GetPower(automaton, false))
		assert.Equal(t, 10000, scn.Match.GetPower(automaton, true), "power attacker +4000 while attacking")
	})

	t.Run("it evolves onto a Xenoparts creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, jabahasBaseUID, jabahasAutomatonSrc)

		automaton, err := player.Player.SpawnCard(jabahasAutomatonUID, match.HAND)
		require.NoError(t, err)
		for range 5 {
			_, err := player.Player.SpawnCard(jabahasManaUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, automaton.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, automaton.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone, "the base goes under the evolution")
	})

	t.Run("it cannot evolve onto an unrelated creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		wrongBase := putCardInBattlezone(t, scn, player.Player, jabahasNonBaseUID, jabahasAutomatonSrc)

		automaton, err := player.Player.SpawnCard(jabahasAutomatonUID, match.HAND)
		require.NoError(t, err)
		for range 5 {
			_, err := player.Player.SpawnCard(jabahasManaUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.Error(t, scn.ActionPlayCard(player, automaton.ID), "there is no legal base to evolve from")
		assert.Equal(t, match.HAND, automaton.Zone)
		assert.Equal(t, match.BATTLEZONE, wrongBase.Zone)
	})
}
