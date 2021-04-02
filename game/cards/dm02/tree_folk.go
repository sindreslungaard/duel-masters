package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ElfX ...
func ElfX(c *match.Card) {

	c.Name = "Elf-X"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = family.TreeFolk
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		// Only active if ElfX is in the battlezone
		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.GetCostEvent); ok {

			// Only reduce your other creatures
			if event.Card.Player != card.Player || !event.Card.HasCondition(cnd.Creature) || event.Card == card {
				return
			}

			reduction := 1

			if event.Cost-reduction > 0 {
				event.Cost -= reduction
			}
		}

	})

}
