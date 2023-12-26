package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func CantankerousGiant(c *match.Card) {

	c.Name = "Cantankerous Giant"
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker)
}

func CliffcrushGiant(c *match.Card) {

	c.Name = "Cliffcrush Giant"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return !card.Tapped })) > 1 {
				card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
			} else {
				card.RemoveCondition(cnd.CantAttackCreatures)
			}
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackCreatures) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))

			ctx.InterruptFlow()

		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return !card.Tapped })) > 1 {
				card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
			} else {
				card.RemoveCondition(cnd.CantAttackPlayers)
			}
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackPlayers) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))

			ctx.InterruptFlow()
		}
	})
}
