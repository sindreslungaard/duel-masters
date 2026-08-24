package cards

import (
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression coverage for the "Survivor" keyword-ability sharing mechanic
// (DM-06's "Survivor (Each of your Survivors has this creature's Survivor
// ability.)" cards) in the ordinary case where Q-tronic Omnistrain is not in
// play. These fixed cards must keep behaving exactly as they did before that
// fix among creatures that are printed Survivors on their own.
const (
	survivorSharingSmashHornQUID     = "5ff3c63c-30df-4f5c-acce-796f5b6c2dac" // Smash Horn Q (Survivor: +1000 power, shared)
	survivorSharingBalloonshroomQUID = "a8f4d303-72c0-4b67-b6af-d99d267de0c5" // Balloonshroom Q (Survivor: destroyed -> mana zone instead, shared)
	survivorSharingPromephiusQUID    = "d176b30a-cac6-4249-a78d-18f34b97546b" // Promephius Q (plain printed Survivor, no shared ability of its own)
	survivorSharingNonSurvivorUID    = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, not a Survivor)
	survivorSharingSetupSrc          = "survivor_ability_sharing_test_setup"
)

func TestSurvivorAbilitySharingWithoutOmnistrain(t *testing.T) {
	t.Run("Smash Horn Q shares its +1000 power with another real Survivor", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		smashHornQ := putCardInBattlezone(t, scn, player.Player, survivorSharingSmashHornQUID, survivorSharingSetupSrc)
		promephiusQ := putCardInBattlezone(t, scn, player.Player, survivorSharingPromephiusQUID, survivorSharingSetupSrc)
		nonSurvivor := putCardInBattlezone(t, scn, player.Player, survivorSharingNonSurvivorUID, survivorSharingSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, smashHornQ.Power+1000, scn.Match.GetPower(smashHornQ, false), "keeps its own Survivor ability")
		assert.Equal(t, promephiusQ.Power+1000, scn.Match.GetPower(promephiusQ, false), "a real Survivor without the ability still receives it")
		assert.Equal(t, nonSurvivor.Power, scn.Match.GetPower(nonSurvivor, false), "a non-Survivor never receives it")
	})

	t.Run("the shared bonus leaves when Smash Horn Q leaves", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		smashHornQ := putCardInBattlezone(t, scn, player.Player, survivorSharingSmashHornQUID, survivorSharingSetupSrc)
		promephiusQ := putCardInBattlezone(t, scn, player.Player, survivorSharingPromephiusQUID, survivorSharingSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, promephiusQ.Power+1000, scn.Match.GetPower(promephiusQ, false))

		_, err := player.Player.MoveCard(smashHornQ.ID, match.BATTLEZONE, match.GRAVEYARD, survivorSharingSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, promephiusQ.Power, scn.Match.GetPower(promephiusQ, false))
	})

	t.Run("Balloonshroom Q saves another real Survivor from destruction", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, survivorSharingBalloonshroomQUID, survivorSharingSetupSrc)
		promephiusQ := putCardInBattlezone(t, scn, player.Player, survivorSharingPromephiusQUID, survivorSharingSetupSrc)

		scn.Match.Destroy(promephiusQ, promephiusQ, match.DestroyedByMiscAbility)

		assert.Equal(t, match.MANAZONE, promephiusQ.Zone, "saved by Balloonshroom Q's shared ability")
	})

	t.Run("Balloonshroom Q does not save a non-Survivor from destruction", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, survivorSharingBalloonshroomQUID, survivorSharingSetupSrc)
		nonSurvivor := putCardInBattlezone(t, scn, player.Player, survivorSharingNonSurvivorUID, survivorSharingSetupSrc)

		scn.Match.Destroy(nonSurvivor, nonSurvivor, match.DestroyedByMiscAbility)

		assert.Equal(t, match.GRAVEYARD, nonSurvivor.Zone, "not a Survivor, so no protection")
	})
}
