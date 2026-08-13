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
	deklowazTheTerminatorUID      = "a808b98c-2de7-412b-970c-a3b925bf43c2"
	deklowazSmallCreatureUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	deklowazEdgeCreatureUID       = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	deklowazBigCreatureUID        = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	deklowazSpellUID              = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	deklowazTheTerminatorSetupSrc = "deklowaz_the_terminator_test_setup"
)

func TestDeklowazTheTerminator(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Deklowaz, the Terminator", deklowaz.Name)
		assert.Equal(t, 5000, deklowaz.Power)
		assert.Equal(t, 6, deklowaz.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, deklowaz.Civs)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, deklowaz.ManaRequirement)
		assert.True(t, deklowaz.IsMulticolored())
		assert.True(t, deklowaz.HasFamily(family.SpiritQuartz))
		assert.True(t, deklowaz.HasCondition(cnd.TapAbility))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(deklowazTheTerminatorUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, deklowazTheTerminatorSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("the tap ability sweeps the battle zone of everything at 3000 or less", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)

		ownSmall := putCardInBattlezone(t, scn, player.Player, deklowazSmallCreatureUID, deklowazTheTerminatorSetupSrc)
		theirEdge := putCardInBattlezone(t, scn, opponent.Player, deklowazEdgeCreatureUID, deklowazTheTerminatorSetupSrc)
		theirBig := putCardInBattlezone(t, scn, opponent.Player, deklowazBigCreatureUID, deklowazTheTerminatorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionUseTapAbility(player, deklowaz.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, ownSmall.Zone, "its own side is swept too")
		assert.Equal(t, match.GRAVEYARD, theirEdge.Zone, "exactly 3000 is within range")
		assert.Equal(t, match.BATTLEZONE, theirBig.Zone)
		assert.Equal(t, match.BATTLEZONE, deklowaz.Zone, "at 5000 power it survives its own sweep")
		assert.True(t, deklowaz.Tapped, "using the tap ability taps it")
	})

	t.Run("the opponent discards the small creatures in their hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)

		small, edge, big, spell := deklowazSeedHand(t, opponent)
		passTurnToSelf(t, scn, player, opponent)

		require.True(t, small.HasCondition(cnd.Creature))
		require.True(t, spell.HasCondition(cnd.Spell))

		require.NoError(t, scn.ActionUseTapAbility(player, deklowaz.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, small.Zone)
		assert.Equal(t, match.GRAVEYARD, edge.Zone, "exactly 3000 is within range")
		assert.Equal(t, match.HAND, big.Zone)
		assert.Equal(t, match.HAND, spell.Zone, "only creatures are discarded")
	})

	t.Run("it does not discard from its controller's own hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)

		ownSmall, _, _, _ := deklowazSeedHand(t, player)
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionUseTapAbility(player, deklowaz.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, ownSmall.Zone)
	})

	t.Run("an empty opposing hand is not a problem", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		for _, card := range hand {
			_, err := opponent.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, deklowazTheTerminatorSetupSrc)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionUseTapAbility(player, deklowaz.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, deklowaz.Tapped)
		assert.Equal(t, match.BATTLEZONE, deklowaz.Zone)
	})

	t.Run("a tapped creature cannot use the ability again", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		deklowaz := putCardInBattlezone(t, scn, player.Player, deklowazTheTerminatorUID, deklowazTheTerminatorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionUseTapAbility(player, deklowaz.ID))
		require.NoError(t, scn.WaitForEventLoop())
		require.True(t, deklowaz.Tapped)

		theirSmall := putCardInBattlezone(t, scn, opponent.Player, deklowazSmallCreatureUID, deklowazTheTerminatorSetupSrc)

		require.Error(t, scn.ActionUseTapAbility(player, deklowaz.ID), "a tapped creature has nothing left to tap")

		assert.Equal(t, match.BATTLEZONE, theirSmall.Zone, "the sweep did not run a second time")
	})
}

// deklowazSeedHand fills a hand with one creature under, one on and one over
// the 3000 line, plus a spell. It is called before the turn handover so an
// untap step gives each card its identity, the way a card that came out of the
// deck would already have it.
func deklowazSeedHand(t *testing.T, player *match.PlayerReference) (*match.Card, *match.Card, *match.Card, *match.Card) {
	t.Helper()

	small, err := player.Player.SpawnCard(deklowazSmallCreatureUID, match.HAND)
	require.NoError(t, err)
	edge, err := player.Player.SpawnCard(deklowazEdgeCreatureUID, match.HAND)
	require.NoError(t, err)
	big, err := player.Player.SpawnCard(deklowazBigCreatureUID, match.HAND)
	require.NoError(t, err)
	spell, err := player.Player.SpawnCard(deklowazSpellUID, match.HAND)
	require.NoError(t, err)

	return small, edge, big, spell
}
