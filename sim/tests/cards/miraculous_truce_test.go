package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	miraculousTruceUID   = "331e0cb8-5a6d-46a8-ad4f-0ab98c774dc6"
	truceFireCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire)
	truceWaterBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (water)
	truceSetupSrc        = "miraculous_truce_test_setup"
)

// The five civilizations are offered in this order by the card.
const truceFireChoice = 3

func TestMiraculousTruce(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(miraculousTruceUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Miraculous Truce", spell.Name)
		assert.Equal(t, 5, spell.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, spell.Civs)
		assert.True(t, spell.IsMulticolored())
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(miraculousTruceUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, truceSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("creatures of the chosen civilization cannot attack the caster", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, opponent.Player, truceFireCreatureUID, truceSetupSrc)

		castSpell(t, scn, player, miraculousTruceUID)
		require.NoError(t, scn.SubmitChoice(player, truceFireChoice))
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))

		_, err := scn.ActionAttackPlayer(opponent, attacker.ID)
		assert.Error(t, err, "a fire creature should be refused")
		assert.False(t, attacker.Tapped, "a refused attack does not tap it")
	})

	t.Run("other civilizations are unaffected", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, opponent.Player, truceWaterBlockerUID, truceSetupSrc)

		castSpell(t, scn, player, miraculousTruceUID)
		require.NoError(t, scn.SubmitChoice(player, truceFireChoice))
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err, "a water creature is not covered by the truce")
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))
	})

	t.Run("a creature that arrives after the spell is covered too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		castSpell(t, scn, player, miraculousTruceUID)
		require.NoError(t, scn.SubmitChoice(player, truceFireChoice))
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))

		// Put into play during the opponent's turn, after the spell resolved.
		latecomer := putCardInBattlezone(t, scn, opponent.Player, truceFireCreatureUID, truceSetupSrc)

		_, err := scn.ActionAttackPlayer(opponent, latecomer.ID)
		assert.Error(t, err, "the printed text covers creatures put into play later")
	})

	t.Run("the truce is over by the caster's next turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, opponent.Player, truceFireCreatureUID, truceSetupSrc)

		castSpell(t, scn, player, miraculousTruceUID)
		require.NoError(t, scn.SubmitChoice(player, truceFireChoice))
		require.NoError(t, scn.WaitForEventLoop())

		// The caster's turn, the opponent's turn, then the caster's turn again:
		// the effect ends as that last one begins.
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.NoError(t, scn.ActionEndTurn(player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err, "the truce has expired")
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))
	})
}
