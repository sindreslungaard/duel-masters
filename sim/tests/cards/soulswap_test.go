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
const phalEegaUID = "1960f3c2-5321-4da1-be56-50a7e98cae0d"
const soulswapTestSrc = "soulswap_test_setup"

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

	t.Run("does not let a triggered random discard consume the resolving spell itself", func(t *testing.T) {
		// Jagila, the Hidden Pillager makes its controller's opponent discard 3
		// random cards from hand when it enters the battle zone. Soulswap can
		// put it there directly from its owner's mana zone, which makes the
		// caster of Soulswap the one who gets hit by that discard. Soulswap
		// itself must already have left the caster's hand for the graveyard by
		// then, or it becomes a candidate for its own opponent's random discard.
		scn, caster, opponent := setupDuel(t)

		// The wave striker fillers and the redirect target need an untap step
		// before Jagila's own "2 or more other creatures" check switches on, so
		// everything opponent-side is placed before passing the turn.
		addWaveStrikerFillers(t, scn, opponent, 2)
		redirectTarget := putCardInBattlezone(t, scn, opponent.Player, scowlingTomatoUID, soulswapTestSrc)
		jagila, err := opponent.Player.SpawnCard(jagilaUID, match.MANAZONE)
		require.NoError(t, err)
		for range 3 {
			_, err := opponent.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
			require.NoError(t, err)
		}
		passTurnToSelf(t, scn, caster, opponent)

		// Emptied last so the drawn opening hand cannot be discarded instead of
		// Soulswap, and so Soulswap is the only card in hand when Jagila's
		// discard fires.
		emptyHand(t, caster, soulswapTestSrc)
		caster.Player.SpawnCard(soulswapUID, match.HAND)
		for range 3 {
			caster.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
		}
		spell, err := scn.FindCard(caster.Player, match.HAND, soulswapUID)
		require.NoError(t, err)

		// Records which card's move first put Soulswap in the graveyard. If
		// it's Jagila's discard rather than Soulswap's own resolution, the bug
		// is present.
		var graveyardMoveSource string
		spell.Use(func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.CardMoved); ok &&
				event.CardID == card.ID && event.To == match.GRAVEYARD && graveyardMoveSource == "" {
				graveyardMoveSource = event.Source
			}
		})

		firstPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		_, err = scn.WaitForMultipartAction(caster, firstPromptStart)
		require.NoError(t, err)

		secondPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, redirectTarget.ID))
		secondAction, err := scn.WaitForAction(caster, secondPromptStart)
		require.NoError(t, err)
		offeredCardIDs := make([]string, 0, len(secondAction.Cards))
		for _, offeredCard := range secondAction.Cards {
			offeredCardIDs = append(offeredCardIDs, offeredCard.CardID)
		}
		require.Contains(t, offeredCardIDs, jagila.ID)

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, jagila.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, spell.Zone)
		assert.Equal(t, spell.ID, graveyardMoveSource, "Soulswap should reach the graveyard through its own resolution, not Jagila's random discard")

		hand, err := caster.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand, "Jagila's discard had nothing left to take once Soulswap had already left hand")
	})

	t.Run("lets a creature it puts into play recur the spell from the graveyard", func(t *testing.T) {
		// Phal Eega, Dawn Guardian returns a spell from its controller's own
		// graveyard to hand when it enters the battle zone. For Phal Eega to
		// end up under the caster's control, the caster has to redirect one of
		// their own battle zone creatures (moving it into their own mana zone)
		// and bring Phal Eega in from that same mana zone. By the time Phal
		// Eega's trigger resolves, Soulswap itself must already be sitting in
		// the caster's graveyard and offered as a legal target for its own
		// recursion.
		scn, caster, _, spell := setupSoulswapTest(t)
		// Phal Eega's filter looks for cnd.Spell, which is otherwise only
		// rebuilt on an untap step; set it directly since this test does not
		// pass a turn.
		spell.AddCondition(cnd.Spell, nil, spell.ID)
		redirectTarget := putSoulswapTestCardInZone(t, scn, caster.Player, scowlingTomatoUID, match.BATTLEZONE)
		redirectTarget.AddCondition(cnd.Creature, nil, nil)
		phalEega := putSoulswapTestCardInZone(t, scn, caster.Player, phalEegaUID, match.MANAZONE)
		phalEega.AddCondition(cnd.Creature, nil, nil)
		for range 3 {
			c := putSoulswapTestCardInZone(t, scn, caster.Player, scowlingTomatoUID, match.MANAZONE)
			c.AddCondition(cnd.Creature, nil, nil)
		}

		firstPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(caster, spell.ID))
		_, err = scn.WaitForMultipartAction(caster, firstPromptStart)
		require.NoError(t, err)

		secondPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, redirectTarget.ID))
		secondAction, err := scn.WaitForAction(caster, secondPromptStart)
		require.NoError(t, err)
		offeredCardIDs := make([]string, 0, len(secondAction.Cards))
		for _, offeredCard := range secondAction.Cards {
			offeredCardIDs = append(offeredCardIDs, offeredCard.CardID)
		}
		require.Contains(t, offeredCardIDs, phalEega.ID)

		thirdPromptStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, phalEega.ID))
		thirdAction, err := scn.WaitForAction(caster, thirdPromptStart)
		require.NoError(t, err, "Phal Eega's own effect should prompt to return a spell from the graveyard")
		recurCardIDs := make([]string, 0, len(thirdAction.Cards))
		for _, offeredCard := range thirdAction.Cards {
			recurCardIDs = append(recurCardIDs, offeredCard.CardID)
		}
		assert.Contains(t, recurCardIDs, spell.ID, "Soulswap should already be in the graveyard by the time Phal Eega's own ability resolves")

		completionStart, err := scn.MessageCount(caster)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(caster, spell.ID))
		require.NoError(t, scn.WaitForMessage(caster, completionStart, "state_update"))

		assert.Equal(t, match.HAND, spell.Zone)
		assert.Equal(t, match.BATTLEZONE, phalEega.Zone)
	})

	t.Run("cast for free from a broken shield, still lets a fetched creature recur it", func(t *testing.T) {
		// Reported bug, in its original shape: "if opponent breaks into my
		// shield and it's swap and I choose to swap in my own Phal Eega, I
		// should be able to retrieve that same swap back but it currently
		// doesn't work."
		scn, defender, attacker := setupDuel(t)

		shieldSpell, err := defender.Player.SpawnCard(soulswapUID, match.SHIELDZONE)
		require.NoError(t, err)
		redirectTarget := putCardInBattlezone(t, scn, defender.Player, scowlingTomatoUID, soulswapTestSrc)
		phalEega, err := defender.Player.SpawnCard(phalEegaUID, match.MANAZONE)
		require.NoError(t, err)
		for range 3 {
			_, err := defender.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
			require.NoError(t, err)
		}
		raider := putCardInBattlezone(t, scn, attacker.Player, scowlingTomatoUID, soulswapTestSrc)

		// A single untap step, on the way into attacker's turn, gives the
		// shield spell its cnd.ShieldTrigger/cnd.Spell, everything else its
		// cnd.Creature, and clears the raider's summoning sickness.
		require.NoError(t, scn.ActionEndTurn(defender))

		_, err = scn.ActionAttackPlayer(attacker, raider.ID)
		require.NoError(t, err)

		triggerPromptStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(attacker, shieldSpell.ID))

		triggerAction, err := scn.WaitForAction(defender, triggerPromptStart)
		require.NoError(t, err, "expected the shield trigger prompt offering Soulswap")
		triggerCardIDs := make([]string, 0, len(triggerAction.Cards))
		for _, offered := range triggerAction.Cards {
			triggerCardIDs = append(triggerCardIDs, offered.CardID)
		}
		require.Contains(t, triggerCardIDs, shieldSpell.ID)

		firstPromptStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(defender, shieldSpell.ID))
		_, err = scn.WaitForMultipartAction(defender, firstPromptStart)
		require.NoError(t, err)

		secondPromptStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(defender, redirectTarget.ID))
		secondAction, err := scn.WaitForAction(defender, secondPromptStart)
		require.NoError(t, err)
		offeredCardIDs := make([]string, 0, len(secondAction.Cards))
		for _, offered := range secondAction.Cards {
			offeredCardIDs = append(offeredCardIDs, offered.CardID)
		}
		require.Contains(t, offeredCardIDs, phalEega.ID)

		thirdPromptStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(defender, phalEega.ID))
		thirdAction, err := scn.WaitForAction(defender, thirdPromptStart)
		require.NoError(t, err, "expected Phal Eega's own effect to prompt to return a spell from the graveyard")
		recurCardIDs := make([]string, 0, len(thirdAction.Cards))
		for _, offered := range thirdAction.Cards {
			recurCardIDs = append(recurCardIDs, offered.CardID)
		}
		assert.Contains(t, recurCardIDs, shieldSpell.ID, "Soulswap should already be in the graveyard by the time Phal Eega's own ability resolves")

		completionStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(defender, shieldSpell.ID))
		require.NoError(t, scn.WaitForMessage(defender, completionStart, "state_update"))

		assert.Equal(t, match.HAND, shieldSpell.Zone)
		assert.Equal(t, match.BATTLEZONE, phalEega.Zone)
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
	moved, err := player.MoveCard(card.ID, match.HAND, zone, soulswapTestSrc)
	require.NoError(t, err)
	require.Equal(t, zone, moved.Zone)
	return moved
}
