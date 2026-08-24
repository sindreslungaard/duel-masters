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
	qTronicOmnistrainUID             = "4c2e1260-29a9-4645-939a-012fa9df932f"
	qTronicOmnistrainSurvivorBaseUID = "d176b30a-cac6-4249-a78d-18f34b97546b" // Promephius Q (Water Survivor, 2000)
	qTronicOmnistrainNonSurvivorUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, 2000)
	qTronicOmnistrainSmashHornQUID   = "5ff3c63c-30df-4f5c-acce-796f5b6c2dac" // Smash Horn Q (Survivor: +1000 power, shared)
	qTronicOmnistrainBronzeArmUID    = "015fd6bb-37a9-45cf-bb6b-a5497412b880" // Bronze-Arm Tribe (Beast Folk, no Survivor keyword)
	qTronicOmnistrainSetupSrc        = "q_tronic_omnistrain_test_setup"
)

func TestQTronicOmnistrain(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Q-tronic Omnistrain", 3000, 6, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.Survivor))
		assert.True(t, card.HasCondition(cnd.ShieldTrigger))
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("it evolves onto a Survivor creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainSurvivorBaseUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		omnistrain := summonWithOwnMana(t, scn, player, qTronicOmnistrainUID)

		assert.Equal(t, match.BATTLEZONE, omnistrain.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone, "the base goes under the evolution")
	})

	t.Run("it cannot evolve onto a non-Survivor creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		wrongBase := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainNonSurvivorUID, qTronicOmnistrainSetupSrc)

		omnistrain, err := player.Player.SpawnCard(qTronicOmnistrainUID, match.HAND)
		require.NoError(t, err)
		for range 6 {
			_, err := player.Player.SpawnCard(qTronicOmnistrainUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.Error(t, scn.ActionPlayCard(player, omnistrain.ID), "there is no legal base to evolve from")
		assert.Equal(t, match.HAND, omnistrain.Zone)
		assert.Equal(t, match.BATTLEZONE, wrongBase.Zone)
	})

	t.Run("grants Survivor to its controller's creatures, in addition to their other races", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		filler := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainNonSurvivorUID, qTronicOmnistrainSetupSrc)
		theirFiller := putCardInBattlezone(t, scn, opponent.Player, qTronicOmnistrainNonSurvivorUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.True(t, filler.HasFamily(family.Survivor))
		assert.True(t, filler.HasFamily(family.Human), "keeps its other races")
		assert.False(t, theirFiller.HasFamily(family.Survivor), "only its controller's creatures")
	})

	t.Run("the grant leaves with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		omnistrain := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		filler := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainNonSurvivorUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, filler.HasFamily(family.Survivor))

		_, err := player.Player.MoveCard(omnistrain.ID, match.BATTLEZONE, match.GRAVEYARD, qTronicOmnistrainSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.False(t, filler.HasFamily(family.Survivor))
		assert.True(t, filler.HasFamily(family.Human))
	})

	t.Run("lets another Survivor evolution use a creature it granted the race to", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		filler := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainNonSurvivorUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, filler.HasFamily(family.Survivor))

		gargantua, err := player.Player.SpawnCard(qTronicGargantuaUID, match.HAND)
		require.NoError(t, err)
		for range 6 {
			_, err := player.Player.SpawnCard(qTronicGargantuaUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionPlayCard(player, gargantua.ID))

		// Both the Omnistrain and the filler it granted the race to are now
		// legal evolution bases, so unlike a single-target evolve this opens a
		// choice instead of resolving on its own.
		_, err = scn.LatestAction(player, 0)
		require.NoError(t, err, "expected a prompt to choose the evolution base")
		require.NoError(t, scn.SubmitAction(player, filler.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, gargantua.Zone)
		assert.Equal(t, match.HIDDENZONE, filler.Zone, "the granted race let it serve as the evolution base")
	})

	// Official ruling: "Your creatures with a Survivor race (with the ability
	// of Q-tronic Omnistrain) will get the abilities from your other creatures
	// that have the Survivor keyword. Your creatures that don't have the
	// Survivor race (without the effect of Q-tronic Omnistrain) will not share
	// any of their abilities to your cards with a (printed) Survivor race."
	t.Run("shares another Survivor's ability with a creature it granted the race to", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		smashHornQ := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainSmashHornQUID, qTronicOmnistrainSetupSrc)
		bronzeArmTribe := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainBronzeArmUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		require.True(t, bronzeArmTribe.HasFamily(family.Survivor), "granted the race by Omnistrain")
		require.False(t, bronzeArmTribe.HasCondition(cnd.Survivor), "it has no printed Survivor keyword of its own")

		assert.Equal(t, bronzeArmTribe.Power+1000, scn.Match.GetPower(bronzeArmTribe, false), "gets Smash Horn Q's shared Survivor ability")
		assert.Equal(t, smashHornQ.Power+1000, scn.Match.GetPower(smashHornQ, false), "and Smash Horn Q keeps its own")
	})

	t.Run("a creature without the printed Survivor race never shares its own ability", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainUID, qTronicOmnistrainSetupSrc)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainSmashHornQUID, qTronicOmnistrainSetupSrc)
		putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainBronzeArmUID, qTronicOmnistrainSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		manaBefore, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		manaCountBefore := len(manaBefore)

		// Bronze-Arm Tribe has no printed Survivor race, so it has no Survivor
		// keyword to share its "put the top card of your deck into your mana
		// zone" ability with, no matter how many real Survivors are in play.
		promephiusQ := putCardInBattlezone(t, scn, player.Player, qTronicOmnistrainSurvivorBaseUID, qTronicOmnistrainSetupSrc)
		settleTurn(t, scn)

		require.True(t, promephiusQ.HasFamily(family.Survivor))

		manaAfter, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, manaAfter, manaCountBefore, "Promephius Q entering must not trigger Bronze-Arm Tribe's ability")
	})
}
