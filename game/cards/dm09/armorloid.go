package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SimianWarriorGrash ...
func SimianWarriorGrash(c *match.Card) {

	c.Name = "Simian Warrior Grash"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				exit()
				return

			}

			if event, ok := ctx2.Event.(*match.CreatureDestroyed); ok && event.Card.Player == c.Player && event.Card.HasFamily(family.Armorloid) {
				fx.OpponentChoosesManaBurn(card, ctx2)
			}

		})

	}))

}
