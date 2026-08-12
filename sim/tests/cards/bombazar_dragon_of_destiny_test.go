package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bombazarUID           = "c94b8b0a-5406-4d24-9ac9-333c5c077ccd"
	bombazarSixThousandID = "f6ff8845-afc0-4958-8673-fad12058193a" // Bloodwing Mantis (6000, attack trigger only)
	bombazarSmallerID     = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	bombazarSetupSrc      = "bombazar_dragon_of_destiny_test_setup"
)

func TestBombazarDragonOfDestiny(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)

		// Kept in hand: putting it into play arms the extra turn, which would
		// take this test somewhere it does not need to go. The keywords are
		// rebuilt in every zone, so they can be read from here.
		bombazar, err := player.Player.SpawnCard(bombazarUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Bombazar, Dragon of Destiny", bombazar.Name)
		assert.Equal(t, 6000, bombazar.Power)
		assert.Equal(t, 7, bombazar.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, bombazar.Civs)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, bombazar.ManaRequirement)
		assert.True(t, bombazar.IsMulticolored())
		assert.True(t, bombazar.HasFamily(family.ArmoredDragon))
		assert.True(t, bombazar.HasFamily(family.EarthDragon))
		assert.True(t, bombazar.HasCondition(cnd.SpeedAttacker))
		assert.True(t, bombazar.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupSilentSkillTest(t)

		card, err := player.Player.SpawnCard(bombazarUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, bombazarSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("entering play destroys every other creature with exactly 6000 power", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)

		ownSix := putCardInBattlezone(t, scn, player.Player, bombazarSixThousandID, bombazarSetupSrc)
		theirSix := putCardInBattlezone(t, scn, opponent.Player, bombazarSixThousandID, bombazarSetupSrc)
		theirSmaller := putCardInBattlezone(t, scn, opponent.Player, bombazarSmallerID, bombazarSetupSrc)

		bombazar := putCardInBattlezone(t, scn, player.Player, bombazarUID, bombazarSetupSrc)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, ownSix.Zone, "its own side is caught too")
		assert.Equal(t, match.GRAVEYARD, theirSix.Zone)
		assert.Equal(t, match.BATTLEZONE, theirSmaller.Zone, "4000 is not 6000")
		assert.Equal(t, match.BATTLEZONE, bombazar.Zone, "\"all other creatures\" spares Bombazar itself")
	})

	t.Run("its controller takes an extra turn and then loses at the end of it", func(t *testing.T) {
		scn, player, _ := setupSilentSkillTest(t)

		bombazar := putCardInBattlezone(t, scn, player.Player, bombazarUID, bombazarSetupSrc)
		require.NoError(t, scn.WaitForEventLoop())
		require.Equal(t, match.BATTLEZONE, bombazar.Zone)

		require.True(t, scn.Match.IsPlayerTurn(player.Player))

		// Ending the turn does not hand it over: the same player starts again.
		require.NoError(t, scn.ActionEndTurn(player))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "the extra turn belongs to the same player")

		// Ending the extra turn is where the promise comes due.
		require.NoError(t, scn.ActionEndTurn(player))

		require.Eventually(t, scn.Match.IsClosed, 2*time.Second, 10*time.Millisecond,
			"the game ends at the end of the extra turn")
	})

	t.Run("the loss still happens if Bombazar has left the battle zone", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)

		bombazar := putCardInBattlezone(t, scn, player.Player, bombazarUID, bombazarSetupSrc)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, scn.Match.IsPlayerTurn(player.Player))

		// Removed during the extra turn: the printed text ties the loss to the
		// turn, not to the creature still being around.
		_, err := player.Player.MoveCard(bombazar.ID, match.BATTLEZONE, match.GRAVEYARD, bombazarSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))

		require.Eventually(t, scn.Match.IsClosed, 2*time.Second, 10*time.Millisecond,
			"leaving play does not call the loss off")

		_ = opponent
	})

	t.Run("the opponent is the winner", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)

		putCardInBattlezone(t, scn, player.Player, bombazarUID, bombazarSetupSrc)
		require.NoError(t, scn.WaitForEventLoop())

		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(player))

		require.Eventually(t, scn.Match.IsClosed, 2*time.Second, 10*time.Millisecond)

		headers, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.Contains(t, headers, "duel_finished", "the opponent is told they won")
	})

	t.Run("only one extra turn is taken", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)

		putCardInBattlezone(t, scn, player.Player, bombazarUID, bombazarSetupSrc)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, scn.Match.IsPlayerTurn(player.Player))

		require.NoError(t, scn.ActionEndTurn(player))
		require.Eventually(t, scn.Match.IsClosed, 2*time.Second, 10*time.Millisecond,
			"the second end of turn ends the game rather than granting another turn")

		_ = opponent
	})
}

var _ = scenario.New
