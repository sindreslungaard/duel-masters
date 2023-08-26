package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LeQuistTheOracle(c *match.Card) {
	c.Name = "Le Quist, the Oracle"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	toTap := ""

	filter := func(x *match.Card) bool { return x.Civ == civ.Fire || x.Civ == civ.Darkness }
	ability := func(ctx *match.Context, blockers []*match.Card) (newblockers []*match.Card) {
		cards := make(map[string][]*match.Card)
		cards["Your creatures"] = fx.FindFilter(c.Player, match.BATTLEZONE, filter)
		cards["Opponent's creatures"] = fx.FindFilter(ctx.Match.Opponent(c.Player), match.BATTLEZONE, filter)

		fx.SelectMultipart(
			c.Player,
			ctx.Match,
			cards,
			"Le Quist, the Oracle: Select a card to tap or close to cancel",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			toTap = x.ID
		}).Or(func() {
			toTap = ""
		})

		b := []*match.Card{}

		for _, x := range blockers {
			if x.ID == toTap {
				continue
			}

			b = append(b, x)
		}

		return b
	}

	c.Use(func(card *match.Card, ctx *match.Context) {
		switch event := ctx.Event.(type) {

		case *match.AttackPlayer:
			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {
				event.Blockers = ability(ctx, event.Blockers)
			})

		case *match.AttackCreature:
			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {
				event.Blockers = ability(ctx, event.Blockers)
			})

		case *match.AttackConfirmed:
			if event.CardID != card.ID || toTap == "" {
				return
			}

			fx.FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, func(x *match.Card) bool {
				return x.ID == toTap
			}).Map(func(x *match.Card) {
				x.Tapped = true
			})

			toTap = ""

		}
	}, fx.Creature)

}
