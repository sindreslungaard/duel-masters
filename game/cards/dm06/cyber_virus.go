package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SteamStar(c *match.Card) {

	c.Name = "Steam Star"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)
}

func RippleLotusQ(c *match.Card) {

	c.Name = "Ripple Lotus Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

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
				creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creature and tap it. Close to not tap any creatures.", 1, 1, true)
				for _, creature := range creatures {
					creature.Tapped = true
				}

			}

		}

	})

}
