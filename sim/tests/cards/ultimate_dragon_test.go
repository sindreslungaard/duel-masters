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
	ultimateDragonUID          = "31683921-16fc-4d4c-bb77-3225d10f7366"
	ultimateDragonDragonUID    = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (Volcano Dragon)
	ultimateDragonDragonoidUID = "c125f786-e6d5-4477-8ab0-1e92d6eed348" // Cavalry General Curatops (Dragonoid)
	ultimateDragonNonDragonUID = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (Rock Beast)
	ultimateDragonSetupSrc     = "ultimate_dragon_test_setup"
)

func TestUltimateDragon(t *testing.T) {
	t.Run("printed characteristics and no bonus while alone", func(t *testing.T) {
		scn, player, _ := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)

		assert.Equal(t, "Ultimate Dragon", ultimate.Name)
		assert.Equal(t, 5000, ultimate.Power)
		assert.Equal(t, 6, ultimate.ManaCost)
		assert.Equal(t, civ.Fire, ultimate.Civ)
		assert.True(t, ultimate.HasFamily(family.ArmoredDragon))

		assert.Equal(t, 5000, scn.Match.GetPower(ultimate, false))
		assert.False(t, ultimate.HasCondition(cnd.ShieldBreakModifier))
	})

	t.Run("gains power and a shield for each other own dragon", func(t *testing.T) {
		scn, player, _ := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)

		firstDragon := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonDragonUID)
		assert.Equal(t, 10000, scn.Match.GetPower(ultimate, false))
		assert.Equal(t, 1, ultimateDragonTestShieldBreakModifier(t, ultimate))

		// A second Ultimate Dragon is an Armored Dragon, so each counts the other.
		secondUltimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)
		assert.Equal(t, 15000, scn.Match.GetPower(ultimate, false))
		assert.Equal(t, 2, ultimateDragonTestShieldBreakModifier(t, ultimate))
		assert.Equal(t, 15000, scn.Match.GetPower(secondUltimate, false))
		assert.Equal(t, 2, ultimateDragonTestShieldBreakModifier(t, secondUltimate))

		moved, err := player.Player.MoveCard(firstDragon.ID, match.BATTLEZONE, match.GRAVEYARD, ultimateDragonSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		assert.Equal(t, 10000, scn.Match.GetPower(ultimate, false), "the bonus shrinks with the battle zone")
		assert.Equal(t, 1, ultimateDragonTestShieldBreakModifier(t, ultimate))
	})

	t.Run("ignores non-dragons, dragonoids and the opponent's dragons", func(t *testing.T) {
		scn, player, opponent := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)

		putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonNonDragonUID)
		assert.Equal(t, 5000, scn.Match.GetPower(ultimate, false), "a Rock Beast is not a Dragon")

		dragonoid := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonDragonoidUID)
		require.True(t, dragonoid.HasFamily(family.Dragonoid))
		assert.Equal(t, 5000, scn.Match.GetPower(ultimate, false), "Dragonoid is not a Dragon race")

		putUltimateDragonTestCardInBattlezone(t, scn, opponent.Player, ultimateDragonDragonUID)
		assert.Equal(t, 5000, scn.Match.GetPower(ultimate, false), "only your own creatures count")
		assert.False(t, ultimate.HasCondition(cnd.ShieldBreakModifier))
	})

	t.Run("breaks one more shield for each other own dragon", func(t *testing.T) {
		scn, player, opponent := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)
		putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonDragonUID)
		putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonDragonUID)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, ultimate.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, action.MinSelections, "1 base shield plus 1 per other dragon")
		assert.Equal(t, 3, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-3)
	})

	t.Run("breaks a single shield with no other dragons", func(t *testing.T) {
		scn, player, opponent := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)

		action, err := scn.ActionAttackPlayer(player, ultimate.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 4)
	})

	t.Run("drops the crew breaker when it leaves the battle zone", func(t *testing.T) {
		scn, player, _ := setupUltimateDragonTest(t)
		ultimate := putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonUID)
		putUltimateDragonTestCardInBattlezone(t, scn, player.Player, ultimateDragonDragonUID)
		require.Equal(t, 1, ultimateDragonTestShieldBreakModifier(t, ultimate))

		moved, err := player.Player.MoveCard(ultimate.ID, match.BATTLEZONE, match.GRAVEYARD, ultimateDragonSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, ultimate.HasCondition(cnd.ShieldBreakModifier))
		assert.Equal(t, 5000, scn.Match.GetPower(ultimate, false))
	})
}

func setupUltimateDragonTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

func putUltimateDragonTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, ultimateDragonSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

func ultimateDragonTestShieldBreakModifier(t *testing.T, card *match.Card) int {
	t.Helper()

	total := 0
	for _, condition := range card.Conditions() {
		if condition.ID != cnd.ShieldBreakModifier {
			continue
		}

		val, ok := condition.Val.(int)
		require.True(t, ok, "shield break modifier must carry an int value")
		total += val
	}

	return total
}
