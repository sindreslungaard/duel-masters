package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	aquaSkydiverUID      = "ef6d1314-1005-4b5d-9afe-827ffe3ba58a"
	aquaSkydiverKillerID = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (5000)
	aquaSkydiverSetupSrc = "aqua_skydiver_test_setup"
)

func TestAquaSkydiver(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupAquaSkydiverTest(t)
		skydiver := putCardInBattlezone(t, scn, player.Player, aquaSkydiverUID, aquaSkydiverSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Aqua Skydiver", skydiver.Name)
		assert.Equal(t, 1000, skydiver.Power)
		assert.Equal(t, 4, skydiver.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Water}, skydiver.Civs)
		assert.True(t, skydiver.HasFamily(family.LiquidPeople))
		assert.True(t, skydiver.HasCondition(cnd.Blocker))
		assert.True(t, skydiver.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("returns to hand instead of being destroyed", func(t *testing.T) {
		scn, player, _ := setupAquaSkydiverTest(t)
		skydiver := putCardInBattlezone(t, scn, player.Player, aquaSkydiverUID, aquaSkydiverSetupSrc)

		scn.Match.Destroy(skydiver, skydiver, match.DestroyedByMiscAbility)

		assert.Equal(t, match.HAND, skydiver.Zone)
		graveyard, err := player.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		assert.Empty(t, graveyard)
	})

	t.Run("returns to hand after losing a battle", func(t *testing.T) {
		scn, defender, attacker := setupAquaSkydiverTest(t)
		skydiver := putCardInBattlezone(t, scn, defender.Player, aquaSkydiverUID, aquaSkydiverSetupSrc)
		killer := putCardInBattlezone(t, scn, attacker.Player, aquaSkydiverKillerID, aquaSkydiverSetupSrc)

		require.NoError(t, scn.ActionEndTurn(defender))
		skydiver.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(attacker, killer.ID, skydiver.ID))

		assert.Equal(t, match.HAND, skydiver.Zone, "a battle loss is still a destruction")
		assert.Equal(t, match.BATTLEZONE, killer.Zone)
	})

	t.Run("is put into the mana zone tapped", func(t *testing.T) {
		scn, player, _ := setupAquaSkydiverTest(t)
		player.Player.SpawnCard(aquaSkydiverUID, match.HAND)
		skydiver, err := scn.FindCard(player.Player, match.HAND, aquaSkydiverUID)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(skydiver.ID, match.HAND, match.MANAZONE, aquaSkydiverSetupSrc)
		require.NoError(t, err)
		assert.True(t, moved.Tapped)
	})
}

func setupAquaSkydiverTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}
