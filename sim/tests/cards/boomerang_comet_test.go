package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const boomerangCometUID = "b275fbf0-5355-45ec-b3a8-a956cf898ae6"

func TestBoomerangComet(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(boomerangCometUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Boomerang Comet", spell.Name)
		assert.Equal(t, 6, spell.ManaCost)
		assert.Equal(t, []string{civ.Light}, spell.Civs)
	})

	t.Run("returns a card from mana to hand and itself ends up in the mana zone", func(t *testing.T) {
		// Regression test: this card moves itself from hand to the mana zone
		// with a hand-rolled MoveCard rather than fx.Charger. Before it was
		// updated to move from the graveyard, fx.Spell having already sent it
		// there first made this call fail silently, stranding it in the
		// graveyard instead of the mana zone.
		scn, player, _ := setupDuel(t)

		manaCard, err := player.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, boomerangCometUID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the mana-zone selection prompt to be open")
		offeredCardIDs := make([]string, 0, len(action.Cards))
		for _, offered := range action.Cards {
			offeredCardIDs = append(offeredCardIDs, offered.CardID)
		}
		require.Contains(t, offeredCardIDs, manaCard.ID)

		require.NoError(t, scn.SubmitAction(player, manaCard.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, manaCard.Zone)
		assert.Equal(t, match.MANAZONE, spell.Zone, "Boomerang Comet goes to the mana zone instead of the graveyard")
	})
}
