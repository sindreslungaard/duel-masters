package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func JusticeJamming(c *match.Card) {

	c.Name = "Justice Jamming"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		option := fx.MultipleChoiceQuestion(
			card.Player,
			ctx.Match,
			"Tap all creatures of a civ type",
			[]string{"All darkness", "All fire"},
		)
		var civToTap string
		if option == 0 {
			civToTap = civ.Darkness
		} else {
			civToTap = civ.Fire
		}

		fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {
			if x.Civ == civToTap && !x.Tapped {
				x.Tapped = true
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was tapped", x.Name))
			}
		})
		fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE).Map(func(x *match.Card) {
			if x.Civ == civToTap && !x.Tapped {
				x.Tapped = true
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was tapped", x.Name))
			}
		})
	}))

}
