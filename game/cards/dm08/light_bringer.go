package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// NarielTheOracle ...
func NarielTheOracle(c *match.Card) {

	c.Name = "Nariel, the Oracle"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			creature, err := ctx.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if ctx.Match.GetPower(creature, false) >= 3000 {
				ctx.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't attack due to %s's effect.", creature.Name, card.Name))
				ctx.InterruptFlow()
			}
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			creature, err := ctx.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if ctx.Match.GetPower(creature, false) >= 3000 {
				ctx.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't attack due to %s's effect.", creature.Name, card.Name))
				ctx.InterruptFlow()
			}
		}

		if event, ok := ctx.Event.(*match.TapAbility); ok {
			creature, err := ctx.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if ctx.Match.GetPower(creature, false) >= 3000 {
				ctx.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't use tap ability due to %s's effect.", creature.Name, card.Name))
				ctx.InterruptFlow()
			}
		}
	})

}
