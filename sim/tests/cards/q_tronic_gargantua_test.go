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
	qTronicGargantuaUID       = "adc7eb5d-862d-430c-bd7d-ca51f0c94b02"
	qTronicGargantuaSurvivorQ = "d176b30a-cac6-4249-a78d-18f34b97546b" // Promephius Q
	qTronicGargantuaSetupSrc  = "q_tronic_gargantua_test_setup"
)

func TestQTronicGargantua(t *testing.T) {
	t.Run("breaks one more shield for each other survivor", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		gargantua := putCardInBattlezone(t, scn, player.Player, qTronicGargantuaUID, qTronicGargantuaSetupSrc)

		assert.Equal(t, "Q-tronic Gargantua", gargantua.Name)
		assert.Equal(t, 9000, gargantua.Power)
		assert.Equal(t, []string{civ.Fire}, gargantua.Civs)
		assert.True(t, gargantua.HasFamily(family.Survivor))
		assert.Equal(t, 0, shieldBreakModifierOf(t, gargantua), "no other survivors yet")

		putCardInBattlezone(t, scn, player.Player, qTronicGargantuaSurvivorQ, qTronicGargantuaSetupSrc)
		assert.Equal(t, 1, shieldBreakModifierOf(t, gargantua))

		second := putCardInBattlezone(t, scn, player.Player, qTronicGargantuaSurvivorQ, qTronicGargantuaSetupSrc)
		assert.Equal(t, 2, shieldBreakModifierOf(t, gargantua))

		moved, err := player.Player.MoveCard(second.ID, match.BATTLEZONE, match.GRAVEYARD, qTronicGargantuaSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		assert.Equal(t, 1, shieldBreakModifierOf(t, gargantua), "the modifier follows the battle zone")
	})

	t.Run("keeps its own survivor and evolution conditions when attacking", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		gargantua := putCardInBattlezone(t, scn, player.Player, qTronicGargantuaUID, qTronicGargantuaSetupSrc)
		putCardInBattlezone(t, scn, player.Player, qTronicGargantuaSurvivorQ, qTronicGargantuaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, gargantua.HasCondition(cnd.Survivor))
		require.True(t, gargantua.HasCondition(cnd.Evolution))

		// 1 base shield plus 1 for the other survivor.
		action, err := scn.ActionAttackPlayer(player, gargantua.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))

		assert.True(t, gargantua.HasCondition(cnd.Survivor), "attacking must not strip its own survivor condition")
		assert.True(t, gargantua.HasCondition(cnd.Evolution), "nor its evolution condition")
		assert.Equal(t, 1, shieldBreakModifierOf(t, gargantua), "and the modifier stays accurate")

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 3)
	})

	t.Run("a cancelled attack leaves no stale modifier", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		gargantua := putCardInBattlezone(t, scn, player.Player, qTronicGargantuaUID, qTronicGargantuaSetupSrc)
		other := putCardInBattlezone(t, scn, player.Player, qTronicGargantuaSurvivorQ, qTronicGargantuaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err := scn.ActionAttackPlayer(player, gargantua.ID)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		// The other survivor leaves, so a modifier held over from the cancelled
		// attack would now be wrong.
		moved, err := player.Player.MoveCard(other.ID, match.BATTLEZONE, match.GRAVEYARD, qTronicGargantuaSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.Equal(t, 0, shieldBreakModifierOf(t, gargantua))
		assert.True(t, gargantua.HasCondition(cnd.Survivor))
	})
}

func shieldBreakModifierOf(t *testing.T, card *match.Card) int {
	t.Helper()

	total := 0
	for _, condition := range card.Conditions() {
		if condition.ID != cnd.ShieldBreakModifier {
			continue
		}

		val, ok := condition.Val.(int)
		require.True(t, ok)
		total += val
	}

	return total
}
