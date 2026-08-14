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
	beratchaTheHiddenGluttonUID = "fb598a3b-7153-4a66-93de-0c7b7157b5ea"
	beratchaBigCreatureUID      = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	beratchaSetupSrc            = "beratcha_the_hidden_glutton_test_setup"
)

func TestBeratchaTheHiddenGlutton(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		beratcha := putCardInBattlezone(t, scn, player.Player, beratchaTheHiddenGluttonUID, beratchaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, beratcha, "Beratcha, the Hidden Glutton", 3000, 5, []string{civ.Darkness})
		assert.True(t, beratcha.HasFamily(family.PandorasBox))
		assert.True(t, beratcha.HasCondition(cnd.Slayer))
	})

	t.Run("slayer takes the winner down with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		beratcha := putCardInBattlezone(t, scn, player.Player, beratchaTheHiddenGluttonUID, beratchaSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, beratchaBigCreatureUID, beratchaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		beratcha.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(opponent, attacker.ID, beratcha.ID))

		assert.Equal(t, match.GRAVEYARD, beratcha.Zone, "4000 beats 3000")
		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "slayer destroys the winner too")
	})
}
