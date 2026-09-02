package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// VortexEvolution implements "Vortex Evolution — Put on one of your <A> and
// one of your <B>."
//
// Unlike a plain evolution, which needs a single matching creature already in
// the battle zone, a vortex evolution needs two distinct creatures: one
// matching aFilter and a different one matching bFilter (the same creature
// can never fill both roles, even if its own characteristics satisfy both
// filters). aDescription/bDescription only name the two requirements in
// prompts, e.g. "Light Bringer" and "Cyber Lord".
//
// The player chooses which of the two goes on top of the other, and the
// vortex evolution creature is placed on top of that stack: it takes the tap
// state of whichever creature it lands on directly, exactly like a regular
// evolution.
//
// It does carry cnd.Evolution, like a plain evolution creature, purely so
// generic "is this an evolution creature" checks (for example, a card
// revealed off the top of a deck) recognize it - matching either requirement
// name is a much looser test than the two distinct creatures actually needed,
// so fx.CanBeSummoned may offer a vortex evolution creature it cannot
// actually complete; the guards below still make an illegal summon
// impossible either way.
func VortexEvolution(aDescription string, aFilter func(*match.Card) bool, bDescription string, bFilter func(*match.Card) bool) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddCondition(cnd.Evolution, []string{aDescription, bDescription}, card.ID)
		}

		handleVortexEvolutionEvents(card, ctx, aDescription, aFilter, bDescription, bFilter)
	}
}

// vortexPartnerExists returns true if the player has a creature matching
// filter in the battle zone other than the one named by excludingID.
func vortexPartnerExists(player *match.Player, excludingID string, filter func(*match.Card) bool) bool {
	for _, c := range FindFilter(player, match.BATTLEZONE, filter) {
		if c.ID != excludingID {
			return true
		}
	}

	return false
}

// canVortexEvolve returns true if there exist two distinct creatures in the
// battle zone, one matching aFilter and a different one matching bFilter.
func canVortexEvolve(player *match.Player, aFilter func(*match.Card) bool, bFilter func(*match.Card) bool) bool {
	for _, a := range FindFilter(player, match.BATTLEZONE, aFilter) {
		if vortexPartnerExists(player, a.ID, bFilter) {
			return true
		}
	}

	return false
}

func handleVortexEvolutionEvents(card *match.Card, ctx *match.Context, aDescription string, aFilter func(*match.Card) bool, bDescription string, bFilter func(*match.Card) bool) {

	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		if event.CardID != card.ID {
			return
		}

		if !canVortexEvolve(card.Player, aFilter, bFilter) {
			ctx.InterruptFlow()
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("There are no cards for %s to evolve from in your battle zone", card.Name))
			return
		}

		ctx.ScheduleAfter(func() {
			card.RemoveCondition(cnd.SummoningSickness)
		})

	}

	if event, ok := ctx.Event.(*match.CardPlayedEvent); ok {

		if event.CardID != card.ID {
			return
		}

		// Only offered when picking it still leaves a distinct match for the
		// other requirement, so this selection can never strand the next one
		// even when a single creature happens to satisfy both filters.
		aCandidates := SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("Choose one of your %s for %s to evolve from", aDescription, card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return aFilter(x) && vortexPartnerExists(card.Player, x.ID, bFilter) },
			false,
		)

		if len(aCandidates) < 1 {
			ctx.InterruptFlow()
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("There are no cards for %s to evolve from in your battle zone", card.Name))
			return
		}

		baseA := aCandidates[0]

		bCandidates := SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("Choose one of your %s for %s to evolve from", bDescription, card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.ID != baseA.ID && bFilter(x) },
			false,
		)

		if len(bCandidates) < 1 {
			ctx.InterruptFlow()
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("There are no cards for %s to evolve from in your battle zone", card.Name))
			return
		}

		baseB := bCandidates[0]

		orderedIDs := OrderCards(
			card.Player,
			ctx.Match,
			[]*match.Card{baseA, baseB},
			fmt.Sprintf("Choose which creature to stack on top before %s is put on top of both", card.Name),
		)

		top, err := card.Player.GetCard(orderedIDs[0], match.BATTLEZONE)
		if err != nil {
			return
		}

		bottom, err := card.Player.GetCard(orderedIDs[1], match.BATTLEZONE)
		if err != nil {
			return
		}

		card.ClearAttachments()
		stagePendingEvolutionTapState(card, top.Tapped)

		card.Player.MoveCard(top.ID, match.BATTLEZONE, match.HIDDENZONE, card.ID)
		card.Player.MoveCard(bottom.ID, match.BATTLEZONE, match.HIDDENZONE, card.ID)
		card.Attach(top, bottom)

		// Announced once the pile is whole, once per base, so anything
		// generically watching for "a creature evolved" sees both - matching
		// how a plain evolution announces its single base.
		matchPlayerID := byte(2)
		if card.Player == ctx.Match.Player1.Player {
			matchPlayerID = 1
		}

		ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.EvolutionEvent{
			CardID:        card.ID,
			BaseID:        top.ID,
			MatchPlayerID: matchPlayerID,
		}))

		ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.EvolutionEvent{
			CardID:        card.ID,
			BaseID:        bottom.ID,
			MatchPlayerID: matchPlayerID,
		}))

	}

	handleEvolutionCardMoved(card, ctx)

}
