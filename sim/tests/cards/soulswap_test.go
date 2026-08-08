package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const soulswapUID = "a82ec211-588c-4308-95be-798581045e31"

func TestSoulswap(t *testing.T) {
	t.Run("swaps an opponent's creature using the caster's choices", func(t *testing.T) {
		scn, caster, opponent, spell := setupSoulswapTest(t)
		selectedCreature := putSoulswapTestCardInZone(t, scn, opponent.Player, scowlingTomatoUID, match.BATTLEZONE)
		selectedCreature.AddCondition(cnd.Creature, nil, nil)
		manaCreature := putSoulswapTestCardInZone(t, scn, opponent.Player, shamanBroccoliUID, match.MANAZONE)
		manaCreature.AddCondition(cnd.Creature, nil, nil)

		assert.Equal(t, "Soulswap", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)

		firstPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		_, err = scn.WaitForMultipartAction(caster, firstPromptStart)
		require.NoError(t, err)

		secondPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, selectedCreature.ID))
		secondAction, err := scn.WaitForAction(caster, secondPromptStart)
		require.NoError(t, err)
		offeredCardIDs := make([]string, 0, len(secondAction.Cards))
		for _, offeredCard := range secondAction.Cards {
			offeredCardIDs = append(offeredCardIDs, offeredCard.CardID)
		}
		assert.Contains(t, offeredCardIDs, manaCreature.ID)

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, manaCreature.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.MANAZONE, selectedCreature.Zone)
		assert.Equal(t, match.BATTLEZONE, manaCreature.Zone)
		assert.True(t, manaCreature.HasCondition(cnd.SummoningSickness))
	})

	t.Run("can be cancelled without moving a creature", func(t *testing.T) {
		scn, caster, _, spell := setupSoulswapTest(t)
		selectedCreature := putSoulswapTestCardInZone(t, scn, caster.Player, scowlingTomatoUID, match.BATTLEZONE)
		selectedCreature.AddCondition(cnd.Creature, nil, nil)

		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(caster))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, selectedCreature.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("resolves when the moved creature leaves no eligible replacement", func(t *testing.T) {
		scn, caster, opponent, spell := setupSoulswapTest(t)
		selectedCreature := putSoulswapTestCardInZone(t, scn, opponent.Player, gaulezalDragonUID, match.BATTLEZONE)
		selectedCreature.AddCondition(cnd.Creature, nil, nil)

		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, selectedCreature.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.MANAZONE, selectedCreature.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("does not continue when moving the chosen creature is prevented", func(t *testing.T) {
		scn, caster, opponent, spell := setupSoulswapTest(t)
		selectedCreature := putSoulswapTestCardInZone(t, scn, opponent.Player, scowlingTomatoUID, match.BATTLEZONE)
		selectedCreature.AddCondition(cnd.Creature, nil, nil)
		manaCreature := putSoulswapTestCardInZone(t, scn, opponent.Player, shamanBroccoliUID, match.MANAZONE)
		manaCreature.AddCondition(cnd.Creature, nil, nil)
		selectedCreature.Use(func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.MoveCard); ok &&
				event.CardID == card.ID &&
				event.From == match.BATTLEZONE &&
				event.To == match.MANAZONE {
				ctx.InterruptFlow()
			}
		})

		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, selectedCreature.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, selectedCreature.Zone)
		assert.Equal(t, match.MANAZONE, manaCreature.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("ignores non-creature cards in the battle zone", func(t *testing.T) {
		scn, caster, _, spell := setupSoulswapTest(t)
		nonCreature := putSoulswapTestCardInZone(t, scn, caster.Player, thirstForTheHuntUID, match.BATTLEZONE)

		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))

		assert.Equal(t, match.BATTLEZONE, nonCreature.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}

func setupSoulswapTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	caster := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(caster.Player))
	caster.Player.SpawnCard(soulswapUID, match.HAND)
	for range 3 {
		caster.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
	}
	spell, err := scn.FindCard(caster.Player, match.HAND, soulswapUID)
	require.NoError(t, err)
	return scn, caster, opponent, spell
}

func putSoulswapTestCardInZone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string, zone string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, zone, "soulswap_test_setup")
	require.NoError(t, err)
	require.Equal(t, zone, moved.Zone)
	return moved
}
