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
	tropicCrawlerUID  = "adddbc05-e53f-485b-be16-f94a7c5ddd88"
	tropicAttackerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	tropicSetupSrc    = "tropic_crawler_test_setup"
)

func TestTropicCrawler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, tropicCrawlerUID, tropicSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Tropic Crawler", 3000, 4, []string{civ.Water})
		assert.True(t, card.HasFamily(family.EarthEater))
		assert.True(t, card.HasCondition(cnd.Blocker))
		assert.True(t, card.HasCondition(cnd.CantAttackPlayers))
		assert.True(t, card.HasCondition(cnd.CantAttackCreatures))
	})

	t.Run("blocking sends one of the opponent's creatures back to their hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		crawler := putCardInBattlezone(t, scn, player.Player, tropicCrawlerUID, tropicSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, tropicAttackerUID, tropicSetupSrc)
		bystander := putCardInBattlezone(t, scn, opponent.Player, tropicAttackerUID, tropicSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err, "the blocker should have been offered")

		// The next prompt belongs to the opponent, so the block is submitted
		// directly and the wait is done on their connection.
		choiceStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, crawler.ID))

		_, err = scn.WaitForAction(opponent, choiceStart)
		require.NoError(t, err, "the opponent should have been asked to choose")
		require.NoError(t, scn.SubmitAction(opponent, bystander.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, bystander.Zone, "the opponent chose, not the blocker's controller")
		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "3000 beats 2000")
		assert.Equal(t, match.BATTLEZONE, crawler.Zone)
		assert.True(t, crawler.Tapped, "blocking taps the blocker")

		shieldsAfter, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shieldsAfter, shieldCount, "a blocked attack never reaches the shields")
	})

	t.Run("the opponent may pick the attacker itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		crawler := putCardInBattlezone(t, scn, player.Player, tropicCrawlerUID, tropicSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, tropicAttackerUID, tropicSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, tropicAttackerUID, tropicSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err)

		choiceStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, crawler.ID))

		_, err = scn.WaitForAction(opponent, choiceStart)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, attacker.ID))
		settleTurn(t, scn)

		// Pulled out of the battle zone before the battle it lost is resolved,
		// so it goes home instead of to the graveyard.
		assert.Equal(t, match.HAND, attacker.Zone)
	})
}
