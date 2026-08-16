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
	meloppeUID      = "0451a36e-fe88-4817-ad94-3dbd9e460fb8"
	meloppeSetupSrc = "meloppe_test_setup"

	// meloppeVanillaAttackerUID is Lah, Purification Enforcer: a plain 5500
	// power creature with no effects of its own, used as an attacker that
	// isn't Meloppe itself.
	meloppeVanillaAttackerUID = "0cc5279e-0a26-41a8-a2a5-f7711120b772"

	// meloppeDoubleBreakerUID is Astrocomet Dragon: double breaker with no
	// other effect, used to confirm the shield count survives the swap.
	meloppeDoubleBreakerUID = "91db2302-6794-4aa4-b17b-6637d356e9ac"

	// meloppeAquaMasterUID is Aqua Master: "Whenever this creature attacks a
	// player and isn't blocked, choose one of your opponent's shields and turn
	// it face up." Used to confirm Meloppe also reverses the chooser for a
	// shield-choosing effect that isn't the default shield-break selection.
	meloppeAquaMasterUID = "2b33507e-0154-43f7-a190-e84178ea81c9"

	// meloppeCosmicDartsUID is Cosmic Darts: "Your opponent chooses one of
	// your shields..." Unlike every other case here, the printed default
	// chooser is the caster's *opponent*, choosing among the *caster's* own
	// shields. Used to confirm Meloppe's other clause: when the caster
	// controls Meloppe, the caster chooses instead of letting the opponent
	// pick.
	meloppeCosmicDartsUID = "0fe53ab5-bef9-4bbd-bb05-c94ccb9b1342"

	// meloppeInvincibleCataclysmUID is Invincible Cataclysm: "Choose up to 3
	// of your opponent's shields and put them into his graveyard." A
	// cross-player choice made outside any attack, to confirm the general
	// ChooseShieldEvent mechanism (not just SelectShields) is reversed.
	meloppeInvincibleCataclysmUID = "d8aee1e9-0799-4acf-af2a-3a8fd1b4eb8e"

	// meloppeReconOperationUID is Recon Operation: "Look at up to 3 of your
	// opponent's shields." A non-mutating "look" effect, to confirm the fix
	// covers effects that never move or break anything.
	meloppeReconOperationUID = "aa1cca8a-3140-4180-9c9f-3f126dd16a68"

	// meloppeGajirabuteUID is Gajirabute, Vile Centurion: "When it comes into
	// play, choose one of your opponent's shields and put it into his
	// graveyard." Exercises fx.PutOpShieldIntoGraveyard, a shared helper fixed
	// once for every card that calls it.
	meloppeGajirabuteUID = "e2552e90-61bf-46ec-80a6-28666bfccd1c"

	// meloppePointaUID is Pointa, the Aqua Shadow: "When it comes into play,
	// look at 1 of your opponent's shields. Then your opponent discards a
	// card at random from their hand." Exercises fx.ShowXShields, another
	// shared helper.
	meloppePointaUID = "bc128af9-0fc2-4a1b-b10e-2f695f05c24e"

	// meloppeStingerBallUID is Stinger Ball: "Whenever this creature attacks,
	// you may look at 1 of your opponent's shields." Exercises
	// fx.WheneverThisAttacksMayLookAtOpShield, which fires on AttackConfirmed
	// rather than SelectShields, to confirm the fix isn't tied to the
	// shield-break flow specifically.
	meloppeStingerBallUID = "fb3a5abd-c16a-4365-98c0-79753fdcd91d"

	// meloppeAdomisUID is Adomis, the Oracle: "$tap Choose a shield and look
	// at it." Names no owner, so the offer spans both shield zones in one
	// blind pick. Used to confirm that when the picker's blind pick lands on
	// the opponent's zone, the opponent chooses which of their own shields is
	// affected instead.
	meloppeAdomisUID = "c30f60d5-d97c-457e-b438-e4a606136171"

	// meloppeUlarusUID is Ularus, Punishment Elemental: "for each creature you
	// have in the battle zone, you may choose a shield and turn it face up."
	// Same owner-agnostic pool as Adomis, but multi-select instead of a single
	// pick.
	meloppeUlarusUID = "05c5496d-e5fa-4691-8542-2d6c6919f402"
)

func TestMeloppe(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Meloppe", 1000, 3, []string{civ.Water})
		assert.True(t, card.HasFamily(family.CyberLord))
	})

	t.Run("defender chooses their own shields instead of the attacker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// The defender controls Meloppe; the attacker is an unrelated vanilla creature.
		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeVanillaAttackerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the attacker, to be prompted for shields")
		assert.False(t, action.Cancellable, "a chooser Meloppe installed can't call off the attacker's attack")
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)
		assert.Len(t, action.Cards, len(shields))
		assert.Contains(t, action.Text, "Meloppe", "the prompt should explain why the defender is the one choosing")

		// The attacker never received a shield-selection prompt of their own,
		// only the "waiting for your opponent" notice.
		attackerHeaders, err := scn.MessageHeaders(player, attackerStart)
		require.NoError(t, err)
		assert.NotContains(t, attackerHeaders, "action")

		chosen := shields[0]
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, chosen.Zone)
	})

	t.Run("attacking with Meloppe itself still leaves the choice to its opponent", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		meloppe := putCardInBattlezone(t, scn, player.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, meloppe.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender to be prompted even though they don't control Meloppe")
		assert.False(t, action.Cancellable)
		assert.Contains(t, action.Text, "Meloppe")

		chosen := shields[0]
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, chosen.Zone)
	})

	t.Run("double breaker still asks the defender for two shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeDoubleBreakerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 2)
		require.True(t, attacker.HasCondition(cnd.DoubleBreaker))

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections)
		assert.Contains(t, action.Text, "Meloppe")

		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, shields[0].Zone)
		assert.Equal(t, match.HAND, shields[1].Zone)
	})

	t.Run("Aqua Master's face-up choice goes to the defender instead of the attacker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		aquaMaster := putCardInBattlezone(t, scn, player.Player, meloppeAquaMasterUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 2)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, aquaMaster.ID))

		// Aqua Master's ShieldsSelectionEffect condition makes it confirm the
		// attack up front, before any shield is involved.
		confirm, err := scn.WaitForAction(player, attackerStart)
		require.NoError(t, err)
		assert.Equal(t, "question", confirm.ActionType)
		require.NoError(t, scn.SubmitAction(player))

		faceUpPrompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the attacker, to choose which shield gets turned face up")
		assert.Contains(t, faceUpPrompt.Text, "turn it face up")
		assert.Contains(t, faceUpPrompt.Text, "Meloppe")
		assert.False(t, faceUpPrompt.Cancellable)
		assert.Equal(t, 1, faceUpPrompt.MinSelections)
		assert.Equal(t, 1, faceUpPrompt.MaxSelections)

		// The attacker never received a prompt of their own for this choice.
		attackerHeaders, err := scn.MessageHeaders(player, attackerStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(attackerHeaders, "action"), "the attacker should only have seen their own attack confirmation")

		chosenToFlip := shields[0]
		midStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, chosenToFlip.ID))

		breakPrompt, err := scn.WaitForAction(opponent, midStart)
		require.NoError(t, err, "expected the defender to also be asked which shield to break")
		assert.Contains(t, breakPrompt.Text, "Select 1 shield")
		assert.Contains(t, breakPrompt.Text, "Meloppe")

		// Flipping it face up makes chosenToFlip visible from now on, but the
		// attacker still wasn't there for the defender's pick, so the chat
		// message is what ties "shield #1" back to what just happened. Waiting
		// for breakPrompt above proves this already arrived.
		chats, err := scn.ChatMessages(opponent, midStart)
		require.NoError(t, err)
		require.NotEmpty(t, chats)
		assert.Contains(t, chats[0], "shield #1")
		assert.Contains(t, chats[0], chosenToFlip.Name)

		chosenToBreak := shields[1]
		require.NoError(t, scn.SubmitAction(opponent, chosenToBreak.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, chosenToFlip.ShieldFaceUp, "the shield the defender chose should be the one turned face up")
		assert.Equal(t, match.HAND, chosenToBreak.Zone)
	})

	t.Run("Meloppe leaving play beforehand restores the attacker as chooser", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		meloppe := putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeVanillaAttackerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		_, err := opponent.Player.MoveCard(meloppe.ID, match.BATTLEZONE, match.GRAVEYARD, meloppeSetupSrc)
		require.NoError(t, err)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(player, attackerStart)
		require.NoError(t, err, "expected the attacker to be prompted again once Meloppe left play")
		assert.True(t, action.Cancellable)
		assert.NotContains(t, action.Text, "Meloppe", "there is nothing to explain once Meloppe is gone")

		require.NoError(t, scn.SubmitAction(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, shields[0].Zone)
	})

	t.Run("Cosmic Darts: the caster chooses their own shield instead of the opponent", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Meloppe belongs to the caster here, exercising the clause where an
		// opponent's effect would choose one of Meloppe's controller's own
		// shields.
		putCardInBattlezone(t, scn, player.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		playerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		castSpell(t, scn, player, meloppeCosmicDartsUID)

		// The opponent is the printed default chooser and never got a prompt
		// of their own.
		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "action")

		prompt, err := scn.LatestAction(player, playerStart)
		require.NoError(t, err, "expected the caster, not the opponent, to choose which of their own shields is revealed")
		assert.False(t, prompt.Cancellable)
		assert.Equal(t, 1, prompt.MinSelections)
		assert.Equal(t, 1, prompt.MaxSelections)
		assert.Len(t, prompt.Cards, len(shields))
		assert.Contains(t, prompt.Text, "Meloppe")

		require.NoError(t, scn.SubmitAction(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// A non-spell shield is only shown, then put back; the choice having
		// gone to the right player is what this test is checking.
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)
	})

	t.Run("Invincible Cataclysm: the defender chooses which of their own shields are destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 3)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		castSpell(t, scn, player, meloppeInvincibleCataclysmUID)

		prompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the caster, to choose which of their own shields are destroyed")
		assert.True(t, prompt.Cancellable)
		assert.Equal(t, 1, prompt.MinSelections)
		assert.Equal(t, 3, prompt.MaxSelections)
		assert.Contains(t, prompt.Text, "Meloppe")

		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, shields[0].Zone)
		assert.Equal(t, match.GRAVEYARD, shields[1].Zone)
	})

	t.Run("Recon Operation: the defender chooses which of their own shields are shown", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 3)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		castSpell(t, scn, player, meloppeReconOperationUID)

		prompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender to choose which of their own shields are looked at")
		assert.True(t, prompt.Cancellable)
		assert.Contains(t, prompt.Text, "Meloppe")

		revealStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// Recon Operation only looks; nothing moves.
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)

		// The caster never saw the defender choose, so the pop-up and chat
		// message are the only way to know which shield #1 it was.
		popups, err := scn.ShowCardsMessages(player, revealStart)
		require.NoError(t, err)
		require.Len(t, popups, 1)
		assert.Contains(t, popups[0], "shield #1")

		chats, err := scn.ChatMessages(player, revealStart)
		require.NoError(t, err)
		require.NotEmpty(t, chats)
		assert.Contains(t, chats[0], "shield #1")
	})

	t.Run("Gajirabute: the defender chooses which of their own shields goes to the graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, meloppeGajirabuteUID)

		prompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the summoner, to choose which of their own shields is discarded")
		assert.False(t, prompt.Cancellable)
		assert.Contains(t, prompt.Text, "Meloppe")

		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, shields[0].Zone)
	})

	t.Run("Pointa: the defender chooses which of their own shields is shown", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, meloppePointaUID)

		prompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender to choose which of their own shields is looked at")
		assert.False(t, prompt.Cancellable)
		assert.Contains(t, prompt.Text, "Meloppe")

		revealStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// Pointa only looks; nothing moves.
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)

		// The summoner never saw the defender choose, so the pop-up is the only
		// way to know which shield #1 it was.
		popups, err := scn.ShowCardsMessages(player, revealStart)
		require.NoError(t, err)
		require.Len(t, popups, 1)
		assert.Contains(t, popups[0], "shield #1")

		chats, err := scn.ChatMessages(player, revealStart)
		require.NoError(t, err)
		require.NotEmpty(t, chats)
		assert.Contains(t, chats[0], "shield #1")
	})

	t.Run("Stinger Ball: the defender chooses which of their own shields is looked at on attack", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeStingerBallUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		// Stinger Ball also carries ShieldsSelectionEffect, so the attack
		// confirms up front before its "may look" ability fires.
		confirm, err := scn.WaitForAction(player, attackerStart)
		require.NoError(t, err)
		assert.Equal(t, "question", confirm.ActionType)
		require.NoError(t, scn.SubmitAction(player))

		lookPrompt, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the attacker, to choose which of their own shields is looked at")
		assert.True(t, lookPrompt.Cancellable)
		assert.Contains(t, lookPrompt.Text, "Meloppe")

		midStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		playerMidStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID))

		// The attacker is shown the revealed shield through a non-dismissible
		// prompt of their own; it has to be closed before the attack sequence
		// can continue to the shield-break selection.
		require.NoError(t, scn.WaitForMessage(player, playerMidStart, "show_cards_non_dismissible"))

		// The attacker never saw the defender choose, so the pop-up and chat
		// message are the only way to know which shield #1 it was.
		popups, err := scn.ShowCardsMessages(player, playerMidStart)
		require.NoError(t, err)
		require.Len(t, popups, 1)
		assert.Contains(t, popups[0], "shield #1")

		chats, err := scn.ChatMessages(player, playerMidStart)
		require.NoError(t, err)
		require.NotEmpty(t, chats)
		assert.Contains(t, chats[0], "shield #1")

		require.NoError(t, scn.CancelAction(player))

		// The mandatory shield-break selection follows right after.
		breakPrompt, err := scn.WaitForAction(opponent, midStart)
		require.NoError(t, err)
		assert.Contains(t, breakPrompt.Text, "Select 1 shield")
		assert.Contains(t, breakPrompt.Text, "Meloppe")

		require.NoError(t, scn.SubmitAction(opponent, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// Stinger Ball's look effect only shows the shield; the break is what moves it.
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)
		assert.Equal(t, match.HAND, shields[1].Zone)
	})

	t.Run("Adomis: a blind pick landing on the opponent's zone lets them choose which shield instead", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		adomis := putCardInBattlezone(t, scn, player.Player, meloppeAdomisUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		playerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionUseTapAbility(player, adomis.ID))
		_, err = scn.WaitForMultipartAction(player, playerStart)
		require.NoError(t, err)

		// The picker's blind pick names a specific shield of the opponent's,
		// even though the client only sees it face down.
		require.NoError(t, scn.SubmitAction(player, shields[0].ID))

		prompt, err := scn.WaitForAction(opponent, opponentStart)
		require.NoError(t, err, "expected the defender, not the picker, to choose which of their own shields is looked at")
		assert.False(t, prompt.Cancellable, "the picker already committed to affecting one of the defender's shields")
		assert.Equal(t, 1, prompt.MinSelections)
		assert.Equal(t, 1, prompt.MaxSelections)
		assert.Contains(t, prompt.Text, "Meloppe")

		// The defender's own choice (shields[1], their shield #2) is what
		// actually gets revealed, not the picker's earlier blind pick.
		revealStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// Adomis only looks; nothing moves either way.
		assert.Equal(t, match.SHIELDZONE, shields[0].Zone)
		assert.Equal(t, match.SHIELDZONE, shields[1].Zone)

		// The picker (player) never saw the defender choose, so the pop-up and
		// chat message are the only way to know which shield #2 turned out to
		// be, and whose it was.
		popups, err := scn.ShowCardsMessages(player, revealStart)
		require.NoError(t, err)
		require.Len(t, popups, 1)
		assert.Contains(t, popups[0], "shield #2")
		assert.Contains(t, popups[0], shields[1].Name)

		chats, err := scn.ChatMessages(player, revealStart)
		require.NoError(t, err)
		require.Len(t, chats, 1)
		assert.Contains(t, chats[0], "shield #2")
		assert.Contains(t, chats[0], shields[1].Name)
		assert.Contains(t, chats[0], opponent.Player.Username())
		assert.Contains(t, chats[0], player.Player.Username())
	})

	t.Run("Ularus: a blind pick landing on the opponent's zone lets them choose which shield is turned face up", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		playerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		// Summoned alone, so its allowance (one shield per own creature in the
		// battle zone, itself included) is exactly 1.
		summonWithOwnMana(t, scn, player, meloppeUlarusUID)
		_, err = scn.WaitForMultipartAction(player, playerStart)
		require.NoError(t, err)

		require.NoError(t, scn.SubmitAction(player, shields[0].ID))

		prompt, err := scn.WaitForAction(opponent, opponentStart)
		require.NoError(t, err, "expected the defender, not the summoner, to choose which of their own shields is turned face up")
		assert.False(t, prompt.Cancellable)
		assert.Equal(t, 1, prompt.MinSelections)
		assert.Equal(t, 1, prompt.MaxSelections)
		assert.Contains(t, prompt.Text, "Meloppe")

		revealStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, shields[0].ShieldFaceUp, "the picker's blind pick should have been discarded, not turned face up")
		assert.True(t, shields[1].ShieldFaceUp, "the defender's own choice should be the one turned face up")

		// Flipping it face up makes shields[1] visible from now on, but the
		// summoner still wasn't there for the defender's pick, so the chat
		// message is what ties "shield #2" back to what just happened.
		chats, err := scn.ChatMessages(player, revealStart)
		require.NoError(t, err)
		require.NotEmpty(t, chats)
		assert.Contains(t, chats[0], "shield #2")
		assert.Contains(t, chats[0], shields[1].Name)
	})
}
