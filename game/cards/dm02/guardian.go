package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LadiaBaleTheInspirational ...
func LadiaBaleTheInspirational(c *match.Card) {

	c.Name = "Ladia Bale, the Inspirational"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.Evolution, fx.Doublebreaker)

}

// PhalEegaDawnGuardian ...
func PhalEegaDawnGuardian(c *match.Card) {

	c.Name = "Phal Eega, Dawn Guardian"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			"Phal Eega, Dawn Guardian: Select a spell to return from your graveyard to your hand.",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
		).Map(func(x *match.Card) {
			if x.ID == card.ID {
				return
			}
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand from their graveyard by Phal Eega, Dawn Guardian", x.Name, x.Player.Username()))
		})

	}))

}

// ResoPacosClearSkyGuardian ...
func ResoPacosClearSkyGuardian(c *match.Card) {

	c.Name = "Reso Pacos, Clear Sky Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// LarbaGeerTheImmaculate ...
func LarbaGeerTheImmaculate(c *match.Card) {

	c.Name = "Larba Geer, the Immaculate"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {

			fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
			).Map(func(x *match.Card) {
				x.Tapped = true
				ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by %s", x.Name, card.Name))
			})

		})

	}))

}
