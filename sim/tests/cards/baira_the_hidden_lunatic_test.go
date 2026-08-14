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
	bairaTheHiddenLunaticUID = "b12e23f7-70cc-4aad-b78f-e71dde4f783e"
	bairaSetupSrc            = "baira_the_hidden_lunatic_test_setup"
)

func TestBairaTheHiddenLunatic(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baira := putCardInBattlezone(t, scn, player.Player, bairaTheHiddenLunaticUID, bairaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, baira, "Baira, the Hidden Lunatic", 5000, 3, []string{civ.Darkness})
		assert.True(t, baira.HasFamily(family.PandorasBox))
		assert.True(t, baira.HasCondition(cnd.Blocker))
		assert.True(t, baira.HasCondition(cnd.CantAttackPlayers))
		assert.True(t, baira.HasCondition(cnd.CantAttackCreatures))
		assert.True(t, baira.HasCondition(cnd.DestroyAfterBattle))
	})

	t.Run("it cannot attack at all", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baira := putCardInBattlezone(t, scn, player.Player, bairaTheHiddenLunaticUID, bairaSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, bairaSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		_, err := scn.ActionAttackPlayer(player, baira.ID)
		require.Error(t, err)

		require.Error(t, scn.ActionAttackCreature(player, baira.ID, victim.ID))
		assert.Equal(t, match.BATTLEZONE, victim.Zone)
		assert.False(t, baira.Tapped)
	})

	t.Run("winning a block still destroys it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baira := putCardInBattlezone(t, scn, player.Player, bairaTheHiddenLunaticUID, bairaSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, bairaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		attackStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		// The defender is asked whether to block.
		_, err = scn.WaitForAction(player, attackStart)
		require.NoError(t, err)
		answerInTurn(t, scn, player, baira.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "5000 beats 2000")
		assert.Equal(t, match.GRAVEYARD, baira.Zone, "it destroys itself after the battle")
	})
}
