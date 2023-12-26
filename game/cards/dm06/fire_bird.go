package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CoccoLupia(c *match.Card) {

	c.Name = "Cocco Lupia"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

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

				if toReduce.HasFamily(family.ArmoredDragon) {
					toReduce.AddUniqueSourceCondition(cnd.ReducedCost, 2, card.ID)
				}

			}

		}

	})

}
