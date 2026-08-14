package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
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
	// handlers run in every zone, so plain state on the card instance is what
	// carries the promise, rather than conditions that are cleared each turn.
	extraTurnPending := false
	loseAfterExtraTurn := false

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

			for _, creature := range toDestroy {
				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
			}

			if !ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			extraTurnPending = true
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s will take an extra turn and then lose the game at the end of it", card.Name, card.Player.Username()))
		}),
		func(card *match.Card, ctx *match.Context) {

			if _, ok := ctx.Event.(*match.EndOfTurnStep); !ok {
				return
			}

			if !ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			// Scheduled so every other end of turn effect has resolved before
			// the turn is either repeated or the game is ended.
			ctx.ScheduleAfter(func() {

				if extraTurnPending {
					extraTurnPending = false
					loseAfterExtraTurn = true

					// Cancelling stops the engine handing the turn over, so the
					// repeated turn below belongs to the same player.
					ctx.InterruptFlow()
					ctx.Match.BeginNewTurn(true)

					return
				}

				if loseAfterExtraTurn {
					loseAfterExtraTurn = false

					opponent := ctx.Match.Opponent(card.Player)
					ctx.Match.End(opponent, fmt.Sprintf("%s won the game because of %s", ctx.Match.PlayerRef(opponent).Socket.User.Username, card.Name))
				}
			})
		},
	)

}
