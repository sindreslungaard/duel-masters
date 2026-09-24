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
	bolmeteusSteelDragonUID      = "4717faef-1065-4153-a509-854c22637e27"
	bolmeteusSteelDragonSetupSrc = "bolmeteus_steel_dragon_test_setup"
)

func TestBolmeteusSteelDragon(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, bolmeteusSteelDragonUID, bolmeteusSteelDragonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, dragon, "Bolmeteus Steel Dragon", 7000, 7, []string{civ.Fire})
		assert.True(t, dragon.HasFamily(family.ArmoredDragon))
		assert.True(t, dragon.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("shields it would break go to the opponent's graveyard instead of hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, bolmeteusSteelDragonUID, bolmeteusSteelDragonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		action, err := scn.ActionAttackPlayer(player, dragon.ID)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(action.Cards), 2)
		targeted := []string{action.Cards[0].CardID, action.Cards[1].CardID}

		require.NoError(t, scn.ResolveAttack(player, targeted...))
		settleTurn(t, scn)

		for _, id := range targeted {
			card, err := opponent.Player.GetCard(id, match.GRAVEYARD)
			require.NoError(t, err, "the shield landed in the graveyard, not the hand")
			assert.Equal(t, match.GRAVEYARD, card.Zone)
		}

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount, "the broken shields never reached hand")
	})

	t.Run("only applies to shields it breaks itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, bolmeteusSteelDragonUID, bolmeteusSteelDragonSetupSrc)
		dragon.Tapped = true // present, but not the attacker
		attacker := putCardInBattlezone(t, scn, player.Player, cruelNagaFillerUID, bolmeteusSteelDragonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		target := shields[0]

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, target.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, target.Zone, "an unrelated attacker's broken shield still goes to hand")
	})
}
