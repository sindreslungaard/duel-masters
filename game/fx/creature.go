package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	// Untap the card, add creature condition
	if _, ok := ctx.Event.(*match.UntapStep); ok {

		if ctx.Match.IsPlayerTurn(card.Player) {
			card.AddCondition(cnd.Creature, nil, nil)
			card.Tapped = false
		}

	}

	// Add to battlezone
	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// Do this last in case any other cards want to interrupt the flow
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

			ctx.Match.NewAction(
				card.Player,
				untappedMana,
				card.ManaCost,
				card.ManaCost,
				fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", card.ManaCost, card.Name, card.Civ),
				true,
			)

			for {

				action := <-card.Player.Action

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

				if len(action.Cards) != card.ManaCost || !match.AssertCardsIn(untappedMana, action.Cards...) || !card.Player.CanPlayCard(card, cards) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				ctx.Match.CloseAction(card.Player)

				for _, mana := range cards {
					mana.Tapped = true
				}

				card.AddCondition(cnd.SummoningSickness, nil, nil)

				card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)

				break

			}

		})

	}

	// Attack the player
	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			opponent := ctx.Match.Opponent(card.Player)

			shieldzone, err := opponent.Container(match.SHIELDZONE)

			if err != nil {
				return
			}

			shieldsAttacked := make([]*match.Card, 0)

			if len(shieldzone) > 0 {

				minmax := 1

				if card.HasCondition(cnd.DoubleBreaker) {
					minmax = 2
				}

				if card.HasCondition(cnd.TripleBreaker) {
					minmax = 3
				}

				if minmax > len(shieldzone) {
					minmax = len(shieldzone)
				}

				ctx.Match.NewBacksideAction(card.Player, shieldzone, minmax, minmax, fmt.Sprintf("Select %v shield(s) to break", minmax), true)

				for {

					action := <-card.Player.Action

					if action.Cancel {
						ctx.Match.CloseAction(card.Player)
						return
					}

					if len(action.Cards) != minmax || !match.AssertCardsIn(shieldzone, action.Cards[0]) {
						ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
						continue
					}

					for _, cardID := range action.Cards {
						shield, err := opponent.GetCard(cardID, match.SHIELDZONE)
						if err != nil {
							logrus.Debug("Could not find specified shield in shieldzone")
							continue
						}
						shieldsAttacked = append(shieldsAttacked, shield)
					}

					ctx.Match.CloseAction(card.Player)

					break

				}

			}

			card.Tapped = true

			// Allow the opponent to block if they can
			if len(event.Blockers) > 0 {

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

				identifierStr := "you"

				if len(shieldsAttacked) > 0 {
					identifierStr = fmt.Sprintf("%v of your shields", len(shieldsAttacked))
				}

				ctx.Match.NewAction(opponent, event.Blockers, 1, 1, fmt.Sprintf("%s (%v) is attacking %s. Choose a creature to block the attack with or close to not block the attack.", card.Name, ctx.Match.GetPower(card, true), identifierStr), true)

				for {

					action := <-opponent.Action

					if action.Cancel {
						ctx.Match.EndWait(card.Player)
						ctx.Match.CloseAction(opponent)

						if len(shieldzone) < 1 {
							// Win
							ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
						} else {
							// Break n shields
							ctx.Match.BreakShields(shieldsAttacked)
						}

						break
					}

					if len(action.Cards) != 1 || !match.AssertCardsIn(event.Blockers, action.Cards[0]) {
						ctx.Match.ActionWarning(opponent, "Your selection of cards does not fulfill the requirements")
						continue
					}

					c, err := opponent.GetCard(action.Cards[0], match.BATTLEZONE)

					if err != nil {
						ctx.Match.ActionWarning(opponent, "The card you selected is not in the battlefield")
						continue
					}

					c.Tapped = true

					ctx.Match.EndWait(card.Player)
					ctx.Match.CloseAction(opponent)

					ctx.Match.Battle(card, c, true)

					break

				}

			} else {

				card.Tapped = true

				if len(shieldzone) < 1 {
					// Win
					ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				} else {
					// Break n shields
					ctx.Match.BreakShields(shieldsAttacked)
				}

			}

		})

	}

	// Attack a creature
	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			opponent := ctx.Match.Opponent(card.Player)

			battlezone, err := opponent.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			attackable := make([]*match.Card, 0)

			for _, c := range battlezone {
				if c.Tapped || card.HasCondition(cnd.AttackUntapped) {
					attackable = append(attackable, c)
				}
			}

			if len(attackable) < 1 {
				ctx.Match.WarnPlayer(card.Player, "None of your opponents creatures can currently be attacked.")
				return
			}

			attackedCreatures := make([]*match.Card, 0)

			ctx.Match.NewAction(card.Player, attackable, 1, 1, "Select the creature to attack", true)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					ctx.Match.CloseAction(card.Player)
					return
				}

				if len(action.Cards) != 1 || !match.AssertCardsIn(attackable, action.Cards[0]) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				c, err := opponent.GetCard(action.Cards[0], match.BATTLEZONE)

				if err != nil {
					return
				}

				attackedCreatures = append(attackedCreatures, c)

				ctx.Match.CloseAction(card.Player)

				break

			}

			if len(attackedCreatures) < 1 {
				return
			}

			c := attackedCreatures[0]

			card.Tapped = true

			// Allow the opponent to block if they can
			if len(event.Blockers) > 0 {

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

				ctx.Match.NewAction(opponent, event.Blockers, 1, 1, fmt.Sprintf("%s (%v) is attacking %s (%v). Choose a creature to block the attack with or close to not block the attack.", card.Name, ctx.Match.GetPower(card, true), c.Name, ctx.Match.GetPower(c, false)), true)

				for {

					action := <-opponent.Action

					if action.Cancel {
						ctx.Match.EndWait(card.Player)
						ctx.Match.CloseAction(opponent)
						break
					}

					if len(action.Cards) != 1 || !match.AssertCardsIn(event.Blockers, action.Cards[0]) {
						ctx.Match.ActionWarning(opponent, "Your selection of cards does not fulfill the requirements")
						continue
					}

					blocker, err := opponent.GetCard(action.Cards[0], match.BATTLEZONE)

					if err != nil {
						ctx.Match.ActionWarning(opponent, "The card you selected is not in the battlefield")
						continue
					}

					blocker.Tapped = true

					ctx.Match.EndWait(card.Player)
					ctx.Match.CloseAction(opponent)

					ctx.Match.Battle(card, blocker, true)

					return

				}

			}

			ctx.Match.Battle(card, c, false)

		})

	}

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.ScheduleAfter(func() {

				card.Player.MoveCard(card.ID, match.BATTLEZONE, match.GRAVEYARD)

				// Slayer
				if card.HasCondition(cnd.Slayer) {

					creature, err := ctx.Match.Opponent(card.Player).GetCard(event.Source.ID, match.BATTLEZONE)

					if err == nil {

						ctx.Match.Destroy(creature, card)

					}

				}

			})

		}

	}

}
