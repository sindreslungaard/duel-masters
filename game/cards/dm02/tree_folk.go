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

		if !ctx.Match.IsPlayerTurn(card.Player) {
			return
		}

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

			if card.Zone == match.BATTLEZONE && event.CardID != card.ID {

				toReduce, err := ctx.Match.CurrentPlayer().Player.GetCard(event.CardID, match.HAND)

				if err != nil {
					return
				}

				toReduce.AddCondition(cnd.ReducedCost, true, card.ID)

			}

		}

	})

}
