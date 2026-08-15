package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	spinalParasiteUID       = "46346474-92da-4465-994c-ba48bfcd1537"
	spinalParasiteTapperUID = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (tap ability, no target needed)
	spinalParasiteSetupSrc  = "spinal_parasite_test_setup"
)

func TestSpinalParasite(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, spinalParasiteUID, spinalParasiteSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Spinal Parasite", 2000, 5, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.BrainJacker))
	})

	t.Run("the chosen creature cannot use a tap ability instead of attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, spinalParasiteUID, spinalParasiteSetupSrc)
		tapper := putCardInBattlezone(t, scn, opponent.Player, spinalParasiteTapperUID, spinalParasiteSetupSrc)

		// The opponent's own start-of-turn choice: which of their attack-capable
		// creatures Spinal Parasite forces to attack this turn. With only one
		// legal candidate the selection resolves without a prompt.
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		// Official ruling (Slime Veil, the same "attacks if able" wording):
		// "On your turn, when you have the option either to attack with your
		// creature or to use its tap ability... you must attack with it."
		require.Error(t, scn.ActionUseTapAbility(opponent, tapper.ID), "the tap ability should be refused")
		assert.False(t, tapper.Tapped, "a refused ability does not tap the creature")

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, tapper.ID)
		require.NoError(t, err, "attacking with it instead is allowed")
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "attacking satisfied the requirement")
	})

	t.Run("does not restrict a creature it did not choose", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, spinalParasiteUID, spinalParasiteSetupSrc)
		tapper := putCardInBattlezone(t, scn, opponent.Player, spinalParasiteTapperUID, spinalParasiteSetupSrc)
		untouched := putCardInBattlezone(t, scn, opponent.Player, spinalParasiteTapperUID, spinalParasiteSetupSrc)

		questionStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		action, err := scn.WaitForAction(opponent, questionStart)
		require.NoError(t, err)
		require.Len(t, action.Cards, 2, "both attack-capable creatures are offered")
		answerInTurn(t, scn, opponent, tapper.ID)

		// The other tap ability is untouched; it may still be used freely.
		require.NoError(t, scn.ActionUseTapAbility(opponent, untouched.ID))
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, untouched.Tapped)

		require.Error(t, scn.ActionUseTapAbility(opponent, tapper.ID), "the chosen creature is still forced to attack")
	})
}
