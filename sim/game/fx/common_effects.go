package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
	"strings"
)

func GiveTapAbilityToAllies(card *match.Card, ctx *match.Context, alliesFilter func(card *match.Card) bool, tapAbility func(card *match.Card, ctx *match.Context)) {
	// This is added for the case where the card is added to the field. There is another creature
	// that doesn't initially have a tap abbility but should receive one. The change doesn't propagate fast
	// enough to the FE and that creature doesn't get tap ability until another action takes places.
	// This is an ugly workaround.
	FindFilter(
		card.Player,
		match.BATTLEZONE,
		alliesFilter,
	).Map(func(x *match.Card) {
		x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
	})

	ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
		if card.Zone != match.BATTLEZONE {
			Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			exit()
			return
		}

		FindFilter(
			card.Player,
			match.BATTLEZONE,
			alliesFilter,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
		})
	})
}

func FilterShieldTriggers(ctx *match.Context, filter func(*match.Card) bool) {

	if event, ok := ctx.Event.(*match.ShieldTriggerEvent); ok {
		validCards, invalidCards := FilterCardList(event.Cards, filter)
		event.Cards = validCards
		event.UnplayableCards = append(event.UnplayableCards, invalidCards...)
	}

}

// DiscardOwnXCards makes the card's controller choose x cards from their own
// hand and discard them. A hand holding fewer than x simply loses what is
// there, because the selection is clamped to the cards that exist, and an empty
// hand opens no prompt at all.
func DiscardOwnXCards(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: Choose %d card(s) from your hand to discard.", card.Name, x),
			x,
			x,
			false,
		).Map(func(discarded *match.Card) {
			moved, err := card.Player.MoveCard(discarded.ID, match.HAND, match.GRAVEYARD, card.ID)

			if err != nil || moved.Zone != match.GRAVEYARD {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was discarded from %s's hand by %s", discarded.Name, card.Player.Username(), card.Name))
		})
	}
}

func PlayerDiscardsRandomCard(card *match.Card, ctx *match.Context) {

	hand, err := card.Player.Container(match.HAND)

	if err != nil || len(hand) < 1 {
		return
	}

	discardedCard, err := card.Player.MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD, card.ID)
	if err == nil {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
	}

}

func OpponentDiscardsRandomCard(card *match.Card, ctx *match.Context) {

	hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

	if err != nil || len(hand) < 1 {
		return
	}

	discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD, card.ID)
	if err == nil {
		ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
	}

}

func OpponentDiscards2RandomCards(card *match.Card, ctx *match.Context) {
	OpponentDiscardsRandomCard(card, ctx)
	OpponentDiscardsRandomCard(card, ctx)
}

// To be used as part of a card effect, not for initial shuffle
func ShuffleDeck(card *match.Card, ctx *match.Context, forOpponent bool) {
	if !forOpponent {
		card.Player.ShuffleDeck()
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s shuffled their deck", card.Player.Username()))
	} else {
		opponent := ctx.Match.Opponent(card.Player)
		opponent.ShuffleDeck()
		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s deck shuffled by %s effect", opponent.Username(), card.Name))
	}

}

// SwapHandAndMana simultaneously exchanges the cards that were in player's
// hand and mana zone when the effect began. Cards entering the mana zone are
// tapped. If source is still in hand while resolving, it is not moved.
func SwapHandAndMana(source *match.Card, player *match.Player) {
	manaCards := Find(player, match.MANAZONE)
	handCards := FindFilter(player, match.HAND, func(card *match.Card) bool {
		return card.ID != source.ID
	})

	for _, manaCard := range manaCards {
		player.MoveCard(manaCard.ID, match.MANAZONE, match.HAND, source.ID)
	}

	for _, handCard := range handCards {
		moved, err := player.MoveCard(handCard.ID, match.HAND, match.MANAZONE, source.ID)
		if err == nil && moved.Zone == match.MANAZONE {
			player.TapCard(moved)
		}
	}
}

func BlockerWhenNoShields(card *match.Card, ctx *match.Context) {
	condition := &match.Condition{ID: cnd.Blocker, Val: true, Src: card.ID}
	HaveSelfConditionsWhenNoShields(card, ctx, []*match.Condition{condition})
}

func HaveSelfConditionsWhenNoShields(card *match.Card, ctx *match.Context, conditions []*match.Condition) {
	ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
		notInTheBZ := card.Zone != match.BATTLEZONE
		if notInTheBZ || IHaveShields(card) {
			for _, cond := range conditions {
				card.RemoveSpecificConditionBySource(cond.ID, card.ID)
			}
		}

		if notInTheBZ {
			exit()
			return
		}

		if IDontHaveShields(card, ctx2) {
			for _, cond := range conditions {
				if cond.ID == cnd.Blocker {
					ForceBlocker(card, ctx2, card.ID)
				} else {
					card.AddUniqueSourceCondition(cond.ID, cond.Val, card.ID)
				}
			}
		}
	})
}

// ResolveShieldChooser returns who should be prompted to choose a face-down
// card out of shieldOwner's shield zone, honoring a persistent effect (e.g.
// Meloppe) that redirects a cross-player shield choice back to the shield's
// own owner.
//
// Choosing your own shields is never redirected, so callers that already pass
// the shield owner as the chooser can skip this and call SelectBackside
// directly; there is nothing for a persistent effect to reverse there.
func ResolveShieldChooser(ctx *match.Context, defaultChooser *match.Player, shieldOwner *match.Player) *match.Player {
	if defaultChooser == shieldOwner {
		return defaultChooser
	}

	event := &match.ChooseShieldEvent{ShieldOwner: shieldOwner, Chooser: defaultChooser}
	ctx.Match.HandleFx(match.NewContext(ctx.Match, event))

	return event.Chooser
}

// ResolveShieldChoices re-delegates the picks of a combined, owner-agnostic
// shield selection (an effect naming no specific owner, such as "choose a
// shield") already made blind, for any pick that landed on someone other
// than defaultChooser and that a persistent effect (e.g. Meloppe) says its
// owner should have chosen instead.
//
// defaultChooser's own picks are returned unchanged, since there is nothing
// to reverse when the chooser already owns the shield. Everything else is
// grouped by owner; an owner Meloppe redirects to themselves gets a fresh,
// mandatory backside pick of the same count from their own shield zone,
// replacing what the original chooser blindly landed on. Nothing about that
// original pick was ever revealed to the chooser, so swapping it out here
// leaks no information.
func ResolveShieldChoices(ctx *match.Context, defaultChooser *match.Player, picks CardCollection) CardCollection {
	result := make(CardCollection, 0, len(picks))

	byOwner := make(map[*match.Player][]*match.Card)
	for _, card := range picks {
		if card.Player == defaultChooser {
			result = append(result, card)
			continue
		}

		byOwner[card.Player] = append(byOwner[card.Player], card)
	}

	for owner, owned := range byOwner {
		chooser := ResolveShieldChooser(ctx, defaultChooser, owner)
		if chooser == defaultChooser {
			result = append(result, owned...)
			continue
		}

		replaced := SelectBackside(
			chooser,
			ctx.Match,
			owner,
			match.SHIELDZONE,
			fmt.Sprintf("Meloppe's effect: Choose %d of your shields instead.", len(owned)),
			len(owned),
			len(owned),
			false,
		)

		result = append(result, replaced...)
	}

	return result
}

// ShieldPossessive phrases a prompt about shieldOwner's shield zone from
// chooser's point of view, so a chooser Meloppe redirected onto their own
// shields reads "your shields" instead of a stale "your opponent's shields"
// that assumed the default chooser.
func ShieldPossessive(chooser *match.Player, shieldOwner *match.Player) string {
	if chooser == shieldOwner {
		return "your"
	}

	return "your opponent's"
}

// MeloppeNote returns a parenthetical to append to a shield-choice prompt
// when chooser isn't the effect's usual chooser, so a player unfamiliar with
// Meloppe isn't left wondering why they're the one being asked. Hardcoded to
// name Meloppe: it is currently the only effect that reverses who chooses a
// shield, so there is nothing more general to say here.
func MeloppeNote(chooser *match.Player, defaultChooser *match.Player) string {
	if chooser == defaultChooser {
		return ""
	}

	return " (Meloppe's effect: you choose instead of your opponent.)"
}

// ShieldNumber returns a shield's 1-based position in its owner's shield
// zone, the same order the owner sees their own face-down shields numbered
// in. Call it before the shield moves out of that zone: a card that has
// already broken, been discarded, or otherwise left the shield zone can no
// longer be found there, and this returns 0.
//
// A player shown a shield's identity without having picked it themselves
// (e.g. Meloppe redirected who chose) has no other way to know which one it
// was, since nobody but the shield's owner ever sees that pick happen.
func ShieldNumber(shield *match.Card) int {
	shields, err := shield.Player.Container(match.SHIELDZONE)
	if err != nil {
		return 0
	}

	for i, s := range shields {
		if s.ID == shield.ID {
			return i + 1
		}
	}

	return 0
}

// DescribeShield labels a shield with its owner-relative position and name,
// e.g. "shield #2 (Aqua Surfer)", for a message about a shield shown to
// someone other than its owner. Call it before the shield moves out of its
// shield zone; see ShieldNumber.
func DescribeShield(shield *match.Card) string {
	return fmt.Sprintf("shield #%d (%s)", ShieldNumber(shield), shield.Name)
}

func RotateShields(card *match.Card, ctx *match.Context, max int) {

	nrShields, err := card.Player.Container(match.SHIELDZONE)
	if err != nil {
		return
	}

	if len(nrShields) < 1 {
		return
	}

	toShield := Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.HAND,
		fmt.Sprintf("%s: You may select up to %d card(s) from your hand and put it into the shield zone", card.Name, max),
		1,
		max,
		true,
	)

	cardsMoved := len(toShield)
	if cardsMoved < 1 {
		return
	}

	for _, c := range toShield {
		c.Player.MoveCard(c.ID, match.HAND, match.SHIELDZONE, card.ID)
	}

	toHand := SelectBackside(
		card.Player,
		ctx.Match,
		card.Player,
		match.SHIELDZONE,
		fmt.Sprintf("%s: Select %d of your shields that will be moved to your hand", card.Name, cardsMoved),
		cardsMoved,
		cardsMoved,
		false,
	)

	for _, c := range toHand {
		c.Player.MoveCard(c.ID, match.SHIELDZONE, match.HAND, card.ID)
	}

}

func DestroyOpShield(card *match.Card, ctx *match.Context) {
	BreakXOpShields(1)(card, ctx)
}

// BreakXOpShields breaks x of the opponent's shields, chosen face down by the
// card's controller. Breaking is not attacking: the shields go to the
// opponent's hand and their shield triggers are offered, but nothing blocks and
// no creature is tapped.
//
// An opponent holding fewer than x shields simply loses the ones they have,
// because the selection is clamped to what is in the shield zone.
func BreakXOpShields(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)
		chooser := ResolveShieldChooser(ctx, card.Player, opponent)

		shields := SelectBackside(
			chooser,
			ctx.Match,
			opponent,
			match.SHIELDZONE,
			fmt.Sprintf("%s effect: select %d shield(s) to break%s", card.Name, x, MeloppeNote(chooser, card.Player)),
			x,
			x,
			false,
		)

		if len(shields) < 1 {
			return
		}

		ctx.Match.BreakShields(shields, card)

		ctx.Match.ReportActionInChat(opponent,
			fmt.Sprintf("%s effect broke %d of %s's shields", card.Name, len(shields), opponent.Username()))
	}
}

// PutOpShieldIntoGraveyard lets the card's controller choose one of the
// opponent's shields face down and puts it straight into the opponent's
// graveyard. Unlike DestroyOpShield the shield is never broken, so it does not
// go to the opponent's hand and its shield trigger is not offered.
func PutOpShieldIntoGraveyard(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)
	chooser := ResolveShieldChooser(ctx, card.Player, opponent)

	SelectBackside(
		chooser,
		ctx.Match,
		opponent,
		match.SHIELDZONE,
		fmt.Sprintf("%s's effect: Choose one of %s shields and put it into their graveyard.%s", card.Name, ShieldPossessive(chooser, opponent), MeloppeNote(chooser, card.Player)),
		1,
		1,
		false,
	).Map(func(shield *match.Card) {
		moved, err := opponent.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
		if err == nil && moved.Zone == match.GRAVEYARD {
			ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("A shield was put into %s's graveyard by %s", opponent.Username(), card.Name))
		}
	})
}

func OpDiscardsXCards(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {

		min := 0
		handCount := ctx.Match.Opponent(card.Player).Denormalized().HandCount

		if x > handCount {
			min = handCount
		} else {
			min = x
		}

		Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			fmt.Sprintf("%s: Select %d card(s) from your hand that will be sent to your graveyard", card.Name, x),
			min,
			x,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved from %s's hand to his graveyard by %s", x.Name, x.Player.Username(), card.Name))
		})
	}
}

// ShowXShields lets the card's controller look at up to x of the opponent's
// shields. Pass cancellable for a printed "you may look at", and false when the
// card makes looking mandatory. The shields stay where they are.
func ShowXShields(x int, cancellable bool) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {

		shieldsID := []string{}
		descriptions := []string{}
		opponent := ctx.Match.Opponent(card.Player)
		chooser := ResolveShieldChooser(ctx, card.Player, opponent)

		SelectBackside(
			chooser,
			ctx.Match,
			opponent,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select %d of %s shields that will be shown to you%s", card.Name, x, ShieldPossessive(chooser, opponent), MeloppeNote(chooser, card.Player)),
			1,
			x,
			cancellable,
		).Map(func(shield *match.Card) {
			// Grabbed before the pop-up so the numbers still match: whoever
			// picked shield (its own owner, if Meloppe redirected) never sees
			// this reveal, so it's the only way card.Player can tell them which
			// shields they were.
			shieldsID = append(shieldsID, shield.ImageID)
			descriptions = append(descriptions, DescribeShield(shield))
		})

		// Nothing was chosen, either because the opponent has no shields or
		// because an optional look was declined.
		if len(shieldsID) < 1 {
			return
		}

		ctx.Match.ShowCards(
			card.Player,
			fmt.Sprintf("Your opponent's %s:", strings.Join(descriptions, ", ")),
			shieldsID,
		)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s of %s was shown to %s", strings.Join(descriptions, ", "), opponent.Username(), card.Player.Username()))
	}

}

func OpponentDiscardsHand(card *match.Card, ctx *match.Context) {

	Find(
		ctx.Match.Opponent(card.Player),
		match.HAND,
	).Map(func(x *match.Card) {
		ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
	})

	ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s's hand was discarded by %s", ctx.Match.Opponent(card.Player).Username(), card.Name))

}

// OpponentDiscardsXRandomCards makes the opponent discard x cards at random. A
// hand holding fewer simply loses what is there.
func OpponentDiscardsXRandomCards(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		for range x {
			OpponentDiscardsRandomCard(card, ctx)
		}
	}
}
