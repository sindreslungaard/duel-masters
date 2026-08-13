package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hydroozeUID         = "5ccf0dd9-4dd3-4094-ab00-62c4f265cfd1"
	hydroozeKinUID      = "74ce4102-f6af-467e-8d81-464ac9ebe25e" // Sopian (Cyber Lord, 2000)
	hydroozeStrangerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	hydroozeWallUID     = "85852774-dd96-4395-8980-eb5b85bf5bfc" // Ferrosaturn, Spectral Knight (blocker, 2000)
	hydroozeSetupSrc    = "hydrooze_the_mutant_emperor_test_setup"
)

func TestHydroozeTheMutantEmperor(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, hydroozeUID, hydroozeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Hydrooze, the Mutant Emperor", 5000, 4, []string{civ.Water, civ.Darkness})
		assert.True(t, card.HasFamily(family.CyberLord))
		assert.True(t, card.HasFamily(family.Hedrian))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(hydroozeUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, hydroozeSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("its own other kin get +2000", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hydrooze := putCardInBattlezone(t, scn, player.Player, hydroozeUID, hydroozeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, hydroozeKinUID, hydroozeSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, hydroozeStrangerUID, hydroozeSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, hydroozeKinUID, hydroozeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		assert.Equal(t, stranger.Power, scn.Match.GetPower(stranger, false), "only its own races")
		assert.Equal(t, theirKin.Power, scn.Match.GetPower(theirKin, false), "only its controller's creatures")
		assert.Equal(t, hydrooze.Power, scn.Match.GetPower(hydrooze, false), "\"other\" spares itself")
	})

	t.Run("its kin cannot be blocked, itself included", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hydrooze := putCardInBattlezone(t, scn, player.Player, hydroozeUID, hydroozeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, hydroozeKinUID, hydroozeSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, hydroozeStrangerUID, hydroozeSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, hydroozeKinUID, hydroozeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.True(t, kin.HasCondition(cnd.CantBeBlocked))
		assert.True(t, hydrooze.HasCondition(cnd.CantBeBlocked), "the clause has no \"other\"")
		assert.False(t, stranger.HasCondition(cnd.CantBeBlocked))
		assert.False(t, theirKin.HasCondition(cnd.CantBeBlocked))
	})

	t.Run("an unblockable kin walks past a blocker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, hydroozeUID, hydroozeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, hydroozeKinUID, hydroozeSetupSrc)
		wall := putCardInBattlezone(t, scn, opponent.Player, hydroozeWallUID, hydroozeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, wall.HasCondition(cnd.Blocker), "the wall has to be able to block for this to mean anything")

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		_, err = scn.ActionAttackPlayer(player, kin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		// No block prompt is ever opened, so the loop settles on its own.
		settleTurn(t, scn)

		shieldsAfter, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		assert.Len(t, shieldsAfter, shieldCount-1)
		assert.Equal(t, match.BATTLEZONE, wall.Zone, "it never got to block")
	})

	t.Run("both grants leave with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hydrooze := putCardInBattlezone(t, scn, player.Player, hydroozeUID, hydroozeSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, hydroozeKinUID, hydroozeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		require.True(t, kin.HasCondition(cnd.CantBeBlocked))

		_, err := player.Player.MoveCard(hydrooze.ID, match.BATTLEZONE, match.GRAVEYARD, hydroozeSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power, scn.Match.GetPower(kin, false))
		assert.False(t, kin.HasCondition(cnd.CantBeBlocked))
	})
}
