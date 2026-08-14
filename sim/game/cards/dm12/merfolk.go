package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// FunkyWizard ...
func FunkyWizard(c *match.Card) {

	c.Name = "Funky Wizard"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// Both players are asked, its controller first.
		for _, player := range []*match.Player{card.Player, ctx.Match.Opponent(card.Player)} {
			if fx.BinaryQuestion(player, ctx.Match, fmt.Sprintf("%s's effect: Do you want to draw a card?", card.Name)) {
				player.DrawCards(1)
			}
		}
	}))

}

// FranticChieftain ...
func FranticChieftain(c *match.Card) {

	c.Name = "Frantic Chieftain"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// Frantic Chieftain costs 2, so it is a legal choice for its own
		// effect when nothing cheaper is on the board.
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Return one of your creatures that costs 4 or less to your hand.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.ManaCost <= 4 },
			false,
		).Map(func(chosen *match.Card) {
			moved, err := card.Player.MoveCard(chosen.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err == nil && moved.Zone == match.HAND {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand by %s", chosen.Name, card.Player.Username(), card.Name))
			}
		})
	}))

}

// WilyCarpenter ...
func WilyCarpenter(c *match.Card) {

	c.Name = "Wily Carpenter"
	c.Power = 1000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.DrawUpTo2(card, ctx)

		// The discard is not conditional on the draw, so a player who declines
		// still pays two cards.
		fx.DiscardOwnXCards(2)(card, ctx)
	}))

}

// SeaMutantDormel ...
func SeaMutantDormel(c *match.Card) {

	c.Name = "Sea Mutant Dormel"
	c.Power = 4000
	c.Civs = []string{civ.Water, civ.Darkness}
	c.Family = []string{family.Merfolk, family.Hedrian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}
