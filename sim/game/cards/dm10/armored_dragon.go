package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"strings"
)

// GaulezalDragon ...
func GaulezalDragon(c *match.Card) {
	c.Name = "Gaulezal Dragon"
	c.Power = 11000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker)
}

const (
	bombazarPendingKey       = "bombazar_pending"
	bombazarPendingExtraTurn = "extra_turn"
	bombazarPendingLose      = "lose"
)

// bombazarPendingQueue reads a Bombazar's outstanding extra-turn/lose promises,
// oldest first. There can be more than one: destroying an earlier copy of this
// card hands its still-unresolved promise to whichever copy destroyed it,
// rather than dropping it, so a second Bombazar inherits the first one's debt
// on top of creating its own.
func bombazarPendingQueue(card *match.Card) []string {
	raw, ok := card.LocalData(bombazarPendingKey)
	if !ok || raw == "" {
		return nil
	}

	return strings.Split(raw, ",")
}

func setBombazarPendingQueue(card *match.Card, queue []string) {
	card.SetLocalData(bombazarPendingKey, strings.Join(queue, ","))
}

func pushBombazarPending(card *match.Card, state string) {
	setBombazarPendingQueue(card, append(bombazarPendingQueue(card), state))
}

// BombazarDragonOfDestiny ...
func BombazarDragonOfDestiny(c *match.Card) {

	c.Name = "Bombazar, Dragon of Destiny"
	c.Power = 6000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.ArmoredDragon, family.EarthDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	// The extra turn and the loss that follows it outlive the creature: once it
	// has entered play, bouncing or destroying it does not call either off. Card
	// handlers run in every zone, so state stored directly on the card instance
	// (via LocalData, which conditions do not offer) is what carries the
	// promise, rather than conditions that are cleared each turn.
	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.SpeedAttacker,
		fx.Doublebreaker,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			// "Destroy all other creatures that have 6000 power" is an exact
			// figure, not a threshold, and it spares Bombazar itself even
			// though it has exactly that power.
			toDestroy := make([]*match.Card, 0)

			for _, player := range []*match.Player{card.Player, ctx.Match.Opponent(card.Player)} {
				fx.FindFilter(
					player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.ID != card.ID && ctx.Match.GetPower(x, false) == 6000
					},
				).Map(func(x *match.Card) {
					toDestroy = append(toDestroy, x)
				})
			}

			// Only a same-controller Bombazar's debt transfers here. Destroying
			// an opponent's copy of this card must not make their own promise
			// (checked against their own turn, not this player's) mine to pay.
			for _, creature := range toDestroy {
				if creature.Name != card.Name || creature.Player != card.Player {
					continue
				}

				if queue := bombazarPendingQueue(creature); len(queue) > 0 {
					setBombazarPendingQueue(card, append(bombazarPendingQueue(card), queue...))
					setBombazarPendingQueue(creature, nil)
				}
			}

			for _, creature := range toDestroy {
				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
			}

			if ctx.Match.IsPlayerTurn(card.Player) {
				pushBombazarPending(card, bombazarPendingExtraTurn)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s will take an extra turn and then lose the game at the end of it", card.Name, card.Player.Username()))
				return
			}

			// Put into the battle zone on someone else's turn, for example by an
			// effect that has its owner put it into play off-turn. "The turn
			// after this one" is then the owner's own next turn, which already
			// follows naturally, so no turn needs to be repeated: only the loss
			// remains to be scheduled for the end of that turn.
			pushBombazarPending(card, bombazarPendingLose)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s will lose the game at the end of their next turn", card.Name, card.Player.Username()))
		}),
		func(card *match.Card, ctx *match.Context) {

			if _, ok := ctx.Event.(*match.EndOfTurnStep); !ok {
				return
			}

			if !ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			// Scheduled so every other end of turn effect has resolved before
			// the oldest promise either repeats the turn or ends the game.
			ctx.ScheduleAfter(func() {

				queue := bombazarPendingQueue(card)

				if len(queue) == 0 {
					return
				}

				switch queue[0] {

				case bombazarPendingExtraTurn:
					queue[0] = bombazarPendingLose
					setBombazarPendingQueue(card, queue)

					// Cancelling stops the engine handing the turn over, so the
					// repeated turn below belongs to the same player. Only the
					// oldest promise acts this step: a newer one behind it,
					// still an extra turn of its own, waits for a step of its
					// own rather than doubling up on this one.
					ctx.InterruptFlow()
					ctx.Match.BeginNewTurn(true)

				case bombazarPendingLose:
					setBombazarPendingQueue(card, queue[1:])

					opponent := ctx.Match.Opponent(card.Player)
					ctx.Match.End(opponent, fmt.Sprintf("%s won the game because of %s", ctx.Match.PlayerRef(opponent).Socket.User.Username, card.Name))
				}
			})
		},
	)

}
