package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KingNeptas ...
func KingNeptas(c *match.Card) {

	c.Name = "King Neptas"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards := fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
		)

		opponentCards := fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
		)

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			c.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose up to 1 creature in the battle zone and return it to its owner hand", card.Name),
			0,
			1,
			true,
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
			if err != nil {
				return
			}
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand", x.Name, x.Player.Username()))
		})

	}))

}

// KingPonitas ...
func KingPonitas(c *match.Card) {

	c.Name = "King Ponitas"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		fx.SearchDeckTakeCards(
			card,
			ctx,
			1,
			func(x *match.Card) bool { return x.Civ == civ.Water },
			"water card",
		)
	}))
}

// LegendaryBynor ...
func LegendaryBynor(c *match.Card) {

	c.Name = "Legendary Bynor"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE || !ctx.Match.IsPlayerTurn(card.Player) {
			return
		}

		legendaryBynorSpecial(card, ctx)

	})

}

func legendaryBynorSpecial(card *match.Card, ctx *match.Context) {

	p := ctx.Match.CurrentPlayer()

	if event, ok := ctx.Event.(*match.AttackConfirmed); ok {

		creature, err := p.Player.GetCard(event.CardID, match.BATTLEZONE)

		if err != nil {
			return
		}

		if creature.Civ != civ.Water || creature.ID == card.ID {
			return
		}

		creature.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)

	}

}
