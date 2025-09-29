package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MessaBanhaExpanseGuardian ...
func MessaBanhaExpanseGuardian(c *match.Card) {

	c.Name = "Messa Banha, Expanse Guardian"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.Block); ok {
			if event.Attacker.Player == card.Player ||
				event.Attacker.Zone != match.BATTLEZONE ||
				card.Zone != match.BATTLEZONE {
				return
			}

			//@TODO test this
			if len(event.Blockers) > 0 && event.Attacker.IsBlockable(ctx) {
				for _, b := range event.Blockers {
					if b.ID == card.ID {
						// Force the battle between the attacker and this card
						ctx.InterruptFlow()
						ctx.Match.Battle(event.Attacker, card, true, len(event.ShieldsAttacked) > 0)
						return
					}
				}
			}
		}
	})

}
