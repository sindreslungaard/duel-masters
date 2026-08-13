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
	bonfireLizardUID = "aa0688bf-5a62-4a62-a613-0b4880c159de"
	bonfireSetup     = "bonfire_lizard_test_setup"
)

func TestBonfireLizard(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, bonfireLizardUID, bonfireSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lizard, "Bonfire Lizard", 4000, 6, []string{civ.Fire})
		assert.True(t, lizard.HasFamily(family.MeltWarrior))
		assert.True(t, lizard.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it clears up to two blockers", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		first := putCardInBattlezone(t, scn, opponent.Player, waveStrikerSmallBlockerUID, bonfireSetup)
		second := putCardInBattlezone(t, scn, opponent.Player, waveStrikerSmallBlockerUID, bonfireSetup)
		plain := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, bonfireSetup)

		lizard := summonBonfireLizard(t, scn, player)
		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))

		answerInTurn(t, scn, player, first.ID, second.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, first.Zone)
		assert.Equal(t, match.GRAVEYARD, second.Zone)
		assert.Equal(t, match.BATTLEZONE, plain.Zone, "only blockers are eligible")
	})

	t.Run("up to two means the choice can be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		blocker := putCardInBattlezone(t, scn, opponent.Player, waveStrikerSmallBlockerUID, bonfireSetup)

		lizard := summonBonfireLizard(t, scn, player)
		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, blocker.Zone)
	})

	t.Run("without the count no blocker is touched", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		blocker := putCardInBattlezone(t, scn, opponent.Player, waveStrikerSmallBlockerUID, bonfireSetup)

		lizard := summonBonfireLizard(t, scn, player)
		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionPlayCard(player, lizard.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, blocker.Zone)
	})
}

// summonBonfireLizard puts the card in hand with the mana to pay for it. Its
// ability opens a prompt as it arrives, so it has to be summoned through the
// event loop rather than moved directly from the test goroutine, which would
// leave that goroutine waiting on a prompt only it could answer.
func summonBonfireLizard(t *testing.T, _ *scenario.TestScenario, player *match.PlayerReference) *match.Card {
	t.Helper()

	lizard := spawnForLater(t, player, bonfireLizardUID)

	for range 6 {
		_, err := player.Player.SpawnCard(bonfireLizardUID, match.MANAZONE)
		require.NoError(t, err)
	}

	return lizard
}
