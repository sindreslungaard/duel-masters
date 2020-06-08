package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CrystalLancer ...
func CrystalLancer(c *match.Card) {

	c.Name = "Crystal Lancer"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature, fx.Evolution, fx.Doublebreaker)

}

// CrystalPaladin ...
func CrystalPaladin(c *match.Card) {

	c.Name = "Crystal Paladin"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.BATTLEZONE {
				return
			}

			myBattlezone, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			opponentBattlezone, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			blockers := append(myBattlezone, opponentBattlezone...)

			for _, blocker := range blockers {

				_, err := blocker.Player.MoveCard(blocker.ID, match.BATTLEZONE, match.HAND)

				if err != nil {
					continue
				}

				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s battle zone to their hand by Crystal Paladin", blocker.Name, blocker.Player.Username()))

			}

		}

	})

}
