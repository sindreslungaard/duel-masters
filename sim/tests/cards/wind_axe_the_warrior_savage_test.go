package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	windAxeTheWarriorSavageUID        = "469b3b73-d071-4318-8e7d-71c848e2b318"
	windAxeTheWarriorSavageBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker)
	windAxeTheWarriorSavageTopUID     = "7956b4f5-b910-403d-b388-b67c837b7e99" // Scissor Eye
	windAxeTheWarriorSavageSetupSrc   = "wind_axe_the_warrior_savage_test_setup"
)

func TestWindAxeTheWarriorSavage(t *testing.T) {
	t.Run("destroys an opposing blocker and charges the top card", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		blocker := putCardInBattlezone(t, scn, opponent.Player, windAxeTheWarriorSavageBlockerUID, windAxeTheWarriorSavageSetupSrc)

		windAxe := spawnMulticolorCardInHand(t, scn, player, windAxeTheWarriorSavageUID)
		spawnCivMana(t, player, civ.Fire, 3)
		spawnCivMana(t, player, civ.Nature, 3)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		manaCard := putCardOnTopOfDeck(t, scn, player, windAxeTheWarriorSavageTopUID, windAxeTheWarriorSavageSetupSrc)

		assert.Equal(t, "Wind Axe, the Warrior Savage", windAxe.Name)
		assert.Equal(t, 2000, windAxe.Power)
		assert.Equal(t, 5, windAxe.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, windAxe.Civs)
		assert.True(t, windAxe.HasFamily(family.Human))
		assert.True(t, windAxe.HasFamily(family.BeastFolk))

		require.NoError(t, scn.ActionPlayCard(player, windAxe.ID))

		assert.Equal(t, match.BATTLEZONE, windAxe.Zone)
		assert.Equal(t, match.GRAVEYARD, blocker.Zone)
		assert.Equal(t, match.MANAZONE, manaCard.Zone)
	})

	t.Run("charges the top card even with no blocker to destroy", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		windAxe := spawnMulticolorCardInHand(t, scn, player, windAxeTheWarriorSavageUID)
		spawnCivMana(t, player, civ.Fire, 3)
		spawnCivMana(t, player, civ.Nature, 3)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		manaCard := putCardOnTopOfDeck(t, scn, player, windAxeTheWarriorSavageTopUID, windAxeTheWarriorSavageSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, windAxe.ID))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		actions := 0
		for _, header := range headers {
			if header == "action" {
				actions++
			}
		}
		assert.Equal(t, 1, actions, "only the mana payment is prompted")
		assert.Equal(t, match.MANAZONE, manaCard.Zone)
	})
}
