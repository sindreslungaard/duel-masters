package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
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
