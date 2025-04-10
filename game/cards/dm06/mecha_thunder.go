package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func LaveilSeekerOfCatastrophe(c *match.Card) {

	c.Name = "Laveil, Seeker of Catastrophe"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.Doublebreaker, fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}

func DavaToreySeekerOfClouds(c *match.Card) {

	c.Name = "Dava Torey, Seeker of Clouds"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {
			if event.CardID == card.ID && event.From == match.HAND && event.To == match.GRAVEYARD {
				if !ctx.Match.IsPlayerTurn(card.Player) {
					ctx.ScheduleAfter(func() {
						card.Player.MoveCard(card.ID, match.GRAVEYARD, match.BATTLEZONE, card.ID)
						ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was discarded and moved to the battle zone", card.Name))
					})

				}
			}
		}

	})

}
