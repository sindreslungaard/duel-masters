package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SparkleFlower ...
func SparkleFlower(c *match.Card) {

	c.Name = "Sparkle Flower"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.StarlightTree
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.UntapStep); ok {

			if match.ContainerHas(card.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Light }) {
				card.RemoveCondition(cnd.Blocker)
			} else {
				card.AddCondition(cnd.Blocker, true, card.ID)	
			}
		}
	
		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
	
			if card.HasCondition(cnd.Blocker) == false {
				return
			}

			// Only add to list of blockers if it is our player that is being attacked, i.e. not our players turn
			if !ctx.Match.IsPlayerTurn(card.Player) && !card.Tapped && card.Player.HasCard(match.BATTLEZONE, card.ID) {
				event.Blockers = append(event.Blockers, card)
			}
	
		}
	
		if event, ok := ctx.Event.(*match.AttackCreature); ok {
	
			if card.HasCondition(cnd.Blocker) == false {
				return
			}

			// Only add to list of blockers if it is our creature that is being attacked, i.e. not our players turn
			if !ctx.Match.IsPlayerTurn(card.Player) && !card.Tapped && card.Player.HasCard(match.BATTLEZONE, card.ID) {
				event.Blockers = append(event.Blockers, card)
			}
	
		}
	})

}
