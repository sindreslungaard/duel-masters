package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Macho Melon's wave striker ability is "Power attacker +3000", which makes
	// it a convenient stand-in for the mechanic: whether the ability is on is
	// just a number.
	machoMelonUID = "fa987e39-2955-4074-bcf2-b7888ae27319"
	// Merlee is a second wave striker to make up the count, and Immortal Baron
	// Vorg is a creature that is not one.
	merleeUID           = "9ec241a3-57e9-4054-8680-aff2c1a7b45b"
	waveStrikerSetupSrc = "wave_striker_test_setup"
)

func TestWaveStriker(t *testing.T) {
	t.Run("the ability is off below the threshold", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)

		passTurn(t, scn, player, opponent)
		require.True(t, melon.HasCondition(cnd.WaveStriker), "it is a wave striker either way")

		assert.Equal(t, 1000, scn.Match.GetPower(melon, true), "alone it has no ability")

		putInBattleZone(t, scn, player, machoMelonUID)
		assert.Equal(t, 1000, scn.Match.GetPower(melon, true), "one other is not enough")
	})

	t.Run("two other wave strikers switch the ability on", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)

		passTurn(t, scn, player, opponent)

		assert.Equal(t, 4000, scn.Match.GetPower(melon, true), "power attacker +3000 is switched on")
		assert.Equal(t, 1000, scn.Match.GetPower(melon, false), "it is still only a power attacker")
	})

	t.Run("the opponent's wave strikers count towards the threshold", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, opponent, machoMelonUID)
		putInBattleZone(t, scn, opponent, machoMelonUID)

		passTurn(t, scn, player, opponent)

		// The printed wording is "in the battle zone", not "you have".
		assert.Equal(t, 4000, scn.Match.GetPower(melon, true))
	})

	t.Run("creatures without the keyword do not count", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, immortalBaronVorgUID)
		putInBattleZone(t, scn, player, immortalBaronVorgUID)

		passTurn(t, scn, player, opponent)

		assert.Equal(t, 1000, scn.Match.GetPower(melon, true))
	})

	t.Run("the ability follows the board in both directions", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		second := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)

		// Held back to restore the count later. It has to be around for an
		// untap step to count as a wave striker at all, exactly as a card drawn
		// out of the deck would have been.
		spare, err := player.Player.SpawnCard(machoMelonUID, match.HAND)
		require.NoError(t, err)

		passTurn(t, scn, player, opponent)
		require.Equal(t, 4000, scn.Match.GetPower(melon, true))

		// Losing one drops the count back under the threshold.
		_, err = player.Player.MoveCard(second.ID, match.BATTLEZONE, match.GRAVEYARD, waveStrikerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, 1000, scn.Match.GetPower(melon, true), "the ability switches back off")

		// And it comes back when the count is restored.
		_, err = player.Player.MoveCard(spare.ID, match.HAND, match.BATTLEZONE, waveStrikerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, 4000, scn.Match.GetPower(melon, true))
	})

	t.Run("a wave striker outside the battle zone has no ability", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)

		passTurn(t, scn, player, opponent)
		require.Equal(t, 4000, scn.Match.GetPower(melon, true))

		_, err := player.Player.MoveCard(melon.ID, match.BATTLEZONE, match.HAND, waveStrikerSetupSrc)
		require.NoError(t, err)

		assert.Equal(t, 1000, scn.Match.GetPower(melon, true))
	})

	t.Run("a granted condition is cleaned up when the count drops", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		melon := putInBattleZone(t, scn, player, machoMelonUID)
		second := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)

		passTurn(t, scn, player, opponent)
		require.True(t, melon.HasCondition(cnd.PowerAttacker), "the keyword is granted, not just the number")

		_, err := player.Player.MoveCard(second.ID, match.BATTLEZONE, match.GRAVEYARD, waveStrikerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, melon.HasCondition(cnd.PowerAttacker), "the grant is withdrawn, not left behind")
	})

	t.Run("a grant to other creatures is withdrawn the same way", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		merlee := putInBattleZone(t, scn, player, merleeUID)
		second := putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)
		ally := putInBattleZone(t, scn, player, immortalBaronVorgUID)
		theirs := putInBattleZone(t, scn, opponent, immortalBaronVorgUID)

		passTurn(t, scn, player, opponent)

		assert.Equal(t, 3000, scn.Match.GetPower(ally, false), "2000 plus Merlee's 1000")
		assert.Equal(t, 2500, scn.Match.GetPower(merlee, false), "each of your creatures includes Merlee")
		assert.Equal(t, 2000, scn.Match.GetPower(theirs, false), "only its controller's creatures")

		_, err := player.Player.MoveCard(second.ID, match.BATTLEZONE, match.GRAVEYARD, waveStrikerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, 2000, scn.Match.GetPower(ally, false), "the bonus goes when the ability does")
	})

	t.Run("a creature that leaves play does not keep the bonus", func(t *testing.T) {
		scn, player, opponent := waveStrikerScenario(t)
		putInBattleZone(t, scn, player, merleeUID)
		putInBattleZone(t, scn, player, machoMelonUID)
		putInBattleZone(t, scn, player, machoMelonUID)
		ally := putInBattleZone(t, scn, player, immortalBaronVorgUID)

		passTurn(t, scn, player, opponent)
		require.True(t, ally.HasCondition(cnd.PowerAmplifier))

		_, err := player.Player.MoveCard(ally.ID, match.BATTLEZONE, match.HAND, waveStrikerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, ally.HasCondition(cnd.PowerAmplifier), "a bounced creature must not carry it back into play")
	})
}

func waveStrikerScenario(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

func putInBattleZone(t *testing.T, _ *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	moved, err := player.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, waveStrikerSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	return moved
}

func passTurn(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference) {
	t.Helper()

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))
}
