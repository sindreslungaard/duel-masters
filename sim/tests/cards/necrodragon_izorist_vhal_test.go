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
	necrodragonIzoristVhalUID      = "56486d95-811d-4a77-a21b-dd773eb67d97"
	necrodragonIzoristVhalSetupSrc = "necrodragon_izorist_vhal_test_setup"
	izoristVhalDarknessCreatureUID = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul (darkness)
	izoristVhalOtherCivCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire)
)

func TestNecrodragonIzoristVhal(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(izoristVhalDarknessCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		// Let the graveyard creature pick up cnd.Creature before Vhal ever inspects
		// it, matching a creature that actually died in a prior turn.
		passTurnToSelf(t, scn, player, opponent)

		card := putCardInBattlezone(t, scn, player.Player, necrodragonIzoristVhalUID, necrodragonIzoristVhalSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Necrodragon Izorist Vhal", card.Name)
		assert.Equal(t, 6, card.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, card.Civs)
		assert.Equal(t, []string{civ.Darkness}, card.ManaRequirement)
		assert.True(t, card.HasFamily(family.ZombieDragon))
		assert.Equal(t, 2000, card.Power)
	})

	t.Run("survives ending the turn it was summoned on with enough darkness creatures in the graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(izoristVhalDarknessCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)
		_, err = player.Player.SpawnCard(izoristVhalDarknessCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		// Warm up cnd.Creature on the graveyard occupants, then summon Vhal and end
		// the turn in which it entered the battle zone. This is the reported bug:
		// Vhal used to destroy itself here despite the graveyard having plenty of
		// darkness creatures, because BeginTurnStep/UntapManaEvent/UntapStep run
		// with cnd.Creature cleared match-wide and not yet rebuilt.
		passTurnToSelf(t, scn, player, opponent)

		card := putCardInBattlezone(t, scn, player.Player, necrodragonIzoristVhalUID, necrodragonIzoristVhalSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, match.BATTLEZONE, card.Zone)
		assert.Equal(t, 4000, card.Power)
	})

	t.Run("destroyed when there are no darkness creatures in the graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(izoristVhalOtherCivCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		card := putCardInBattlezone(t, scn, player.Player, necrodragonIzoristVhalUID, necrodragonIzoristVhalSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, match.GRAVEYARD, card.Zone)
	})

	t.Run("gains double breaker at 6000 power or more", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		for range 3 {
			_, err := player.Player.SpawnCard(izoristVhalDarknessCreatureUID, match.GRAVEYARD)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		card := putCardInBattlezone(t, scn, player.Player, necrodragonIzoristVhalUID, necrodragonIzoristVhalSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 6000, card.Power)
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})
}
