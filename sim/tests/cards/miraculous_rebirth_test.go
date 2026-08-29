package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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
	// Emeral (cost 2, water): when summoned, may add a card from hand to
	// shields face down. Used to prove that effect can no longer see the
	// resolving Rebirth itself in hand.
	rebirthEmeralUID = "d96cdbba-0e60-4d8e-9394-4830b11f559d"
	rebirthSetupSrc  = "miraculous_rebirth_test_setup"
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

	t.Run("a creature it fetches cannot use Rebirth itself for its own effect", func(t *testing.T) {
		// Reported bug: casting Rebirth and fetching Emeral let Emeral's own
		// "add a card from your hand to your shields" effect put Rebirth
		// itself into the shields, because Rebirth was still sitting in hand
		// while its own search-and-summon text was still resolving.
		scn, player, opponent := setupDuel(t)
		victim := putCardInBattlezone(t, scn, opponent.Player, rebirthVictimUID, rebirthSetupSrc)

		// Emeral (cost 2) buried in the deck among cost-4 filler, so it is
		// the only creature matching the victim's cost.
		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}
		emeral, err := player.Player.SpawnCard(rebirthEmeralUID, match.DECK)
		require.NoError(t, err)
		for range 4 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, miraculousRebirthUID)
		answerInTurn(t, scn, player, emeral.ID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected Emeral's own shield prompt to be open")
		for _, offered := range action.Cards {
			assert.NotEqual(t, spell.ID, offered.CardID, "Miraculous Rebirth must not be offered to Emeral's own effect while it is resolving")
		}

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.BATTLEZONE, emeral.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("lets a creature it fetches recur Rebirth itself from the graveyard", func(t *testing.T) {
		// Reported bug: "if I rebirth a 5-drop and pull out my own Phal, I
		// should be able to recur that same Rebirth but currently can't."
		scn, player, opponent := setupDuel(t)
		victim := putCardInBattlezone(t, scn, opponent.Player, jagilaUID, rebirthSetupSrc)

		// Phal Eega (cost 5) buried in the deck among cost-4 filler, so it is
		// the only creature matching the victim's cost.
		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}
		phalEega, err := player.Player.SpawnCard(phalEegaUID, match.DECK)
		require.NoError(t, err)
		for range 4 {
			_, err := player.Player.SpawnCard(rebirthBigUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, miraculousRebirthUID)
		// Rebirth was spawned into hand fresh for this cast and never saw an
		// untap step, which is otherwise what rebuilds cnd.Spell; set it
		// directly so Phal Eega's own "return a spell" filter can recognize
		// it once it reaches the graveyard.
		spell.AddCondition(cnd.Spell, nil, spell.ID)

		answerInTurn(t, scn, player, phalEega.ID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected Phal Eega's own effect to prompt to return a spell from the graveyard")
		offeredCardIDs := make([]string, 0, len(action.Cards))
		for _, offered := range action.Cards {
			offeredCardIDs = append(offeredCardIDs, offered.CardID)
		}
		assert.Contains(t, offeredCardIDs, spell.ID, "Miraculous Rebirth should already be in the graveyard by the time Phal Eega's own ability resolves")

		answerInTurn(t, scn, player, spell.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.HAND, spell.Zone)
		assert.Equal(t, match.BATTLEZONE, phalEega.Zone)
	})
}
