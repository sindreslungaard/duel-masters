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
	engbeltTheSpydroidUID = "08d2dcbc-e643-4576-9745-1317bfa7968e"
	engbeltSetupSrc       = "engbelt_the_spydroid_test_setup"
)

func TestEngbeltTheSpydroid(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		engbelt := putCardInBattlezone(t, scn, player.Player, engbeltTheSpydroidUID, engbeltSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, engbelt, "Engbelt, the Spydroid", 5500, 4, []string{civ.Light})
		assert.True(t, engbelt.HasFamily(family.Soltrooper))
		assert.True(t, engbelt.HasCondition(cnd.Blocker))
		assert.True(t, engbelt.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("it cannot attack the player", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		engbelt := putCardInBattlezone(t, scn, player.Player, engbeltTheSpydroidUID, engbeltSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		_, err := scn.ActionAttackPlayer(player, engbelt.ID)
		require.Error(t, err, "the attack should be rejected")
		assert.False(t, engbelt.Tapped, "a rejected attack does not tap it")
	})

	t.Run("it can still attack a creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		engbelt := putCardInBattlezone(t, scn, player.Player, engbeltTheSpydroidUID, engbeltSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, engbeltSetupSrc)
		victim.Tapped = true

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, engbelt.ID, victim.ID))

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
	})
}
