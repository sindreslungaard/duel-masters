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
	agiraTheWarlordCrawlerUID = "d2cda920-585e-449e-834f-af6e84b573da"
	agiraKinUID               = "da51845c-4a6b-4c36-9c7d-fbb654ba2aa2" // Kanesill, the Explorer (Gladiator, 4000, blocker)
	agiraStrangerUID          = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	agiraOtherBlockerUID      = "85852774-dd96-4395-8980-eb5b85bf5bfc" // Ferrosaturn, Spectral Knight (Rainbow Phantom, blocker, 2000)
	agiraSetupSrc             = "agira_the_warlord_crawler_test_setup"
)

func TestAgiraTheWarlordCrawler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Agira, the Warlord Crawler", 5500, 4, []string{civ.Light, civ.Water})
		assert.True(t, card.HasFamily(family.Gladiator))
		assert.True(t, card.HasFamily(family.EarthEater))
		assert.True(t, card.IsMulticolored())
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(agiraTheWarlordCrawlerUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, agiraSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("its own other kin get +2000", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		agira := putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, agiraKinUID, agiraSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, agiraStrangerUID, agiraSetupSrc)
		theirKin := putCardInBattlezone(t, scn, opponent.Player, agiraKinUID, agiraSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))
		assert.Equal(t, stranger.Power, scn.Match.GetPower(stranger, false), "only its own races")
		assert.Equal(t, theirKin.Power, scn.Match.GetPower(theirKin, false), "only its controller's creatures")
		assert.Equal(t, agira.Power, scn.Match.GetPower(agira, false), "\"other\" spares itself")
	})

	t.Run("the bonus leaves with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		agira := putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, agiraKinUID, agiraSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, kin.Power+2000, scn.Match.GetPower(kin, false))

		_, err := player.Player.MoveCard(agira.ID, match.BATTLEZONE, match.GRAVEYARD, agiraSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, kin.Power, scn.Match.GetPower(kin, false))
	})

	t.Run("a blocking kin may draw a card", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, agiraKinUID, agiraSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, agiraStrangerUID, agiraSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handBefore := len(hand)

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err, "the blocker should have been offered")

		questionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, kin.ID))

		_, err = scn.WaitForAction(player, questionStart)
		require.NoError(t, err, "the draw should have been offered")
		require.NoError(t, scn.SubmitAction(player))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, hand, handBefore+1)
		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "6000 beats 2000")
	})

	t.Run("the draw may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		kin := putCardInBattlezone(t, scn, player.Player, agiraKinUID, agiraSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, agiraStrangerUID, agiraSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handBefore := len(hand)

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err)

		questionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, kin.ID))

		_, err = scn.WaitForAction(player, questionStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, hand, handBefore)
	})

	t.Run("a blocker that is not its kin offers nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, agiraTheWarlordCrawlerUID, agiraSetupSrc)
		stranger := putCardInBattlezone(t, scn, player.Player, agiraOtherBlockerUID, agiraSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, agiraStrangerUID, agiraSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, stranger.HasCondition(cnd.Blocker))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err)

		require.NoError(t, scn.SubmitAction(player, stranger.ID))

		// The event loop only settles if nothing is waiting on an answer.
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "2000 against 2000 destroys both")
		assert.Equal(t, match.GRAVEYARD, stranger.Zone)
	})
}
