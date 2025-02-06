package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func LegionnaireLizard(c *match.Card) {

	c.Name = "Legionnaire Lizard"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.DuneGecko}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability to give creature  \"Speed Attacker this turn\"", card.Player.Username(), card.Name))
		creatures := fx.SelectFilter(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Speed attacker\"", 1, 1, false, func(x *match.Card) bool { return x.ID != card.ID }, false)
		for _, creature := range creatures {

			creature.AddCondition(cnd.SpeedAttacker, true, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given \"Speed Attacker by %s\"", creature.Name, card.Name))

		}
	}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.TapAbility)
}

func BadlandsLizard(c *match.Card) {

	c.Name = "Badlands Lizard"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.DuneGecko}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker3000, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if !event.Blocked || event.Attacker != card {
				return
			}

			ctx.InterruptFlow()

			event.Attacker.Tapped = true
			event.Defender.Tapped = true

		}
	})
}
