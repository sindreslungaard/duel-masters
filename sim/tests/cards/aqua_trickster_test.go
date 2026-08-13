package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	aquaTricksterUID = "5d09f6bb-3b6a-48ee-be64-0c9017a02708"
	aquaTricksterSrc = "aqua_trickster_test_setup"
)

func TestAquaTrickster(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		trickster := putCardInBattlezone(t, scn, player.Player, aquaTricksterUID, aquaTricksterSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, trickster, "Aqua Trickster", 1000, 2, []string{civ.Water})
		assert.True(t, trickster.HasFamily(family.LiquidPeople))
		assert.True(t, trickster.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it taps a creature on arrival", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, aquaTricksterSrc)

		trickster := spawnForLater(t, player, aquaTricksterUID)
		passTurnToSelf(t, scn, player, opponent)

		chatStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		putIntoPlay(t, scn, player, trickster)

		assert.True(t, theirs.Tapped, "the only legal target is taken without asking")

		// Tapping is otherwise invisible, so the affected player has to be told
		// what happened to their board.
		chat, err := scn.ChatMessages(opponent, chatStart)
		require.NoError(t, err)
		assert.Contains(t, chat, "Immortal Baron, Vorg was tapped by Aqua Trickster")
	})

	t.Run("without the count nothing happens", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, aquaTricksterSrc)

		trickster := spawnForLater(t, player, aquaTricksterUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, trickster)

		assert.False(t, theirs.Tapped)
	})

	t.Run("the arriving creature counts the two already there", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, aquaTricksterSrc)

		trickster := spawnForLater(t, player, aquaTricksterUID)
		passTurnToSelf(t, scn, player, opponent)

		// Exactly at the threshold: two others, and itself makes three.
		require.NoError(t, scn.WaitForEventLoop())
		putIntoPlay(t, scn, player, trickster)

		assert.True(t, theirs.Tapped)
	})
}
