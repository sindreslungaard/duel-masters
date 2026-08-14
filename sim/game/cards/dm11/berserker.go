package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LamielDestinyEnforcer ...
func LamielDestinyEnforcer(c *match.Card) {

	c.Name = "Lamiel, Destiny Enforcer"
	c.Power = 3000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(func(card *match.Card, ctx *match.Context) {
		event, ok := ctx.Event.(*match.CreatureDestroyed)

		// "One of your creatures" includes Lamiel itself, and the trigger is
		// limited to the opponent's turn.
		if !ok ||
			event.Card.Player != card.Player ||
			ctx.Match.IsPlayerTurn(card.Player) {
			return
		}

		if !fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s's effect: %s was destroyed. Do you want to draw a card?", card.Name, event.Card.Name)) {
			return
		}

		card.Player.DrawCards(1)
	}))

}
