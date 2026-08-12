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
	dolmarksTheShadowWarriorUID      = "32b5aa84-3e0f-4b81-a37a-0ec6f21bba52"
	dolmarksVictimUID                = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	dolmarksTheShadowWarriorSetupSrc = "dolmarks_the_shadow_warrior_test_setup"
)

func TestDolmarksTheShadowWarrior(t *testing.T) {
	t.Run("both players lose a creature and a mana card", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		dolmarks := spawnMulticolorCardInHand(t, scn, player, dolmarksTheShadowWarriorUID)
		spawnCivMana(t, player, civ.Darkness, 3)
		spawnCivMana(t, player, civ.Fire, 3)
		spawnCivMana(t, opponent, civ.Water, 2)

		ownVictim := putCardInBattlezone(t, scn, player.Player, dolmarksVictimUID, dolmarksTheShadowWarriorSetupSrc)
		opponentVictim := putCardInBattlezone(t, scn, opponent.Player, dolmarksVictimUID, dolmarksTheShadowWarriorSetupSrc)

		assert.Equal(t, "Dolmarks, the Shadow Warrior", dolmarks.Name)
		assert.Equal(t, 4000, dolmarks.Power)
		assert.Equal(t, 4, dolmarks.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, dolmarks.Civs)
		assert.True(t, dolmarks.HasFamily(family.Ghost))
		assert.True(t, dolmarks.HasFamily(family.Human))

		playerManaBefore, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		playerManaCount := len(playerManaBefore)
		opponentManaBefore, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		opponentManaCount := len(opponentManaBefore)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, dolmarks.ID))

		// Dolmarks is in the battle zone by now, so its controller really does get
		// to pick which of their two creatures dies.
		destroyAction, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.False(t, destroyAction.Cancellable, "the destruction is mandatory")
		offered := make([]string, 0, len(destroyAction.Cards))
		for _, card := range destroyAction.Cards {
			offered = append(offered, card.CardID)
		}
		assert.ElementsMatch(t, []string{dolmarks.ID, ownVictim.ID}, offered)

		manaPromptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, ownVictim.ID))

		manaAction, err := scn.WaitForAction(player, manaPromptStart)
		require.NoError(t, err)
		require.NotEmpty(t, manaAction.Cards)

		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, manaAction.Cards[0].CardID))

		// The opponent controls a single creature, so that half is forced and
		// their only prompt is the mana card.
		opponentManaAction, err := scn.WaitForAction(opponent, opponentStart)
		require.NoError(t, err)
		require.NotEmpty(t, opponentManaAction.Cards)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, opponentManaAction.Cards[0].CardID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, dolmarks.Zone)
		assert.Equal(t, match.GRAVEYARD, ownVictim.Zone, "the controller destroys one of their own")
		assert.Equal(t, match.GRAVEYARD, opponentVictim.Zone, "and the opponent destroys one of theirs")

		playerManaAfter, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, playerManaAfter, playerManaCount-1)
		opponentManaAfter, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, opponentManaAfter, opponentManaCount-1)
	})

	t.Run("resolves when neither player has a creature or spare mana to lose", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		dolmarks := spawnMulticolorCardInHand(t, scn, player, dolmarksTheShadowWarriorUID)
		spawnCivMana(t, player, civ.Darkness, 2)
		spawnCivMana(t, player, civ.Fire, 2)
		clearZone(t, opponent.Player, match.MANAZONE, dolmarksTheShadowWarriorSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, dolmarks.ID))

		// Dolmarks itself is the controller's only creature, so it destroys itself.
		manaAction, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		require.NotEmpty(t, manaAction.Cards)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, manaAction.Cards[0].CardID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.GRAVEYARD, dolmarks.Zone)
		opponentMana, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Empty(t, opponentMana)
	})
}
