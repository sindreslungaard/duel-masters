package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TerradragonRegarion ...
func TerradragonRegarion(c *match.Card) {

	c.Name = "Terradragon Regarion"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000)
}

// SuperTerradragonBailasGale ...
func SuperTerradragonBailasGale(c *match.Card) {

	c.Name = "Super Terradragon Bailas Gale"
	c.Power = 9000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.DragonEvolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.SpellCast); ok && c.Zone == match.BATTLEZONE && event.FromShield {

			// check which player played the spell
			var p *match.Player
			if event.MatchPlayerID == 1 {
				p = ctx.Match.Player1.Player
			} else {
				p = ctx.Match.Player2.Player
			}

			if p == c.Player {

				spell, err := p.GetCard(event.CardID, match.HAND)
				if err != nil {
					return
				}

				// prevents card from being sent to grave
				// uses the fact that cards in the battlezone are handled before ones in hand
				ctx.ScheduleAfter(func() {
					ctx.InterruptFlow()
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the hand by %s", spell.Name, c.Name))
				})

			}

		}

	})
}
