package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

// BallomMasterOfDeath ...
func BallomMasterOfDeath(c *match.Card) {

	c.Name = "Ballom, Master of Death"
	c.Power = 12000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ != civ.Darkness },
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Ballom, Master of Death", x.Name))
		})

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ != civ.Darkness },
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Ballom, Master of Death", x.Name))
		})
	}))

}

// TroxGeneralOfDestruction ...
func TroxGeneralOfDestruction(c *match.Card) {

	c.Name = "Trox, General of Destruction"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Darkness && x.ID != card.ID },
		).Map(func(x *match.Card) {
			hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)
			if err != nil || len(hand) == 0 {
				return
			}

			discarded, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
			if err != nil {
				return
			}

			ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand by Trox, General of Destruction", discarded.Name, discarded.Player.Username()))
		})

	}))
}

// PhotocideLordOfTheWastes ...
func PhotocideLordOfTheWastes(c *match.Card) {

	c.Name = "Photocide, Lord of the Wastes"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.AttackCreature)

		if !ok || event.CardID != card.ID {
			return
		}

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.Civ == civ.Light && x.Tapped == false
			},
		).Map(func(x *match.Card) {
			// don't add if already in the list of attackable creatures
			for _, creature := range event.AttackableCreatures {
				if creature.ID == x.ID {
					return
				}
			}

			event.AttackableCreatures = append(event.AttackableCreatures, x)
		})

	})
}
