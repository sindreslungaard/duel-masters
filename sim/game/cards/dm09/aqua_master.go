package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AquaMaster ...
func AquaMaster(c *match.Card) {

	c.Name = "Aqua Master"
	c.Power = 4000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldsSelectionEffect,
		fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked,
			func(card *match.Card, ctx *match.Context) {
				opponent := ctx.Match.Opponent(card.Player)

				// A persistent effect (e.g. Meloppe) may override who chooses the
				// shield via the SelectShields event's Chooser field, same as the
				// default shield-break selection in fx.Creature.
				chooser := card.Player
				if event, ok := ctx.Event.(*match.SelectShields); ok && event.Chooser != nil {
					chooser = event.Chooser
				}

				fx.SelectBackside(
					chooser,
					ctx.Match,
					opponent,
					match.SHIELDZONE,
					fmt.Sprintf("%s's effect: Choose one of %s shield and turn it face up.%s", card.Name, fx.ShieldPossessive(chooser, opponent), fx.MeloppeNote(chooser, card.Player)),
					1,
					1,
					false,
				).Map(func(x *match.Card) {
					// Grabbed before it's flipped face up: not that it matters here
					// since the shield stays put, but it keeps this consistent with
					// every other reveal that has to grab it before the shield can
					// move or change.
					description := fx.DescribeShield(x)

					x.ShieldFaceUp = true
					ctx.Match.BroadcastState()
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was turned face up from %s's shieldzone.", description, opponent.Username()))
				})
			}))
}
