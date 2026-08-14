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
	quillspikeRumblerUID = "ffbfa50c-bda4-408d-a2b0-64f940f0c305"
	quillspikeVictimUID  = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	quillspikeSetupSrc   = "quillspike_rumbler_test_setup"
)

func TestQuillspikeRumbler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		rumbler := putCardInBattlezone(t, scn, player.Player, quillspikeRumblerUID, quillspikeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, rumbler, "Quillspike Rumbler", 3000, 4, []string{civ.Nature})
		assert.True(t, rumbler.HasFamily(family.BeastFolk))
		assert.Equal(t, 3000, scn.Match.GetPower(rumbler, false))
	})

	t.Run("attacking a creature wins the battle it would otherwise tie", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		rumbler := putCardInBattlezone(t, scn, player.Player, quillspikeRumblerUID, quillspikeSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, quillspikeVictimUID, quillspikeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, rumbler.ID, victim.ID))

		assert.Equal(t, match.GRAVEYARD, victim.Zone, "6000 beats 3000")
		assert.Equal(t, match.BATTLEZONE, rumbler.Zone, "without the bonus this would have been a tie")
		assert.Equal(t, 6000, scn.Match.GetPower(rumbler, false), "the bonus lasts the rest of the turn")
	})

	t.Run("the bonus is gone by the next turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		rumbler := putCardInBattlezone(t, scn, player.Player, quillspikeRumblerUID, quillspikeSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, quillspikeVictimUID, quillspikeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, rumbler.ID, victim.ID))
		require.Equal(t, 6000, scn.Match.GetPower(rumbler, false))

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 3000, scn.Match.GetPower(rumbler, false))
	})

	t.Run("attacking the player does not grant the bonus", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		rumbler := putCardInBattlezone(t, scn, player.Player, quillspikeRumblerUID, quillspikeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, rumbler.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		assert.Equal(t, 3000, scn.Match.GetPower(rumbler, false), "the printed trigger is attacking a creature")
	})
}
