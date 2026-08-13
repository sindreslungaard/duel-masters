package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	miraculousRebirthUID = "a45f2208-2b9a-462e-86a2-4e1d004ae3a1"
	rebirthVictimUID     = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000, cost 2)
	rebirthBigUID        = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000, cost 4)
	rebirthTooBigUID     = "f6ff8845-afc0-4958-8673-fad12058193a" // Bloodwing Mantis (6000)
	rebirthSetupSrc      = "miraculous_rebirth_test_setup"
)

func TestMiraculousRebirth(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(miraculousRebirthUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Miraculous Rebirth", spell.Name)
		assert.Equal(t, 6, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, spell.Civs)
	})

	t.Run("it destroys a creature and fetches one of the same cost", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		victim := putCardInBattlezone(t, scn, opponent.Player, rebirthVictimUID, rebirthSetupSrc)

		// A matching creature buried in the deck, below the top so the draw
		// step cannot take it first.
		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}
		replacement, err := player.Player.SpawnCard(rebirthVictimUID, match.DECK)
		require.NoError(t, err)
		for range 4 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, miraculousRebirthUID)

		// Only one creature to destroy, so that is taken without asking; the
		// search then offers the deck.
		answerInTurn(t, scn, player, replacement.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.BATTLEZONE, replacement.Zone)
		assert.Equal(t, player.Player, replacement.Player, "the fetched creature is the caster's")
	})

	t.Run("the search may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		victim := putCardInBattlezone(t, scn, opponent.Player, rebirthVictimUID, rebirthSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, miraculousRebirthUID)
		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone, "the destruction still happened")
	})

	t.Run("a creature above 5000 power is not a legal target", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tough := putCardInBattlezone(t, scn, opponent.Player, rebirthTooBigUID, rebirthSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, miraculousRebirthUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.BATTLEZONE, tough.Zone)
	})
}
