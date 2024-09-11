package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func ArmoredDecimatorValkaizer(c *match.Card) {

	c.Name = "Armored Decimator Valkaizer"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Armored Decimator Valkaizer: You may select 1 opponent's creature with 4000 or less power and destroy it.",
				0,
				1,
				true,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 4000 },
			).Map(func(x *match.Card) {

				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Armored Decimator Valkaizer", x.Name))
			})

		}
	})

}

func MigasaAdeptOfChaos(c *match.Card) {

	c.Name = "Migasa, Adept of Chaos"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		ctx.Match.Chat("Server", fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))
		creatures := match.Filter(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 fire creature from your battlezone that will gain double breaker", 1, 1, false, func(x *match.Card) bool { return x.Civ == civ.Fire && x.ID != card.ID })
		for _, creature := range creatures {
			if creature.Civ == civ.Fire {
				creature.AddCondition(cnd.DoubleBreaker, true, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given double breaker power by %s until end of turn", creature.Name, card.Name))
			}
		}
	}

	c.Use(fx.Creature, fx.TapAbility)
}

func ChoyaTheUnheeding(c *match.Card) {

	c.Name = "Choya, the Unheeding"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker1000, func(card *match.Card, ctx *match.Context) {

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
