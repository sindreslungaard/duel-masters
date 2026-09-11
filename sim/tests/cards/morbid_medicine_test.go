package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	morbidMedicineUID      = "7b95be5c-d378-4af4-98ce-dc18d2e3d172"
	morbidMedicineSpellUID = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	morbidMedicineSetupSrc = "morbid_medicine_test_setup"
)

func TestMorbidMedicine(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(morbidMedicineUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Morbid Medicine", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, spell.Civs)
	})

	t.Run("it returns up to two creatures from the graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		first, err := player.Player.SpawnCard(immortalBaronVorgUID, match.GRAVEYARD)
		require.NoError(t, err)
		second, err := player.Player.SpawnCard(immortalBaronVorgUID, match.GRAVEYARD)
		require.NoError(t, err)
		third, err := player.Player.SpawnCard(immortalBaronVorgUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, morbidMedicineUID)
		answerInTurn(t, scn, player, first.ID, second.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, first.Zone)
		assert.Equal(t, match.HAND, second.Zone)
		assert.Equal(t, match.GRAVEYARD, third.Zone)
	})

	t.Run("it leaves spells in the graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spellInGrave, err := player.Player.SpawnCard(morbidMedicineSpellUID, match.GRAVEYARD)
		require.NoError(t, err)
		creature, err := player.Player.SpawnCard(immortalBaronVorgUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, morbidMedicineUID)
		answerInTurn(t, scn, player, creature.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, creature.Zone)
		assert.Equal(t, match.GRAVEYARD, spellInGrave.Zone)
	})

	t.Run("an empty graveyard offers nothing selectable", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, morbidMedicineUID)

		// By the time Morbid Medicine's own text resolves, it is itself
		// already sitting in the caster's graveyard (fx.Spell moves a cast
		// spell there immediately, before its own effect body runs), so the
		// return-creatures prompt still opens to show that graveyard - just
		// with nothing selectable in it, since a spell isn't a creature.
		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the return-creatures prompt to be open")
		assert.Empty(t, action.Cards, "nothing should be selectable in an otherwise empty graveyard")

		cancelInTurn(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
