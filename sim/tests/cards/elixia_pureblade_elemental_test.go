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
	elixiaPurebladeElementalUID      = "b9380b9c-894f-421b-b060-1450a4eb32cd"
	elixiaPurebladeElementalSetupSrc = "elixia_pureblade_elemental_test_setup"
)

// One vanilla creature per civilization, so the mana zone's civilization count
// can be controlled exactly.
var elixiaPurebladeElementalManaUIDs = map[string]string{
	civ.Light:    "7b58e8c2-0b1e-4ef5-812f-e667c2092c73", // Reusol, the Oracle
	civ.Water:    "9781089f-1aa9-4a75-b106-35e9d431e31d", // Aqua Vehicle
	civ.Darkness: "e2b992ee-91a3-49d3-8228-7be60a0b9ec5", // Writhing Bone Ghoul
	civ.Fire:     "af3bc221-1cc2-4f58-83ea-2673ac2c66c5", // Immortal Baron, Vorg
	civ.Nature:   "1d72eb3e-5185-449a-a16f-391bd2338343", // Burning Mane
}

func TestElixiaPurebladeElemental(t *testing.T) {
	t.Run("the mana bonus is inactive outside the battle zone", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		player.Player.SpawnCard(elixiaPurebladeElementalUID, match.HAND)
		elixia, err := scn.FindCard(player.Player, match.HAND, elixiaPurebladeElementalUID)
		require.NoError(t, err)
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Fire)

		assert.Equal(t, "Elixia, Pureblade Elemental", elixia.Name)
		assert.Equal(t, 1000, elixia.Power)
		assert.Equal(t, 6, elixia.ManaCost)
		assert.Equal(t, []string{civ.Light}, elixia.Civs)
		assert.True(t, elixia.HasFamily(family.AngelCommand))

		assert.Equal(t, 1000, scn.Match.GetPower(elixia, false))
		assert.False(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))
	})

	t.Run("counts each civilization in the mana zone once", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		elixia := putElixiaPurebladeElementalTestCardInBattlezone(t, scn, player.Player)

		assert.Equal(t, 1000, scn.Match.GetPower(elixia, false), "an empty mana zone grants nothing")

		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light)
		assert.Equal(t, 4000, scn.Match.GetPower(elixia, false))

		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Light, civ.Light)
		assert.Equal(t, 4000, scn.Match.GetPower(elixia, false), "duplicate civilizations only count once")

		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water)
		assert.Equal(t, 7000, scn.Match.GetPower(elixia, false))
		assert.Equal(t, 1000, elixia.Power, "dynamic power must not mutate printed power")

		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Darkness, civ.Fire, civ.Nature)
		assert.Equal(t, 16000, scn.Match.GetPower(elixia, false))
	})

	t.Run("only its controller's mana zone counts", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
		elixia := putElixiaPurebladeElementalTestCardInBattlezone(t, scn, player.Player)

		setElixiaPurebladeElementalTestMana(t, scn, opponent, civ.Light, civ.Water, civ.Darkness)
		assert.Equal(t, 1000, scn.Match.GetPower(elixia, false))

		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Fire)
		assert.Equal(t, 4000, scn.Match.GetPower(elixia, false))
	})

	t.Run("switches between no breaker, double breaker and triple breaker", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		elixia := putElixiaPurebladeElementalTestCardInBattlezone(t, scn, player.Player)

		// 1 civilization -> 4000 power
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light)
		assert.False(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))

		// 2 civilizations -> 7000 power
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water)
		assert.True(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))

		// 4 civilizations -> 13000 power, still short of the triple breaker tier
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Darkness, civ.Fire)
		assert.True(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))

		// 5 civilizations -> 16000 power
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Darkness, civ.Fire, civ.Nature)
		assert.True(t, elixia.HasCondition(cnd.TripleBreaker))
		assert.False(t, elixia.HasCondition(cnd.DoubleBreaker), "triple breaker replaces double breaker")

		// Back down to a single civilization
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light)
		assert.False(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))
	})

	t.Run("breaks three shields while it has 15000 power or more", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
		elixia := putElixiaPurebladeElementalTestCardInBattlezone(t, scn, player.Player)
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Darkness, civ.Fire, civ.Nature)
		require.Equal(t, 16000, scn.Match.GetPower(elixia, false))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, elixia.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, action.MinSelections)
		assert.Equal(t, 3, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-3)
	})

	t.Run("loses its breakers when it leaves the battle zone", func(t *testing.T) {
		scn, player := setupElixiaPurebladeElementalTest(t)
		elixia := putElixiaPurebladeElementalTestCardInBattlezone(t, scn, player.Player)
		setElixiaPurebladeElementalTestMana(t, scn, player, civ.Light, civ.Water, civ.Darkness, civ.Fire, civ.Nature)
		require.True(t, elixia.HasCondition(cnd.TripleBreaker))

		moved, err := player.Player.MoveCard(elixia.ID, match.BATTLEZONE, match.GRAVEYARD, elixiaPurebladeElementalSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, elixia.HasCondition(cnd.DoubleBreaker))
		assert.False(t, elixia.HasCondition(cnd.TripleBreaker))
		assert.Equal(t, 1000, scn.Match.GetPower(elixia, false))
	})
}

func setupElixiaPurebladeElementalTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	return scn, player
}

func putElixiaPurebladeElementalTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player) *match.Card {
	t.Helper()

	player.SpawnCard(elixiaPurebladeElementalUID, match.HAND)
	card, err := scn.FindCard(player, match.HAND, elixiaPurebladeElementalUID)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, elixiaPurebladeElementalSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

// setElixiaPurebladeElementalTestMana rebuilds a player's mana zone so it holds
// exactly one card of each listed civilization, moving cards through the engine
// so the continuous effect is re-evaluated.
func setElixiaPurebladeElementalTestMana(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, civilizations ...string) {
	t.Helper()

	mana, err := player.Player.Container(match.MANAZONE)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), mana...) {
		moved, err := player.Player.MoveCard(card.ID, match.MANAZONE, match.GRAVEYARD, elixiaPurebladeElementalSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}

	for _, civilization := range civilizations {
		uid, ok := elixiaPurebladeElementalManaUIDs[civilization]
		require.True(t, ok, "no test card registered for civilization %s", civilization)

		player.Player.SpawnCard(uid, match.HAND)
		card, err := scn.FindCard(player.Player, match.HAND, uid)
		require.NoError(t, err)
		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, elixiaPurebladeElementalSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.MANAZONE, moved.Zone)
	}
}
