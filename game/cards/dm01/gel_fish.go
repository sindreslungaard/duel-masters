package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// IllusionaryMerfolk ...
func IllusionaryMerfolk(c *match.Card) {

	c.Name = "Illusionary Merfolk"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				battlezone, err := card.Player.Container(match.BATTLEZONE)

				if err != nil {
					return
				}

				for _, creature := range battlezone {

					if creature.HasFamily(family.CyberLord) {
						card.Player.DrawCards(3)
						return
					}

				}

			}

		}

	})

}

// PhantomFish ...
func PhantomFish(c *match.Card) {

	c.Name = "Phantom Fish"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures)

}

// RevolverFish ...
func RevolverFish(c *match.Card) {

	c.Name = "Revolver Fish"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures)

}

// SaucerHeadShark ...
func SaucerHeadShark(c *match.Card) {

	c.Name = "Saucer-Head Shark"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				myBattlezone, err := card.Player.Container(match.BATTLEZONE)
				if err != nil {
					return
				}

				opponentBattlezone, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)
				if err != nil {
					return
				}

				for _, creature := range myBattlezone {
					if ctx.Match.GetPower(creature, false) <= 2000 {
						creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND)
						ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by Saucer-Head Shark", creature.Name, creature.Player.Username()))
					}
				}

				for _, creature := range opponentBattlezone {
					if ctx.Match.GetPower(creature, false) <= 2000 {
						creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND)
						ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by Saucer-Head Shark", creature.Name, creature.Player.Username()))
					}
				}

			}

		}

	})

}
