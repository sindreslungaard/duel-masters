package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BonePiercer ...
func BonePiercer(c *match.Card) {

	c.Name = "Bone Piercer"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = family.BrainJacker
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.From == match.BATTLEZONE && event.To == match.GRAVEYARD {

				cards := match.Search(card.Player, ctx.Match, card.Player, match.MANAZONE, "Bone Piercer: Select 1 card from your manazone that will be sent to your hand", 0, 1, true)

				for _, crd := range cards {
					card.Player.MoveCard(crd.ID, match.MANAZONE, match.HAND)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their mana zone", crd.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				}

			}

		}

	})

}