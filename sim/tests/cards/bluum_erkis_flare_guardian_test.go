package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bluumErkisFlareGuardianUID  = "3231bac1-1bad-4991-9e7d-95249e11d4b6"
	bluumErkisBrainSerumUID     = "7f225860-af37-47ac-9b36-1480872576b6" // Brain Serum: shield trigger spell (Water), "Draw up to 2 cards."
	bluumErkisAmberGrassUID     = "5da6bf2e-9de1-4166-ade0-a9ded5b99612" // Amber Grass: shield trigger creature, no other text
	bluumErkisBondsOfJusticeUID = "cf8b1a10-a425-4670-bb50-af7c4fb132d6" // Bonds of Justice: shield trigger spell (Light), no prompts of its own
	bluumErkisIceVaporUID       = "ab6c7559-1714-4238-a063-393cfe8adc08" // Ice Vapor, Shadow of Anguish: "whenever your opponent casts a spell, he discards and mana burns"
	bluumErkisAlcadeiasUID      = "7d4b64b0-1672-47d4-b54d-1758b0bb08cf" // Alcadeias, Lord of Spirits: "players can't cast spells other than light spells"
	bluumErkisCrypticTotemUID   = "bb339fe8-3e63-4657-8c5d-d2386c22ff38" // Cryptic Totem: "while tapped, your opponent can't use the shield trigger ability"
	bluumErkisBoomerangCometUID = "b275fbf0-5355-45ec-b3a8-a956cf898ae6" // Boomerang Comet: "after you cast this spell, put it into your mana zone instead of your graveyard"
	bluumErkisFillerUID         = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg: the default test deck's filler creature
	bluumErkisSetupSrc          = "bluum_erkis_flare_guardian_test_setup"
)

func TestBluumErkisFlareGuardian(t *testing.T) {

	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, guardian, "Bluum Erkis, Flare Guardian", 8500, 6, []string{civ.Light, civ.Water})
		assert.True(t, guardian.HasFamily(family.Guardian))
		assert.True(t, guardian.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("forces the cast of an opponent's shield trigger spell, benefiting its own controller", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		brainSerum, err := opponent.Player.SpawnCard(bluumErkisBrainSerumUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCountBefore := len(handBefore)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(action.Cards), 2)

		filler := otherOfferedShield(t, action, brainSerum.ID)

		require.NoError(t, scn.ResolveAttack(player, brainSerum.ID, filler))
		// Brain Serum's own "draw up to 2" is answered by the attacker, not by
		// its real owner: Bluum Erkis's controller is the one casting it.
		answerDrawUpTo(t, scn, player, 2, true)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, brainSerum.Zone, "the stolen spell resolves into its real owner's graveyard, not the caster's")
		assert.True(t, opponent.Player == brainSerum.Player, "ownership is handed back to the real owner once it resolves")

		graveyard, err := opponent.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		assert.Contains(t, graveyard, brainSerum)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCountBefore+2, "the attacker drew the cards, not the shield's owner")

		fillerCard, err := opponent.Player.GetCard(filler, match.HAND)
		require.NoError(t, err)
		assert.Equal(t, match.HAND, fillerCard.Zone, "a shield without shield trigger is still simply revealed and reaches its owner's hand")
	})

	t.Run("a creature with shield trigger is unaffected and still offers its owner the normal optional summon", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		amberGrass, err := opponent.Player.SpawnCard(bluumErkisAmberGrassUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(action.Cards), 2)

		filler := otherOfferedShield(t, action, amberGrass.ID)

		require.NoError(t, scn.ResolveAttack(player, amberGrass.ID, filler))

		// Amber Grass reached the opponent's hand and it is the opponent, not
		// the attacker, who is offered its shield trigger.
		answerInTurn(t, scn, opponent, amberGrass.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, amberGrass.Zone, "Bluum Erkis only lets its controller cast spells, so the creature keeps its own shield trigger offer")
		assert.True(t, opponent.Player == amberGrass.Player)
	})

	t.Run("only applies to shields it breaks itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		guardian.Tapped = true // present, but not the attacker
		attacker := putCardInBattlezone(t, scn, player.Player, bluumErkisFillerUID, bluumErkisSetupSrc)
		brainSerum, err := opponent.Player.SpawnCard(bluumErkisBrainSerumUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, brainSerum.ID))

		// An unrelated attacker breaking this shield offers it to its real
		// owner as a normal, optional shield trigger cast; declining leaves
		// it sitting in hand instead of being force-cast by the attacker.
		cancelInTurn(t, scn, opponent)
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, brainSerum.Zone, "an unrelated attacker's broken shield still just reaches hand")
		assert.True(t, opponent.Player == brainSerum.Player)
	})

	t.Run("casts the shield trigger spell mandatorily, with no separate cast-or-decline prompt", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		bondsOfJustice, err := opponent.Player.SpawnCard(bluumErkisBondsOfJusticeUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, bondsOfJustice.ID)

		// Bonds of Justice has no prompts of its own. If casting it were ever
		// offered as an optional choice, resolving the attack would stop and
		// wait on that choice instead of finishing outright.
		require.NoError(t, scn.ResolveAttack(player, bondsOfJustice.ID, filler))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, bondsOfJustice.Zone, "the cast happened and resolved without ever being offered as optional")
		assert.True(t, opponent.Player == bondsOfJustice.Player)
	})

	t.Run("triggers a reactive \"whenever you cast a spell\" ability as its own controller, even when that hurts the caster", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, bluumErkisIceVaporUID, bluumErkisSetupSrc)
		bondsOfJustice, err := opponent.Player.SpawnCard(bluumErkisBondsOfJusticeUID, match.SHIELDZONE)
		require.NoError(t, err)
		manaCard, err := player.Player.SpawnCard(bluumErkisFillerUID, match.MANAZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		require.NotEmpty(t, handBefore, "the attacker needs a card in hand for Ice Vapor's mandatory discard")
		discardTarget := handBefore[0].ID

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, bondsOfJustice.ID)

		require.NoError(t, scn.ResolveAttack(player, bondsOfJustice.ID, filler))

		// Ice Vapor's own printed text, from its controller's (the shield
		// owner's) perspective, is "whenever your opponent casts a spell" —
		// and the attacker is the one casting now, so it is the attacker who
		// is forced through Ice Vapor's mandatory discard (there is no cancel
		// option: even though this is bad for the caster, it still happens).
		answerInTurn(t, scn, player, discardTarget)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, bondsOfJustice.Zone)
		assert.True(t, opponent.Player == bondsOfJustice.Player)

		discarded, err := player.Player.GetCard(discardTarget, match.GRAVEYARD)
		require.NoError(t, err)
		assert.Equal(t, match.GRAVEYARD, discarded.Zone, "the caster, not the shield's owner, paid Ice Vapor's discard")

		manaGraveyard, err := player.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		assert.Contains(t, manaGraveyard, manaCard, "the caster, not the shield's owner, also paid Ice Vapor's mana burn")
	})

	t.Run("does not trigger a reactive \"whenever your opponent casts a spell\" ability owned by the caster", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		// Ice Vapor belongs to the attacker here, so from the attacker's own
		// perspective it watches "your opponent" (the shield's owner) casting
		// a spell. The shield's owner never casts anything: the attacker does.
		putCardInBattlezone(t, scn, player.Player, bluumErkisIceVaporUID, bluumErkisSetupSrc)
		bondsOfJustice, err := opponent.Player.SpawnCard(bluumErkisBondsOfJusticeUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		opponentHandBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		opponentManaBefore, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, bondsOfJustice.ID)

		require.NoError(t, scn.ResolveAttack(player, bondsOfJustice.ID, filler))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, bondsOfJustice.Zone)

		opponentHandAfter, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		// +1 only from the filler shield normally reaching hand (it is not a
		// shield trigger), not from anything Ice Vapor did.
		assert.Len(t, opponentHandAfter, len(opponentHandBefore)+1, "the shield's owner never cast anything, so Ice Vapor never asked them to discard")

		opponentManaAfter, err := opponent.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, opponentManaAfter, len(opponentManaBefore), "the shield's owner never cast anything, so Ice Vapor never mana burned them")
	})

	t.Run("a static \"can't cast\" restriction prevents the cast, and the spell still reaches its owner's graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, bluumErkisAlcadeiasUID, bluumErkisSetupSrc)
		// Brain Serum is a Water spell: Alcadeias, Lord of Spirits ("players
		// can't cast spells other than light spells") should stop the cast.
		brainSerum, err := opponent.Player.SpawnCard(bluumErkisBrainSerumUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, brainSerum.ID)

		// If the cast were not prevented, Brain Serum's own "draw up to 2"
		// would open a prompt here and the attack would not resolve outright.
		require.NoError(t, scn.ResolveAttack(player, brainSerum.ID, filler))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, brainSerum.Zone, "Bluum Erkis still sends it to its owner's graveyard even though the cast itself was prevented")
		assert.True(t, opponent.Player == brainSerum.Player)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore), "the prevented cast never resolved its own effect, so the attacker drew nothing")
	})

	t.Run("a shield-trigger-ability restriction like tapped Cryptic Totem does not block the forced cast", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		totem := putCardInBattlezone(t, scn, opponent.Player, bluumErkisCrypticTotemUID, bluumErkisSetupSrc)
		totem.Tapped = true
		bondsOfJustice, err := opponent.Player.SpawnCard(bluumErkisBondsOfJusticeUID, match.SHIELDZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, bondsOfJustice.ID)

		// Bluum Erkis's ability is not "using the shield trigger ability of
		// that spell", so Cryptic Totem's restriction on the shield trigger
		// ability never comes into play here; the forced cast still resolves
		// outright.
		require.NoError(t, scn.ResolveAttack(player, bondsOfJustice.ID, filler))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, bondsOfJustice.Zone)
		assert.True(t, opponent.Player == bondsOfJustice.Player)
	})

	t.Run("overrides Boomerang Comet's own \"put into mana zone instead of graveyard\", but the caster still keeps its own benefit", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		guardian := putCardInBattlezone(t, scn, player.Player, bluumErkisFlareGuardianUID, bluumErkisSetupSrc)
		boomerangComet, err := opponent.Player.SpawnCard(bluumErkisBoomerangCometUID, match.SHIELDZONE)
		require.NoError(t, err)
		manaCard, err := player.Player.SpawnCard(bluumErkisFillerUID, match.MANAZONE)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, guardian.ID)
		require.NoError(t, err)
		filler := otherOfferedShield(t, action, boomerangComet.ID)

		// "Return a card from your mana zone to your hand" is answered
		// automatically here since the attacker's mana zone holds exactly one
		// legal, mandatory choice.
		require.NoError(t, scn.ResolveAttack(player, boomerangComet.ID, filler))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, boomerangComet.Zone, "its own redirect to the mana zone never applies: Bluum Erkis was always going to put it into its owner's graveyard")
		assert.True(t, opponent.Player == boomerangComet.Player)

		handCard, err := player.Player.GetCard(manaCard.ID, match.HAND)
		require.NoError(t, err)
		assert.Equal(t, match.HAND, handCard.Zone, "the caster, not the shield's owner, got their own mana zone card back")
	})
}

// otherOfferedShield returns the id of a shield in the attack's shield
// selection prompt other than exclude, failing the test if either is missing.
func otherOfferedShield(t *testing.T, action *server.ActionMessage, exclude string) string {
	t.Helper()

	var found bool
	var other string

	for _, c := range action.Cards {
		if c.CardID == exclude {
			found = true
			continue
		}
		if other == "" {
			other = c.CardID
		}
	}

	require.True(t, found, "expected shield to be offered as breakable")
	require.NotEmpty(t, other, "expected another shield to be offered alongside it")

	return other
}
