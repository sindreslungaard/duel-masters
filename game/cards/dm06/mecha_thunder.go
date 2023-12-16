package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func LaveilSeekerOfCatastrophe(c *match.Card) {

	c.Name = "Laveil Seeker Of Catastrophe"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.Doublebreaker, fx.Untap)

}

func DavaToreySeekerOfClouds(c *match.Card) {

	c.Name = "Dava Torey Seeker Of Cloudss"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx2 *match.Context) {

		if !ctx2.Match.IsPlayerTurn(card.Player) {
			if event, ok := ctx2.Event.(*match.CardMoved); ok {
				if event.CardID == card.ID && event.From == match.HAND && event.To == match.GRAVEYARD {
					ctx2.InterruptFlow()
					card.Player.MoveCard(card.ID, match.GRAVEYARD, match.BATTLEZONE)
					ctx2.Match.Chat("Server", fmt.Sprintf("%s was discarded and moved to the battle zone", card.Name))
				}
			}
		}

	})

}
