package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MizoyTheOracle ...
func MizoyTheOracle(c *match.Card) {

	c.Name = "Mizoy, the Oracle"
	c.Power = 2500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.LightBringer}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		matching := func(x *match.Card) bool { return x.HasCiv(civ.Darkness) || x.HasCiv(civ.Fire) }

		// "In the battle zone" is both sides of it.
		cards := map[string][]*match.Card{
			"Your creatures":            fx.FindFilter(card.Player, match.BATTLEZONE, matching),
			"Your opponent's creatures": fx.FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, matching),
		}

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: You may choose a darkness or fire creature in the battle zone and tap it.", card.Name),
			1,
			1,
			true,
		).Map(func(chosen *match.Card) {
			chosen.Tapped = true
			ctx.Match.ReportActionInChat(chosen.Player, fmt.Sprintf("%s was tapped by %s", chosen.Name, card.Name))
		})
	}))

}
