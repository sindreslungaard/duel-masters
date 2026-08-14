package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	reapAndSowUID      = "dac86d1e-f00c-43a5-a2bc-cc8b9ac377c1"
	reapAndSowSetupSrc = "reap_and_sow_test_setup"
)

func TestReapAndSow(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(reapAndSowUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Reap and Sow", spell.Name)
		assert.Equal(t, 5, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, spell.Civs)
	})

	t.Run("it burns one of the opponent's mana and grows its own", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		theirMana, err := opponent.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)

		player.Player.DestroyDeck()
		for range 6 {
			_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		// Read after the handover: the draw step already took a card off the
		// top, so the card this effect will move is not the one seeded first.
		top := player.Player.PeekDeck(1)
		require.Len(t, top, 1)
		topOfDeck := top[0]

		manaBefore, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, reapAndSowUID)
		require.NoError(t, scn.WaitForEventLoop())

		// The opponent has exactly one card in mana, so it is taken without
		// asking.
		assert.Equal(t, match.GRAVEYARD, theirMana.Zone)
		assert.Equal(t, match.MANAZONE, topOfDeck.Zone, "the top of the deck replaces it on the caster's side")

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, len(manaBefore)+6, "five paid for the spell plus the one it puts back")
	})

	t.Run("an empty opposing mana zone still grows the caster's", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		theirMana, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		require.Empty(t, theirMana)

		player.Player.DestroyDeck()
		for range 6 {
			_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		top := player.Player.PeekDeck(1)
		require.Len(t, top, 1)
		topOfDeck := top[0]

		castSpell(t, scn, player, reapAndSowUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.MANAZONE, topOfDeck.Zone)
	})
}
