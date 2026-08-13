package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tanzanyteTheAwakenerUID      = "4a223c53-e0fc-4d85-b531-c325c94ed188"
	tanzanyteCreatureUID         = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	tanzanyteOtherCreatureUID    = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur
	tanzanyteSpellUID            = "b7f236fd-e7eb-41cc-912a-5239c134f265" // Energy Stream
	tanzanyteTheAwakenerSetupSrc = "tanzanyte_the_awakener_test_setup"
)

func TestTanzanyteTheAwakener(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupTanzanyteTest(t)
		tanzanyte := putCardInBattlezone(t, scn, player.Player, tanzanyteTheAwakenerUID, tanzanyteTheAwakenerSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Tanzanyte, the Awakener", tanzanyte.Name)
		assert.Equal(t, 9000, tanzanyte.Power)
		assert.Equal(t, 7, tanzanyte.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, tanzanyte.Civs)
		assert.True(t, tanzanyte.HasFamily(family.SpiritQuartz))
		assert.True(t, tanzanyte.HasCondition(cnd.DoubleBreaker))
		assert.True(t, tanzanyte.HasCondition(cnd.TapAbility))
	})

	t.Run("returns every copy of the chosen creature from the graveyard", func(t *testing.T) {
		scn, player, opponent := setupTanzanyteTest(t)
		tanzanyte := putCardInBattlezone(t, scn, player.Player, tanzanyteTheAwakenerUID, tanzanyteTheAwakenerSetupSrc)

		copies := spawnInGraveyard(t, scn, player, tanzanyteCreatureUID, 3)
		other := spawnInGraveyard(t, scn, player, tanzanyteOtherCreatureUID, 1)[0]
		spell := spawnInGraveyard(t, scn, player, tanzanyteSpellUID, 1)[0]

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionUseTapAbility(player, tanzanyte.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.False(t, action.Cancellable, "choosing is mandatory")
		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		assert.Contains(t, offered, copies[0].ID)
		assert.Contains(t, offered, other.ID)
		assert.NotContains(t, offered, spell.ID, "only creatures may be chosen")

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, copies[0].ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		for _, copy := range copies {
			assert.Equal(t, match.HAND, copy.Zone, "every creature with that name returns")
		}
		assert.Equal(t, match.GRAVEYARD, other.Zone, "a different name stays behind")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
		assert.True(t, tanzanyte.Tapped)
	})

	t.Run("resolves without a prompt when the graveyard holds no creature", func(t *testing.T) {
		scn, player, opponent := setupTanzanyteTest(t)
		tanzanyte := putCardInBattlezone(t, scn, player.Player, tanzanyteTheAwakenerUID, tanzanyteTheAwakenerSetupSrc)
		spell := spawnInGraveyard(t, scn, player, tanzanyteSpellUID, 1)[0]

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionUseTapAbility(player, tanzanyte.ID))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
		assert.True(t, tanzanyte.Tapped)
	})
}

func setupTanzanyteTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

// spawnInGraveyard puts count copies of a card straight into a player's
// graveyard and returns them.
func spawnInGraveyard(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string, count int) []*match.Card {
	t.Helper()

	spawned := make([]*match.Card, 0, count)
	for range count {
		card, err := player.Player.SpawnCard(uid, match.GRAVEYARD)
		require.NoError(t, err)
		spawned = append(spawned, card)
	}

	return spawned
}
