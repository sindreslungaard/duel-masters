package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {

	// Clear conditions
	if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
		card.ClearConditions()
	}

	// Add spell condition to this card
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Spell, nil, card.ID)
	}

	// When the spell is played from hand
	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// make sure we haven't attacked yet
		if _, ok := ctx.Match.Step.(*match.AttackStep); ok {
			ctx.Match.WarnPlayer(card.Player, "You can't cast spells after attacking or using tap ability.")
			ctx.InterruptFlow()
			return
		}

		ctx.ScheduleAfter(func() {

			manazone, err := card.Player.Container(match.MANAZONE)

			if err != nil {
				return
			}

			untappedMana := make([]*match.Card, 0)
			for _, c := range manazone {
				if !c.Tapped {
					untappedMana = append(untappedMana, c)
				}
			}

			if !card.Player.CanPlayCard(card, untappedMana) {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("You do not have sufficient mana to play %s", card.Name))
				ctx.InterruptFlow()
				return
			}

			manaCost := card.EffectiveManaCost()

			ctx.Match.NewAction(
				card.Player,
				untappedMana,
				manaCost,
				manaCost,
				fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least %v civilization card(s).", manaCost, card.Name, ManaRequirementText(card)),
				true,
			)

			for {

				action := card.Player.NextAction()

				if action.Cancel {
					ctx.Match.CloseAction(card.Player)
					ctx.InterruptFlow()
					break
				}

				cards := make([]*match.Card, 0)

				for _, id := range action.Cards {
					mana, err := card.Player.GetCard(id, match.MANAZONE)

					if err != nil {
						continue
					}

					cards = append(cards, mana)
				}

				if len(action.Cards) != manaCost || !match.AssertCardsIn(untappedMana, action.Cards...) || !card.Player.CanPlayCard(card, cards) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				ctx.Match.CloseAction(card.Player)

				for _, mana := range cards {
					card.Player.TapCard(mana)
				}

				ctx.Match.CastSpell(card, false)

				break

			}

		})

	}

	// On spell cast
	if event, ok := ctx.Event.(*match.SpellCast); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// Once a spell is cast it is no longer in hand, by rule, even while its
		// own text is still resolving: something that spell's effect puts into
		// play can trigger another card's own ability re-entrantly (a creature's
		// "put into the battle zone" trigger, for instance), and that nested
		// ability must see this card already sitting in the graveyard rather
		// than mid-cast in hand. Move it synchronously, before the card's own
		// effect body (registered after fx.Spell in the same c.Use chain) runs.
		// A card whose text sends it elsewhere instead of the graveyard
		// (fx.Charger, etc.) relocates it again once SpellResolved fires below.
		card.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, card.ID)

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s casted the spell %s", card.Player.Username(), card.Name))

		ctx.ScheduleAfter(func() {
			e := &match.SpellResolved{
				CardID:        event.CardID,
				FromShield:    event.FromShield,
				MatchPlayerID: event.MatchPlayerID,
			}

			ctx.Match.HandleFx(match.NewContext(ctx.Match, e))
		})

	}

}
