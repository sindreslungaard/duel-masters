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
	dolzarkUID = "b78076a1-eecf-4cc1-9bbc-009b90930edf"
	// Bolshack Dragon (Armored Dragon, 6000) is itself a double breaker, so
	// attacks with it break 2 shields.
	dolzarkKinUID    = "0ffdcae3-9db2-401b-8a82-dfad707b83cd"
	dolzarkStrangerU = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	dolzarkWallUID   = "85852774-dd96-4395-8980-eb5b85bf5bfc" // Ferrosaturn, Spectral Knight (blocker, 2000)
	dolzarkWeakUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000, at or under the 5000 threshold)
	dolzarkHugeUID   = "f6ff8845-afc0-4958-8673-fad12058193a" // Bloodwing Mantis (6000, over the 5000 threshold)
	dolzarkSetupSrc  = "super_dragon_machine_dolzark_test_setup"
)

func TestSuperDragonMachineDolzark(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Super Dragon Machine Dolzark", 7000, 6, []string{civ.Fire, civ.Nature})
		assert.True(t, card.HasFamily(family.ArmoredDragon))
		assert.True(t, card.HasFamily(family.EarthDragon))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(dolzarkUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, dolzarkSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("an attack by another dragon lets its controller mana a weak opposing creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, dolzarkKinUID, dolzarkSetupSrc)
		weak := putCardInBattlezone(t, scn, opponent.Player, dolzarkWeakUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		choiceStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))

		_, err = scn.WaitForAction(player, choiceStart)
		require.NoError(t, err, "Dolzark's controller chooses which of the opponent's creatures is manaed")
		require.NoError(t, scn.SubmitAction(player, weak.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.MANAZONE, weak.Zone)
	})

	t.Run("a creature over the power threshold cannot be chosen", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, dolzarkKinUID, dolzarkSetupSrc)
		huge := putCardInBattlezone(t, scn, opponent.Player, dolzarkHugeUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, huge.Zone, "no legal target means no prompt and nothing moves")
	})

	t.Run("its controller may decline the choice", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, dolzarkKinUID, dolzarkSetupSrc)
		weak := putCardInBattlezone(t, scn, opponent.Player, dolzarkWeakUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		choiceStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))

		_, err = scn.WaitForAction(player, choiceStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, weak.Zone, "declining leaves the creature in play")
	})

	t.Run("an attack by a non-dragon does nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, dolzarkStrangerU, dolzarkSetupSrc)
		weak := putCardInBattlezone(t, scn, opponent.Player, dolzarkWeakUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, stranger.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, weak.Zone, "only a dragon race attacker triggers it")
	})

	t.Run("dolzark's own attack does not trigger itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dolzark := putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		weak := putCardInBattlezone(t, scn, opponent.Player, dolzarkWeakUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, dolzark.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, weak.Zone, "\"other\" spares itself, so its own attack does not trigger")
	})

	t.Run("a blocked dragon attack still triggers, since the ability fires on attacking, not on being unblocked", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dolzarkUID, dolzarkSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, dolzarkKinUID, dolzarkSetupSrc)
		wall := putCardInBattlezone(t, scn, opponent.Player, dolzarkWallUID, dolzarkSetupSrc)
		weak := putCardInBattlezone(t, scn, opponent.Player, dolzarkWeakUID, dolzarkSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, wall.HasCondition(cnd.Blocker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		blockStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		choiceStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))

		_, err = scn.WaitForAction(player, choiceStart)
		require.NoError(t, err, "the ability should still trigger even though the attack will end up blocked")
		require.NoError(t, scn.SubmitAction(player, weak.ID))

		_, err = scn.WaitForAction(opponent, blockStart)
		require.NoError(t, err, "the blocker should have been offered")
		require.NoError(t, scn.SubmitAction(opponent, wall.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.MANAZONE, weak.Zone)
	})
}
