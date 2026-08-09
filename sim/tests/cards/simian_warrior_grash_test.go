package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	simianWarriorGrashUID             = "fb6614a6-644d-4fde-9958-225cf2295d47"
	simianWarriorGrashArmorloidUID    = "206b19fb-b389-44f3-9a45-06a7df6dc0f0"
	simianWarriorGrashDestroySpellUID = "e3d5a9f5-cfd1-491d-baa3-2b14b66171f2"
	simianWarriorGrashNonArmorloidUID = "79c48731-193b-4dc6-b26f-1eb820357367"
)

func TestSimianWarriorGrash(t *testing.T) {
	t.Run("burns mana when another own Armorloid is destroyed", func(t *testing.T) {
		scn, owner, opponent, grash, spell, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)
		target := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashArmorloidUID)

		assert.Equal(t, "Simian Warrior Grash", grash.Name)
		assert.Equal(t, 3000, grash.Power)
		assert.Equal(t, 4, grash.ManaCost)
		assert.Equal(t, []string{civ.Fire}, grash.Civs)
		assert.True(t, grash.HasFamily(family.Armorloid))

		resolveSimianWarriorGrashTestDestruction(t, scn, owner, opponent, spell, target, opponentMana, true)

		assert.Equal(t, match.BATTLEZONE, grash.Zone)
		assert.Equal(t, match.GRAVEYARD, target.Zone)
		assert.Equal(t, match.GRAVEYARD, opponentMana.Zone)
	})

	t.Run("burns mana when Grash itself is destroyed", func(t *testing.T) {
		scn, owner, opponent, grash, spell, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)

		resolveSimianWarriorGrashTestDestruction(t, scn, owner, opponent, spell, grash, opponentMana, true)

		assert.Equal(t, match.GRAVEYARD, grash.Zone)
		assert.Equal(t, match.GRAVEYARD, opponentMana.Zone)
	})

	t.Run("multiple copies each trigger for the same Armorloid", func(t *testing.T) {
		scn, owner, opponent, _, spell, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)
		secondGrash := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashUID)
		target := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashArmorloidUID)
		opponentManaCards, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		require.Len(t, opponentManaCards, 2)

		resolveSimianWarriorGrashTestDestruction(t, scn, owner, opponent, spell, target, opponentMana, true)

		assert.Equal(t, match.BATTLEZONE, secondGrash.Zone)
		for _, manaCard := range opponentManaCards {
			assert.Equal(t, match.GRAVEYARD, manaCard.Zone)
		}
	})

	for _, zone := range []string{match.HAND, match.MANAZONE, match.GRAVEYARD, match.DECK, match.SHIELDZONE} {
		t.Run("does not trigger from "+zone, func(t *testing.T) {
			scn, owner, opponent, grash, spell, opponentMana := setupSimianWarriorGrashTest(t, zone)
			target := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashArmorloidUID)

			resolveSimianWarriorGrashTestDestruction(t, scn, owner, opponent, spell, target, opponentMana, false)

			assert.Equal(t, zone, grash.Zone)
			assert.Equal(t, match.MANAZONE, opponentMana.Zone)
		})
	}

	t.Run("does not trigger for an opponent's Armorloid", func(t *testing.T) {
		scn, _, opponent, grash, _, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)
		target := putSimianWarriorGrashTestCardInBattlezone(t, scn, opponent.Player, simianWarriorGrashArmorloidUID)

		scn.Match.Destroy(target, grash, match.DestroyedByMiscAbility)

		assert.Equal(t, match.GRAVEYARD, target.Zone)
		assert.Equal(t, match.MANAZONE, opponentMana.Zone)
	})

	t.Run("does not trigger for a non-Armorloid", func(t *testing.T) {
		scn, owner, _, grash, _, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)
		target := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashNonArmorloidUID)

		scn.Match.Destroy(target, grash, match.DestroyedByMiscAbility)

		assert.Equal(t, match.GRAVEYARD, target.Zone)
		assert.Equal(t, match.MANAZONE, opponentMana.Zone)
	})

	t.Run("does not trigger when destruction is replaced", func(t *testing.T) {
		scn, owner, _, grash, _, opponentMana := setupSimianWarriorGrashTest(t, match.BATTLEZONE)
		target := putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashArmorloidUID)
		target.Use(fx.When(fx.WouldBeDestroyed, fx.ReturnToMana))

		scn.Match.Destroy(target, grash, match.DestroyedByMiscAbility)

		assert.Equal(t, match.MANAZONE, target.Zone)
		assert.Equal(t, match.MANAZONE, opponentMana.Zone)
	})
}

func setupSimianWarriorGrashTest(t *testing.T, grashZone string) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card, *match.Card, *match.Card) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))
	clearSimianWarriorGrashTestHand(t, owner.Player)

	var grash *match.Card
	if grashZone == match.BATTLEZONE {
		grash = putSimianWarriorGrashTestCardInBattlezone(t, scn, owner.Player, simianWarriorGrashUID)
	} else {
		owner.Player.SpawnCard(simianWarriorGrashUID, grashZone)
		var err error
		grash, err = scn.FindCard(owner.Player, grashZone, simianWarriorGrashUID)
		require.NoError(t, err)
	}

	owner.Player.SpawnCard(simianWarriorGrashDestroySpellUID, match.HAND)
	for range 3 {
		owner.Player.SpawnCard(simianWarriorGrashDestroySpellUID, match.MANAZONE)
	}
	spell, err := scn.FindCard(owner.Player, match.HAND, simianWarriorGrashDestroySpellUID)
	require.NoError(t, err)

	for range 2 {
		opponent.Player.SpawnCard(simianWarriorGrashDestroySpellUID, match.MANAZONE)
	}
	opponentMana, err := scn.FindCard(opponent.Player, match.MANAZONE, simianWarriorGrashDestroySpellUID)
	require.NoError(t, err)

	return scn, owner, opponent, grash, spell, opponentMana
}

func clearSimianWarriorGrashTestHand(t *testing.T, player *match.Player) {
	t.Helper()

	hand, err := player.Container(match.HAND)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), hand...) {
		_, err := player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, "simian_warrior_grash_test_setup")
		require.NoError(t, err)
	}
}

func putSimianWarriorGrashTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "simian_warrior_grash_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

func resolveSimianWarriorGrashTestDestruction(
	t *testing.T,
	scn *scenario.TestScenario,
	owner *match.PlayerReference,
	opponent *match.PlayerReference,
	spell *match.Card,
	target *match.Card,
	opponentMana *match.Card,
	expectTrigger bool,
) {
	t.Helper()

	ownerPromptStart, err := scn.MessageCount(owner)
	require.NoError(t, err)
	opponentPromptStart, err := scn.MessageCount(opponent)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(owner, spell.ID))

	destroyAction, err := scn.LatestAction(owner, ownerPromptStart)
	require.NoError(t, err)
	offered := make([]string, 0, len(destroyAction.Cards))
	for _, card := range destroyAction.Cards {
		offered = append(offered, card.CardID)
	}
	assert.Contains(t, offered, target.ID)

	completionStart, err := scn.MessageCount(owner)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(owner, target.ID))

	if !expectTrigger {
		require.NoError(t, scn.WaitForMessage(owner, completionStart, "state_update"))
		headers, err := scn.MessageHeaders(opponent, opponentPromptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		return
	}

	manaAction, err := scn.WaitForAction(opponent, opponentPromptStart)
	require.NoError(t, err)
	offered = offered[:0]
	for _, card := range manaAction.Cards {
		offered = append(offered, card.CardID)
	}
	assert.Contains(t, offered, opponentMana.ID)

	completionStart, err = scn.MessageCount(owner)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(opponent, opponentMana.ID))
	require.NoError(t, scn.WaitForMessage(owner, completionStart, "state_update"))
}
