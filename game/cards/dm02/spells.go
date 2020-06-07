package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BurstShot ...
func BurstShot(c *match.Card) {

	c.Name = "Burst Shot"
	c.Civ = civ.Fire
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			opponent := ctx.Match.Opponent(card.Player)

			myCreatures, err := card.Player.Container(match.BATTLEZONE)
			if err != nil {
				return
			}

			opponentCreatures, err := opponent.Container(match.BATTLEZONE)
			if err != nil {
				return
			}

			for _, creature := range myCreatures {
				if ctx.Match.GetPower(creature, false) <= 2000 {
					ctx.Match.Destroy(creature, card)
				}
			}

			for _, creature := range opponentCreatures {
				if ctx.Match.GetPower(creature, false) <= 2000 {
					ctx.Match.Destroy(creature, card)
				}
			}

		}

	})

}

// LogicCube ...
func LogicCube(c *match.Card) {

	c.Name = "Logic Cube"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, card.Player, match.DECK, "Select 1 spell from your deck that will be shown to your opponent and sent to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

			for _, creature := range creatures {

				card.Player.MoveCard(creature.ID, match.DECK, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the deck to their hand", card.Player.Username(), creature.Name))

			}

			card.Player.ShuffleDeck()

		}

	})

}
