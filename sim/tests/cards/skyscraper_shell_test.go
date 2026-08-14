package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	skyscraperShellUID = "9c8b4458-dc55-407d-a42a-a29a26bc5d99"
	skyscraperSetup    = "skyscraper_shell_test_setup"
)

func TestSkyscraperShell(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		shell := putCardInBattlezone(t, scn, player.Player, skyscraperShellUID, skyscraperSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, shell, "Skyscraper Shell", 2000, 4, []string{civ.Nature})
		assert.True(t, shell.HasFamily(family.ColonyBeetle))
		assert.True(t, shell.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count the opponent buries one of their own creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, skyscraperSetup)

		shell := spawnForLater(t, player, skyscraperShellUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, shell)

		assert.Equal(t, match.MANAZONE, theirs.Zone)
		assert.Equal(t, opponent.Player, theirs.Player, "it goes to its owner's mana zone")
	})

	t.Run("without the count nothing moves", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, skyscraperSetup)

		shell := spawnForLater(t, player, skyscraperShellUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, shell)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
	})
}
