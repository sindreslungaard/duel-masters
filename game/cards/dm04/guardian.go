package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GulanRiasSpeedGuardian ...
func GulanRiasSpeedGuardian(c *match.Card) {

	c.Name = "Gulan Rias, Speed Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = family.Guardian
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(
		fx.Creature,
		fx.CantBeBlockedIf(func(blocker *match.Card) bool {
			return blocker.Civ == civ.Darkness
		}),
		fx.CantBeAttackedIf(func(attacker *match.Card) bool {
			return attacker.Civ == civ.Darkness
		}),
	)
}

// MistRiasSonicGuardian ...
func MistRiasSonicGuardian(c *match.Card) {

	c.Name = "Mist Rias, Sonic Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = family.Guardian
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok && event.To == match.BATTLEZONE && event.CardID != card.ID {

			if !ctx.Match.IsPlayerTurn(card.Player) {
				ctx.Match.Wait(ctx.Match.Opponent(card.Player), "Waiting for your opponent to make an action")
				defer ctx.Match.EndWait(ctx.Match.Opponent(card.Player))
			}

			result := fx.SelectBacksideFilter(card.Player, ctx.Match, card.Player, match.BATTLEZONE, fmt.Sprintf("%s: You may draw a card. Click close to not draw a card.", card.Name), 1, 1, true, func(x *match.Card) bool {
				return x.ID == c.ID
			})

			if len(result) > 0 {
				card.Player.DrawCards(1)
				ctx.Match.Chat("Server", fmt.Sprintf("%s chose to draw a card from %s's ability", card.Player.Username(), card.Name))
			}
		}

	})
}
