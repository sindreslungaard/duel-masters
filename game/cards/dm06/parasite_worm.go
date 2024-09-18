package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func GraveWormQ(c *match.Card) {

	c.Name = "Grave Worm Q"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm, family.Survivor}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			if creature.HasFamily(family.Survivor) && event.To == match.BATTLEZONE {
				fx.SelectFilterSelectablesOnly(
					card.Player,
					ctx.Match,
					card.Player,
					match.GRAVEYARD,
					"You may return a survivor from your graveyard to your hand",
					1,
					1,
					true,
					func(x *match.Card) bool { return x.HasFamily(family.Survivor) && x.HasCondition(cnd.Creature) },
				).Map(func(x *match.Card) {
					card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their graveyard by Grave Worm Q's survivor ability", x.Name, card.Player.Username()))

				})

			}

		}

	})

}

func TentacleWorm(c *match.Card) {

	c.Name = "Tentacle Worm"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)

}
