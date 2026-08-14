package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	gankloakRogueCommandoUID = "e663f082-c6a0-4c07-a19d-4f8287f1bbc7"
	gankloakWaterAllyUID     = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (water)
	gankloakSetupSrc         = "gankloak_rogue_commando_test_setup"
)

func TestGankloakRogueCommando(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gankloak := putCardInBattlezone(t, scn, player.Player, gankloakRogueCommandoUID, gankloakSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, gankloak, "Gankloak, Rogue Commando", 2000, 3, []string{civ.Fire})
		assert.True(t, gankloak.HasFamily(family.Human))
		assert.True(t, gankloak.HasCondition(cnd.SilentSkill))
	})

	t.Run("its silent skill gives the fire creatures double breaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gankloak := putCardInBattlezone(t, scn, player.Player, gankloakRogueCommandoUID, gankloakSetupSrc)
		gankloak.Tapped = true

		fireAlly := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, gankloakSetupSrc)
		waterAlly := putCardInBattlezone(t, scn, player.Player, gankloakWaterAllyUID, gankloakSetupSrc)
		theirFire := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, gankloakSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.True(t, fireAlly.HasCondition(cnd.DoubleBreaker))
		assert.True(t, gankloak.HasCondition(cnd.DoubleBreaker), "it is a fire creature of its own")
		assert.False(t, waterAlly.HasCondition(cnd.DoubleBreaker), "only fire creatures")
		assert.False(t, theirFire.HasCondition(cnd.DoubleBreaker), "only its controller's creatures")
	})

	t.Run("declining grants nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gankloak := putCardInBattlezone(t, scn, player.Player, gankloakRogueCommandoUID, gankloakSetupSrc)
		gankloak.Tapped = true

		fireAlly := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, gankloakSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.False(t, fireAlly.HasCondition(cnd.DoubleBreaker))
	})
}
