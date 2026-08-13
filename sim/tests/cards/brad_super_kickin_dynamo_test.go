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
	bradSuperKickinDynamoUID        = "bf2605e3-3b96-48ee-a73c-74f15ee8fed5"
	bradSuperKickinDynamoBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker)
	bradSuperKickinDynamoPlainUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	bradSuperKickinDynamoSetupSrc   = "brad_super_kickin_dynamo_test_setup"
)

func TestBradSuperKickinDynamo(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		brad := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoUID, bradSuperKickinDynamoSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Brad, Super Kickin' Dynamo", brad.Name)
		assert.Equal(t, 2000, brad.Power)
		assert.Equal(t, 3, brad.ManaCost)
		assert.Equal(t, []string{civ.Fire}, brad.Civs)
		assert.Equal(t, []string{civ.Fire}, brad.ManaRequirement)
		assert.True(t, brad.HasFamily(family.Human))
		assert.True(t, brad.HasCondition(cnd.SilentSkill))
	})

	t.Run("destroys an opposing blocker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		brad := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoUID, bradSuperKickinDynamoSetupSrc)
		brad.Tapped = true

		blocker := putCardInBattlezone(t, scn, opponent.Player, bradSuperKickinDynamoBlockerUID, bradSuperKickinDynamoSetupSrc)
		plain := putCardInBattlezone(t, scn, opponent.Player, bradSuperKickinDynamoPlainUID, bradSuperKickinDynamoSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, blocker.HasCondition(cnd.Blocker), "the target has to actually be a blocker")

		useSilentSkill(t, scn, player)

		// Exactly one blocker, and the choice is mandatory, so the engine takes
		// it without opening a prompt.
		assert.Equal(t, match.GRAVEYARD, blocker.Zone)
		assert.Equal(t, match.BATTLEZONE, plain.Zone, "a creature without blocker is not a legal target")
		assert.True(t, brad.Tapped)
	})

	t.Run("asks which blocker when there are several", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		brad := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoUID, bradSuperKickinDynamoSetupSrc)
		brad.Tapped = true

		first := putCardInBattlezone(t, scn, opponent.Player, bradSuperKickinDynamoBlockerUID, bradSuperKickinDynamoSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, bradSuperKickinDynamoBlockerUID, bradSuperKickinDynamoSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, second.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, second.Zone)
		assert.Equal(t, match.BATTLEZONE, first.Zone, "only the chosen blocker is destroyed")
	})

	t.Run("no blocker means nothing happens", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		brad := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoUID, bradSuperKickinDynamoSetupSrc)
		brad.Tapped = true

		plain := putCardInBattlezone(t, scn, opponent.Player, bradSuperKickinDynamoPlainUID, bradSuperKickinDynamoSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, plain.Zone)
		assert.True(t, brad.Tapped)
	})

	t.Run("it cannot destroy its controller's own blocker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		brad := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoUID, bradSuperKickinDynamoSetupSrc)
		brad.Tapped = true

		ownBlocker := putCardInBattlezone(t, scn, player.Player, bradSuperKickinDynamoBlockerUID, bradSuperKickinDynamoSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, ownBlocker.Zone)
	})
}
