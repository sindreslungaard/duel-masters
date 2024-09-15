package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	// Clear conditions
	if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
		card.ClearConditions()
	}

	// Untap the card, add creature condition
	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Creature, nil, nil)

		if ctx.Match.IsPlayerTurn(card.Player) {
			card.Tapped = false
		}

	}

	// Add to battlezone
	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// make sure we haven't attacked yet
		if _, ok := ctx.Match.Step.(*match.AttackStep); ok {
			ctx.Match.WarnPlayer(card.Player, "You can't summon creatures after attacking")
			ctx.InterruptFlow()
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

			manaCost := card.ManaCost
			for _, condition := range card.Conditions() {
				if condition.ID == cnd.ReducedCost {
					manaCost -= condition.Val.(int)
					if manaCost < 1 {
						manaCost = 1
					}
				}

				if condition.ID == cnd.IncreasedCost {
					manaCost += condition.Val.(int)
				}
			}

			ctx.Match.NewAction(
				card.Player,
				untappedMana,
				manaCost,
				manaCost,
				fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", manaCost, card.Name, card.Civ),
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

				if len(action.Cards) != manaCost || !match.AssertCardsIn(untappedMana, action.Cards...) || !card.Player.CanPlayCard(card, cards) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				ctx.Match.CloseAction(card.Player)

				cardPlayedCtx := match.NewContext(ctx.Match, &match.CardPlayedEvent{
					CardID: card.ID,
				})
				ctx.Match.HandleFx(cardPlayedCtx)

				if !cardPlayedCtx.Cancelled() {

					for _, mana := range cards {
						mana.Tapped = true
					}

					card.AddCondition(cnd.SummoningSickness, nil, nil)

					card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, card.ID)
					ctx.Match.Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player.Username(), card.Name))

				}

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

		if card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack while tapped", card.Name))
			ctx.InterruptFlow()
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		opponent := ctx.Match.Opponent(card.Player)

		// Add blockers to the attack
		FindFilter(
			opponent,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) && !x.Tapped },
		).Map(func(x *match.Card) {
			event.Blockers = append(event.Blockers, x)
		})

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

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

				for _, condition := range card.Conditions() {
					if condition.ID != cnd.ShieldBreakModifier {
						continue
					}

					if val, ok := condition.Val.(int); ok {
						minmax += val
					}
				}

				if minmax > len(shieldzone) {
					minmax = len(shieldzone)
				}

				ctx.Match.NewBacksideAction(card.Player, shieldzone, minmax, minmax, fmt.Sprintf("Select %v shield(s) to break", minmax), true)

				for {

					action := <-card.Player.Action

					if action.Cancel {
						ctx.InterruptFlow()
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

			// Broadcast state so that opponent can see that this card is tapped if they get any shield triggers
			ctx.Match.BroadcastState()

			// Allow the opponent to block if they can
			if len(event.Blockers) > 0 && !card.HasCondition(cnd.CantBeBlocked) {

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

						ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.AttackConfirmed{CardID: card.ID, Player: true, Creature: false}))

						if len(shieldzone) < 1 {
							// Win
							ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
						} else {
							// Break n shields
							ctx.Match.BreakShields(shieldsAttacked, card.ID)
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

				ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.AttackConfirmed{CardID: card.ID, Player: true, Creature: false}))

				if len(shieldzone) < 1 {
					// Win
					ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				} else {
					// Break n shields
					ctx.Match.BreakShields(shieldsAttacked, card.ID)
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

		if card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack while tapped", card.Name))
			ctx.InterruptFlow()
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		opponent := ctx.Match.Opponent(card.Player)

		// Add blockers to the attack
		FindFilter(
			opponent,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) && !x.Tapped },
		).Map(func(x *match.Card) {
			event.Blockers = append(event.Blockers, x)
		})

		battlezone, err := opponent.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		// Add attackable creatures
		for _, c := range battlezone {
			if c.Tapped || card.HasCondition(cnd.AttackUntapped) {
				event.AttackableCreatures = append(event.AttackableCreatures, c)
			}
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			if len(event.AttackableCreatures) < 1 {
				ctx.Match.WarnPlayer(card.Player, "None of your opponents creatures can currently be attacked.")
				return
			}

			attackedCreatures := make([]*match.Card, 0)

			ctx.Match.NewAction(card.Player, event.AttackableCreatures, 1, 1, "Select the creature to attack", true)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					ctx.InterruptFlow()
					ctx.Match.CloseAction(card.Player)
					return
				}

				if len(action.Cards) != 1 || !match.AssertCardsIn(event.AttackableCreatures, action.Cards[0]) {
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
			if len(event.Blockers) > 0 && !card.HasCondition(cnd.CantBeBlocked) {

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

	if event, ok := ctx.Event.(*match.TapAbility); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't use tap ability because it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		if card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't use tap ability because it is already tapped", card.Name))
			ctx.InterruptFlow()
			return
		}

		if !card.HasCondition(cnd.TapAbility) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s doesn't have any related tap abilities", card.Name))
			ctx.InterruptFlow()
			return
		}

		// According to this https://duelmasters.fandom.com/wiki/Tap_Ability#Details if a creature
		// can't legally attack, it can't use a tap ability either.
		if card.HasCondition(cnd.CantAttackCreatures) && card.HasCondition(cnd.CantAttackPlayers) {
			ctx.Match.WarnPlayer(card.Player, "A card that can't attack can't use tap abilities")
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			tapConditions := make([]*match.Condition, 0)
			tapConditionsSourceCards := make([]*match.Card, 0)

			for _, condition := range card.Conditions() {
				if condition.ID != cnd.TapAbility {
					continue
				}
				tapConditions = append(tapConditions, &condition)
				id, _ := condition.Src.(string)
				sourceCard, err := card.Player.GetCard(id, match.BATTLEZONE)
				if err == nil {
					tapConditionsSourceCards = append(tapConditionsSourceCards, sourceCard)
				}
			}
			var tapEffect interface{}

			if len(tapConditions) > 1 {

				ctx.Match.NewAction(
					card.Player,
					tapConditionsSourceCards,
					1,
					1,
					"Select the source of the tap effect",
					true)

				for {

					action := <-card.Player.Action

					if action.Cancel {
						ctx.InterruptFlow()
						ctx.Match.CloseAction(card.Player)
						return
					}

					if len(action.Cards) != 1 || !match.AssertCardsIn(tapConditionsSourceCards, action.Cards[0]) {
						ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
						continue
					}

					for _, condition := range tapConditions {

						if condition.Src == action.Cards[0] {
							tapEffect = condition.Val
						}

					}

					ctx.Match.CloseAction(card.Player)

					break

				}
			} else {
				tapEffect = tapConditions[0].Val
			}

			if f, ok := tapEffect.(func(card *match.Card, ctx *match.Context)); ok {
				f(card, ctx)
			}

			card.Tapped = true
		})

	}

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {
		if event.Card != card {
			return
		}

		ctx.ScheduleAfter(func() {
			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.GRAVEYARD, card.ID)

			// Slayer
			if event.Context == match.DestroyedInBattle {
				for _, condition := range card.Conditions() {
					if condition.ID != cnd.Slayer {
						continue
					}

					creature, err := ctx.Match.Opponent(card.Player).GetCard(event.Source.ID, match.BATTLEZONE)
					if err != nil {
						break
					}

					if f, ok := condition.Val.(SlayerCondition); ok {
						// conditional slayer
						if f(creature) {
							ctx.Match.Destroy(creature, card, match.DestroyedBySlayer)
						}
					} else {
						// regular slayer
						ctx.Match.Destroy(creature, card, match.DestroyedBySlayer)
					}
				}
			}
		})
	}

	if event, ok := ctx.Event.(*match.Battle); ok && event.Attacker.ID == card.ID {
		ctx.ScheduleAfter(func() {
			m := ctx.Match
			attacker := event.Attacker
			attackerPower := event.AttackerPower
			defender := event.Defender
			defenderPower := event.DefenderPower
			blocked := event.Blocked

			if attackerPower > defenderPower {
				m.HandleFx(match.NewContext(m, &match.CreatureDestroyed{Card: defender, Source: attacker, Blocked: blocked}))
				m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", ctx.Match.FormatDisplayableCard(defender), defenderPower, ctx.Match.FormatDisplayableCard(attacker), attackerPower))
			} else if attackerPower == defenderPower {
				m.HandleFx(match.NewContext(m, &match.CreatureDestroyed{Card: attacker, Source: defender, Blocked: blocked}))
				m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", ctx.Match.FormatDisplayableCard(attacker), attackerPower, ctx.Match.FormatDisplayableCard(defender), defenderPower))
				m.HandleFx(match.NewContext(m, &match.CreatureDestroyed{Card: defender, Source: attacker, Blocked: blocked}))
				m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", ctx.Match.FormatDisplayableCard(defender), defenderPower, ctx.Match.FormatDisplayableCard(attacker), attackerPower))
			} else if attackerPower < defenderPower {
				m.HandleFx(match.NewContext(m, &match.CreatureDestroyed{Card: attacker, Source: defender, Blocked: blocked}))
				m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", ctx.Match.FormatDisplayableCard(attacker), attackerPower, ctx.Match.FormatDisplayableCard(defender), defenderPower))
			}
		})
	}

}
