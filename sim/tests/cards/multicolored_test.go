package cards

import (
	"duel-masters/game/cards"
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Gonta, the Warrior Savage is a fire/nature creature with no printed
	// ability beyond the multicolored rules themselves.
	gontaTheWarriorSavageUID = "dfe767f5-8883-4d3c-80ee-df3b277ff425"
	// Melnia, the Aqua Shadow is water/darkness.
	melniaTheAquaShadowUID = "ddccdc18-92ef-431e-913e-71ba5bb6b1b1"

	multicolorFireUID     = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	multicolorNatureUID   = "1d72eb3e-5185-449a-a16f-391bd2338343" // Burning Mane
	multicolorWaterUID    = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	multicolorDarknessUID = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul
	multicolorSetupSrc    = "multicolored_test_setup"
)

func TestMulticoloredCards(t *testing.T) {
	t.Run("printed characteristics carry every civilization", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		assert.Equal(t, "Gonta, the Warrior Savage", gonta.Name)
		assert.Equal(t, 4000, gonta.Power)
		assert.Equal(t, 2, gonta.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, gonta.Civs)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, gonta.ManaRequirement)
		assert.True(t, gonta.HasFamily(family.Human))
		assert.True(t, gonta.HasFamily(family.BeastFolk))

		assert.True(t, gonta.IsMulticolored())
		assert.True(t, gonta.HasCiv(civ.Fire), "a multicolored card is a card of each of its civilizations")
		assert.True(t, gonta.HasCiv(civ.Nature))
		assert.False(t, gonta.HasCiv(civ.Water))
	})

	t.Run("a mono coloured card is not multicoloured", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(multicolorFireUID, match.HAND)
		vorg, err := scn.FindCard(player.Player, match.HAND, multicolorFireUID)
		require.NoError(t, err)

		assert.False(t, vorg.IsMulticolored())
		assert.True(t, vorg.HasCiv(civ.Fire))
		assert.False(t, vorg.HasCiv(civ.Nature))
	})

	t.Run("needs one mana of each civilization, not just one of them", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		fireMana := spawnMulticolorTestMana(t, scn, player, multicolorFireUID, 2)
		assert.False(t, player.Player.CanPlayCard(gonta, fireMana), "two fire mana cannot pay a fire/nature cost")

		natureMana := spawnMulticolorTestMana(t, scn, player, multicolorNatureUID, 2)
		assert.False(t, player.Player.CanPlayCard(gonta, natureMana), "two nature mana cannot pay it either")

		mixed := []*match.Card{fireMana[0], natureMana[0]}
		assert.True(t, player.Player.CanPlayCard(gonta, mixed), "one of each pays the cost")
	})

	t.Run("mono coloured cards still need only their own civilization", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(multicolorFireUID, match.HAND)
		vorg, err := scn.FindCard(player.Player, match.HAND, multicolorFireUID)
		require.NoError(t, err)

		fireMana := spawnMulticolorTestMana(t, scn, player, multicolorFireUID, 2)
		assert.True(t, player.Player.CanPlayCard(vorg, fireMana))

		natureMana := spawnMulticolorTestMana(t, scn, player, multicolorNatureUID, 2)
		assert.False(t, player.Player.CanPlayCard(vorg, natureMana))
	})

	t.Run("a single multicolored mana card cannot pay for two of its civilizations at once", func(t *testing.T) {
		// Official ruling: when tapped for mana, a multicolored card is used as
		// only one of its civilizations, never several simultaneously. A report
		// surfaced a live bug where Gonta, a fire/nature 2-drop, could be
		// summoned by tapping a single fire/nature card (Wind Axe) plus any
		// other unrelated card, because the old civilization check only asked
		// "does some selected card have this civilization" per requirement
		// without tracking which card had already been claimed.
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		multicolorMana := spawnMulticolorTestMana(t, scn, player, gontaTheWarriorSavageUID, 1)
		multicolorMana[0].Tapped = false
		anyMana := spawnMulticolorTestMana(t, scn, player, multicolorWaterUID, 1)

		assert.False(t, player.Player.CanPlayCard(gonta, append(multicolorMana, anyMana...)),
			"the fire/nature mana card can only stand in for one requirement, and the water card fulfills neither")
	})

	t.Run("two separate multicolored mana cards may each cover a different civilization", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		multicolorMana := spawnMulticolorTestMana(t, scn, player, gontaTheWarriorSavageUID, 2)
		for _, mana := range multicolorMana {
			mana.Tapped = false
		}

		assert.True(t, player.Player.CanPlayCard(gonta, multicolorMana),
			"one fire/nature card can cover the fire requirement while a second, distinct fire/nature card covers nature")
	})

	t.Run("is put into the mana zone tapped, however it gets there", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(gonta.ID, match.HAND, match.MANAZONE, multicolorSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.MANAZONE, moved.Zone)
		assert.True(t, moved.Tapped, "multicolored cards arrive in the mana zone tapped")

		// A mono coloured card arriving the same way is untapped.
		player.Player.SpawnCard(multicolorFireUID, match.HAND)
		vorg, err := scn.FindCard(player.Player, match.HAND, multicolorFireUID)
		require.NoError(t, err)
		movedVorg, err := player.Player.MoveCard(vorg.ID, match.HAND, match.MANAZONE, multicolorSetupSrc)
		require.NoError(t, err)
		assert.False(t, movedVorg.Tapped)

		// It untaps normally on its controller's next untap step.
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.False(t, gonta.Tapped, "the tapped state only applies on arrival")
	})

	t.Run("moving to any other zone still untaps", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(gonta.ID, match.HAND, match.BATTLEZONE, multicolorSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, moved.Zone)
		assert.False(t, moved.Tapped)
	})

	t.Run("can be summoned through the real play flow", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)

		spawnMulticolorTestMana(t, scn, player, multicolorFireUID, 2)
		spawnMulticolorTestMana(t, scn, player, multicolorNatureUID, 2)

		require.NoError(t, scn.ActionPlayCard(player, gonta.ID))
		assert.Equal(t, match.BATTLEZONE, gonta.Zone)

		// The mana selection helper must have chosen at least one of each.
		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		tappedFire, tappedNature := 0, 0
		for _, manaCard := range mana {
			if !manaCard.Tapped {
				continue
			}
			if manaCard.HasCiv(civ.Fire) {
				tappedFire++
			}
			if manaCard.HasCiv(civ.Nature) {
				tappedNature++
			}
		}
		assert.Equal(t, 1, tappedFire)
		assert.Equal(t, 1, tappedNature)
	})

	t.Run("counts as each of its civilizations for other cards", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		// Hourglass Mutant grants slayer to its controller's water and fire
		// creatures. A water/darkness creature is a water creature.
		putMulticolorTestCardInBattlezone(t, scn, player.Player, hourglassMutantUID)
		melnia := putMulticolorTestCardInBattlezone(t, scn, player.Player, melniaTheAquaShadowUID)
		gonta := putMulticolorTestCardInBattlezone(t, scn, player.Player, gontaTheWarriorSavageUID)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.True(t, melnia.HasCondition(cnd.Slayer), "water/darkness counts as water")
		assert.True(t, gonta.HasCondition(cnd.Slayer), "fire/nature counts as fire")
	})

	t.Run("Melnia has its printed abilities", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		melnia := putMulticolorTestCardInBattlezone(t, scn, player.Player, melniaTheAquaShadowUID)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Melnia, the Aqua Shadow", melnia.Name)
		assert.Equal(t, 1000, melnia.Power)
		assert.Equal(t, 2, melnia.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, melnia.Civs)
		assert.True(t, melnia.HasFamily(family.LiquidPeople))
		assert.True(t, melnia.HasFamily(family.Ghost))
		assert.True(t, melnia.HasCondition(cnd.CantBeBlocked))
		assert.True(t, melnia.HasCondition(cnd.Slayer))
	})

	t.Run("cost reduction on a multicolored card floors at its number of civilizations, not 1", func(t *testing.T) {
		// Gonta is fire/nature and must be paid with one distinct mana card of
		// each civilization. A cost reduction that floored at a flat 1 would ask
		// the player to select a single mana card, which can never satisfy two
		// civilizations, making the card unplayable even with enough mana.
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(gontaTheWarriorSavageUID, match.HAND)
		gonta, err := scn.FindCard(player.Player, match.HAND, gontaTheWarriorSavageUID)
		require.NoError(t, err)
		require.Equal(t, 2, gonta.ManaCost)

		gonta.AddUniqueSourceCondition(cnd.ReducedCost, 5, multicolorSetupSrc)
		assert.Equal(t, 2, gonta.EffectiveManaCost(), "cost cannot drop below one mana card per required civilization")

		fireMana := spawnMulticolorTestMana(t, scn, player, multicolorFireUID, 1)
		natureMana := spawnMulticolorTestMana(t, scn, player, multicolorNatureUID, 1)
		fireMana[0].Tapped = false
		natureMana[0].Tapped = false

		require.NoError(t, scn.ActionPlayCard(player, gonta.ID))
		assert.Equal(t, match.BATTLEZONE, gonta.Zone)
		assert.True(t, fireMana[0].Tapped)
		assert.True(t, natureMana[0].Tapped)
	})

	t.Run("cost reduction on a mono coloured card still floors at 1", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()

		player.Player.SpawnCard(multicolorFireUID, match.HAND)
		vorg, err := scn.FindCard(player.Player, match.HAND, multicolorFireUID)
		require.NoError(t, err)

		vorg.AddUniqueSourceCondition(cnd.ReducedCost, 99, multicolorSetupSrc)
		assert.Equal(t, 1, vorg.EffectiveManaCost())
	})
}

// TestEveryRegisteredCardMatchesItsCatalogCivilizations walks the whole card
// registry rather than a hand picked sample, so a card registered with the wrong
// civilizations, or a multicolored card that forgot fx.PutIntoManaZoneTapped,
// fails here without needing its own test.
func TestEveryRegisteredCardMatchesItsCatalogCivilizations(t *testing.T) {
	catalog := readCatalogCivilizations(t)

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	multicoloredUIDs := make([]string, 0)
	checked := 0

	for setID, set := range cards.Sets {
		for uid, constructor := range *set {
			if constructor == nil {
				continue
			}

			expected, ok := catalog[uid]
			require.True(t, ok, "%s card %s is registered but missing from the card catalog", setID, uid)

			// Spawning into the hand dispatches no events, so the whole registry
			// can be classified cheaply in one match.
			card, err := player.Player.SpawnCard(uid, match.HAND)
			require.NoError(t, err)

			assert.Equal(t, len(expected), len(card.Civs), "%s (%s) civilization count", card.Name, uid)
			for _, civilization := range expected {
				assert.True(t, card.HasCiv(civilization), "%s (%s) should be a %s card", card.Name, uid, civilization)
			}
			assert.Equal(t, len(expected) > 1, card.IsMulticolored(), "%s (%s) multicolored", card.Name, uid)

			checked++
			if len(expected) > 1 {
				multicoloredUIDs = append(multicoloredUIDs, uid)
			}
		}
	}

	t.Logf("checked %d registered cards, %d of them multicolored", checked, len(multicoloredUIDs))
	require.NotEmpty(t, multicoloredUIDs, "the sweep must actually cover multicolored cards")

	// Each card is moved in its own match so that no other card can react to the
	// move, and through MoveCard so the CardMoved event really is what taps it.
	for _, uid := range multicoloredUIDs {
		card := moveIntoManaZone(t, uid)
		assert.True(t, card.Tapped, "%s (%s) must be put into the mana zone tapped", card.Name, uid)
	}

	mono := moveIntoManaZone(t, multicolorFireUID)
	assert.False(t, mono.Tapped, "a mono coloured card is not tapped on arrival")
}

// moveIntoManaZone charges a single card into an otherwise empty mana zone.
func moveIntoManaZone(t *testing.T, uid string) *match.Card {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	spawned, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	moved, err := player.Player.MoveCard(spawned.ID, match.HAND, match.MANAZONE, multicolorSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.MANAZONE, moved.Zone)

	return moved
}

// readCatalogCivilizations maps every catalog card id to its civilizations,
// lowercased to match the civ package constants.
func readCatalogCivilizations(t *testing.T) map[string][]string {
	t.Helper()

	data, err := os.ReadFile("../../DuelMastersCards.json")
	require.NoError(t, err)

	var parsed struct {
		Cards []struct {
			ID            string   `json:"id"`
			Civilizations []string `json:"civilizations"`
		} `json:"cards"`
	}
	require.NoError(t, json.Unmarshal(data, &parsed))

	civilizations := make(map[string][]string, len(parsed.Cards))
	for _, card := range parsed.Cards {
		if card.ID == "" {
			continue
		}

		lowered := make([]string, 0, len(card.Civilizations))
		for _, civilization := range card.Civilizations {
			lowered = append(lowered, strings.ToLower(civilization))
		}
		civilizations[card.ID] = lowered
	}

	return civilizations
}

func spawnMulticolorTestMana(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string, count int) []*match.Card {
	t.Helper()

	spawned := make([]*match.Card, 0, count)
	for range count {
		card, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
		spawned = append(spawned, card)
	}

	return spawned
}

func putMulticolorTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, multicolorSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
