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
			ctx.Match.WarnPlayer(card.Player, "You can't summon creatures after attacking or using tap ability.")
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
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s summoned %s to the battle zone", card.Player.Username(), card.Name))

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

		if card.HasCondition(cnd.CantAttackPlayers) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))
			ctx.InterruptFlow()
			return
		}

		if card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack while tapped", card.Name))
			ctx.InterruptFlow()
			return
		}

		if HasSummoningSickness(card) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {
			if card.HasCondition(cnd.HasShieldsSelectionEffect) {
				// Only for cards that have a specific effect that needs Attack Confirmation BEFORE Selecting Shields
				// Prompt the user for confirmation of attack player sequence
				// After this step, user will NOT be able to cancel the attack player sequence anymore
				if BinaryQuestion(
					card.Player,
					ctx.Match,
					fmt.Sprintf("Do you CONFIRM to Attack your Opponent? This is irreversible and cannot be cancelled afterwards.\r\nYou will be prompted to select shields after %s's effect will be resolved.", card.Name),
				) {
					if !tapCardAndConfirmAttack(card, ctx, true) {
						return
					}

					selectBlockersEvent := match.SelectBlockers{Blockers: make([]*match.Card, 0), Attacker: card, AttackedCardID: ""}
					ctx.Match.HandleFx(match.NewContext(ctx.Match, &selectBlockersEvent))
					ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.Block{Blockers: selectBlockersEvent.Blockers, Attacker: selectBlockersEvent.Attacker, AttackedCardID: selectBlockersEvent.AttackedCardID}))
				} else {
					ctx.InterruptFlow()
					return
				}
			} else {
				selectShieldsEvent := match.SelectShields{Attacker: card, Cancellable: true}
				selectShieldsCtx := match.NewContext(ctx.Match, &selectShieldsEvent)
				ctx.Match.HandleFx(selectShieldsCtx)

				if !tapCardAndConfirmAttack(card, selectShieldsCtx, true) {
					return
				}

				selectBlockersEvent := match.SelectBlockers{Blockers: make([]*match.Card, 0), Attacker: card, AttackedCardID: ""}
				ctx.Match.HandleFx(match.NewContext(ctx.Match, &selectBlockersEvent))
				ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.Block{Blockers: selectBlockersEvent.Blockers, Attacker: selectBlockersEvent.Attacker, AttackedCardID: selectBlockersEvent.AttackedCardID, ShieldsAttacked: selectShieldsEvent.ShieldsAttacked}))
			}
		})

	}

	// Attack a creature
	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if card.HasCondition(cnd.CantAttackCreatures) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))
			ctx.InterruptFlow()
			return
		}

		if card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack while tapped", card.Name))
			ctx.InterruptFlow()
			return
		}

		if HasSummoningSickness(card) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot attack this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		opponent := ctx.Match.Opponent(card.Player)

		battlezone, err := opponent.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		// Add attackable creatures
		for _, c := range battlezone {
			if c.HasCondition(cnd.TreatedAsTapped) || (!c.HasCondition(cnd.CantBeAttacked) && (c.Tapped || card.HasCondition(cnd.AttackUntapped))) {
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

			if !tapCardAndConfirmAttack(card, ctx, false) {
				return
			}

			// In case the creature was removed by an AttackConfirmed effect
			if c.Zone != match.BATTLEZONE {
				return
			}

			selectBlockersEvent := match.SelectBlockers{Blockers: make([]*match.Card, 0), Attacker: card, AttackedCardID: c.ID}
			ctx.Match.HandleFx(match.NewContext(ctx.Match, &selectBlockersEvent))
			ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.Block{Blockers: selectBlockersEvent.Blockers, Attacker: selectBlockersEvent.Attacker, AttackedCardID: selectBlockersEvent.AttackedCardID}))
		})

	}

	if event, ok := ctx.Event.(*match.AttackConfirmed); ok {
		if event.CardID != card.ID || card.Zone != match.BATTLEZONE {
			return
		}

		// To prevent the AttackStep to be incorrectly set whenever
		// someone would cancel the AttackCreature / AttackPlayer sequence
		if !ctx.Cancelled() {
			ctx.Match.Step = &match.AttackStep{}
			card.Player.CanChargeMana = false
		}
	}

	if event, ok := ctx.Event.(*match.Block); ok {

		// Is this event for me or someone else?
		if event.Attacker != card ||
			event.Attacker.Zone != match.BATTLEZONE {
			return
		}

		opponent := ctx.Match.Opponent(card.Player)
		oppShieldZone, _ := opponent.Container(match.SHIELDZONE)
		oppShieldsZoneLength := len(oppShieldZone)

		// (Attack Player case) Allow the opponent to block if they can
		if event.AttackedCardID == "" {

			if len(event.Blockers) > 0 && !card.HasCondition(cnd.CantBeBlocked) && !stealthActive(card, ctx) {

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

				ctx.Match.NewAction(opponent, event.Blockers, 1, 1, fmt.Sprintf("%s (%v) is attacking you. Choose a creature to block the attack with or close to not block the attack.", card.Name, ctx.Match.GetPower(card, true)), true)

				for {

					action := <-opponent.Action

					if action.Cancel {
						ctx.Match.EndWait(card.Player)
						ctx.Match.CloseAction(opponent)

						if oppShieldsZoneLength < 1 {
							// Win
							ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
						} else {
							// Break n shields
							if card.HasCondition(cnd.HasShieldsSelectionEffect) {
								selectShieldsEvent := match.SelectShields{Attacker: card, Cancellable: false}
								ctx.Match.HandleFx(match.NewContext(ctx.Match, &selectShieldsEvent))

								if len(selectShieldsEvent.ShieldsAttacked) > 0 {
									ctx.Match.BreakShields(selectShieldsEvent.ShieldsAttacked, card)
								}
							} else {
								ctx.Match.BreakShields(event.ShieldsAttacked, card)
							}
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

					ctx.Match.Battle(card, c, true, true)

					break

				}

				return

			} else {

				if oppShieldsZoneLength < 1 {
					// Win
					ctx.Match.End(card.Player, fmt.Sprintf("%s won the game", ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				} else {
					// Break n shields
					if card.HasCondition(cnd.HasShieldsSelectionEffect) {
						selectShieldsEvent := match.SelectShields{Attacker: card, Cancellable: false}
						ctx.Match.HandleFx(match.NewContext(ctx.Match, &selectShieldsEvent))

						if len(selectShieldsEvent.ShieldsAttacked) > 0 {
							ctx.Match.BreakShields(selectShieldsEvent.ShieldsAttacked, card)
						}
					} else {
						ctx.Match.BreakShields(event.ShieldsAttacked, card)
					}
				}

				return

			}

		} else {
			// (Attack Creature case) Allow the opponent to block if they can
			attackedCard, err := opponent.GetCard(event.AttackedCardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			if len(event.Blockers) > 0 && !card.HasCondition(cnd.CantBeBlocked) && !stealthActive(card, ctx) {

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

				ctx.Match.NewAction(opponent, event.Blockers, 1, 1, fmt.Sprintf("%s (%v) is attacking %s (%v). Choose a creature to block the attack with or close to not block the attack.", card.Name, ctx.Match.GetPower(card, true), attackedCard.Name, ctx.Match.GetPower(attackedCard, false)), true)

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

					ctx.Match.Battle(card, blocker, true, false)

					return

				}

			}

			ctx.Match.Battle(card, attackedCard, false, false)

		}
	}

	if event, ok := ctx.Event.(*match.SelectShields); ok {

		// Is this event for me or someone else?
		if event.Attacker.ID != card.ID {
			return
		}

		ctx.ScheduleAfter(func() {
			shieldsAttacked := SelectAndReturnShields(
				card,
				ctx,
				event.Cancellable,
			)

			event.ShieldsAttacked = shieldsAttacked
		})

	}

	if event, ok := ctx.Event.(*match.TapAbility); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if HasSummoningSickness(card) {
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
				ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.ActivatesTapAbility{
					CardID: card.ID,
				}))
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activates tap effect", card.Name))
				f(card, ctx)
			}

			card.Tapped = true
		})

	}

	if event, ok := ctx.Event.(*match.ActivatesTapAbility); ok {
		if event.CardID != card.ID || card.Zone != match.BATTLEZONE {
			return
		}

		// Same as in AttackConfirmed event,
		// To prevent the AttackStep to be incorrectly set whenever
		// someone would cancel the TapAbility selection (in case a creature has multiple tap abilities)
		if !ctx.Cancelled() {
			// Tap abilities can only be used during attack step
			// https://duelmasters.fandom.com/wiki/Step#Step_7_(Attack_step)
			ctx.Match.Step = &match.AttackStep{}
			card.Player.CanChargeMana = false
		}
	}

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {
		if event.Card != card {
			return
		}

		ctx.ScheduleAfter(func() {
			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.GRAVEYARD, event.Source.ID)
		})
	}

	if event, ok := ctx.Event.(*match.Battle); ok && event.Attacker.ID == card.ID {
		ctx.ScheduleAfter(func() {
			attacker := event.Attacker
			attackerPower := event.AttackerPower
			defender := event.Defender
			defenderPower := event.DefenderPower
			blocked := event.Blocked

			if attackerPower > defenderPower {
				handleBattle(ctx, attacker, attackerPower, defender, defenderPower, blocked)
			} else if attackerPower == defenderPower {
				handleBattle(ctx, attacker, attackerPower, defender, defenderPower, blocked)
				handleBattle(ctx, defender, defenderPower, attacker, attackerPower, blocked)
			} else if attackerPower < defenderPower {
				handleBattle(ctx, defender, defenderPower, attacker, attackerPower, blocked)
			}
		})
	}

}

func handleBattle(ctx *match.Context, winner *match.Card, winnerPower int, looser *match.Card, looserPower int, blocked bool) {
	ctx.Match.ReportActionInChat(looser.Player, fmt.Sprintf("%s (%v) was destroyed by %s (%v)", looser.Name, looserPower, winner.Name, winnerPower))
	ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.CreatureDestroyed{Card: looser, Source: winner, Blocked: blocked}))

	// Destroy after battle
	for _, condition := range winner.Conditions() {
		if condition.ID == cnd.DestroyAfterBattle {
			ctx.Match.Destroy(winner, winner, match.DestroyedInBattle)

			break
		}
	}

	// Slayer
	hasSlayer := false

	for _, condition := range looser.Conditions() {
		if condition.ID != cnd.Slayer {
			continue
		}

		if f, ok := condition.Val.(SlayerCondition); ok {
			// conditional slayer
			hasSlayer = hasSlayer || f(winner)
		} else {
			// regular slayer
			hasSlayer = true
		}
	}

	if hasSlayer {
		ctx.Match.Destroy(winner, looser, match.DestroyedBySlayer)
	}
}

func stealthActive(card *match.Card, ctx *match.Context) bool {
	if !card.HasCondition(cnd.Stealth) {
		return false
	}

	for _, cond := range card.Conditions() {
		if cond.ID != cnd.Stealth {
			continue
		}
		if match.ContainerHas(
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			func(c *match.Card) bool { return c.Civ == cond.Val },
		) {
			return true
		}
	}

	return false
}

func tapCardAndConfirmAttack(card *match.Card, ctx *match.Context, attackPlayer bool) bool {
	if ctx.Cancelled() {
		return false
	}

	card.Tapped = true

	// Broadcast state so that opponent can see that this card is tapped if they get any shield triggers
	ctx.Match.BroadcastState()

	ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.AttackConfirmed{CardID: card.ID, Player: attackPlayer, Creature: !attackPlayer}))

	// In case AttackConfirmed effect removes itself (current attacking card) from the Battlezone
	return card.Zone == match.BATTLEZONE
}

func HasSummoningSickness(card *match.Card) bool {
	return card.HasCondition(cnd.SummoningSickness) && !card.HasCondition(cnd.SpeedAttacker)
}

func SelectAndReturnShields(card *match.Card, ctx *match.Context, cancellable bool) []*match.Card {
	opponent := ctx.Match.Opponent(card.Player)
	shieldzone, err := opponent.Container(match.SHIELDZONE)

	if err != nil || len(shieldzone) == 0 {
		return []*match.Card{}
	}

	noOfShields := 1

	if card.HasCondition(cnd.DoubleBreaker) {
		noOfShields = 2
	}

	if card.HasCondition(cnd.TripleBreaker) {
		noOfShields = 3
	}

	for _, condition := range card.Conditions() {
		if condition.ID != cnd.ShieldBreakModifier {
			continue
		}

		if val, ok := condition.Val.(int); ok {
			noOfShields += val
		}
	}

	if noOfShields > len(shieldzone) {
		noOfShields = len(shieldzone)
	}

	shieldsAttacked := make([]*match.Card, 0)

	ctx.Match.NewBacksideAction(card.Player, shieldzone, noOfShields, noOfShields, fmt.Sprintf("Select %v shield(s) to break", noOfShields), cancellable)

	for {
		action := <-card.Player.Action

		if action.Cancel {
			ctx.InterruptFlow()
			ctx.Match.CloseAction(card.Player)
			return []*match.Card{}
		}

		if len(action.Cards) != noOfShields || !match.AssertCardsIn(shieldzone, action.Cards[0]) {
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

		return shieldsAttacked
	}
}
