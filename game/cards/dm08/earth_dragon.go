package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TerradragonRegarion ...
func TerradragonRegarion(c *match.Card) {

	c.Name = "Terradragon Regarion"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000)
}

// TerradragonGamiratar ...
func TerradragonGamiratar(c *match.Card) {

	c.Name = "Terradragon Gamiratar"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			fmt.Sprintf("%s's effect: You may choose a creature from your hand and put it into your battlezone.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.BATTLEZONE, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was put into the battlezone from %s's hand.", x.Name, ctx.Match.Opponent(card.Player).Username()))
		})
	}))
}
