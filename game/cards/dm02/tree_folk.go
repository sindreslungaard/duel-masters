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
	c.Family = []string{family.TreeFolk}
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

				if toReduce.HasCondition(cnd.Creature) {
					toReduce.AddUniqueSourceCondition(cnd.ReducedCost, 1, card.ID)
				}

			}

		}

	})

}

// EssenceElf ...
func EssenceElf(c *match.Card) {

	c.Name = "Essence Elf"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 2
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

				if toReduce.HasCondition(cnd.Spell) {
					toReduce.AddUniqueSourceCondition(cnd.ReducedCost, 1, card.ID)
				}

			}

		}

	})

}
