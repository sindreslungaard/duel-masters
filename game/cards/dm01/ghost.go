package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DarkRavenShadowOfGrief ...
func DarkRavenShadowOfGrief(c *match.Card) {

	c.Name = "Dark Raven, Shadow of Grief"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker)

}

// MaskedHorrorShadowOfScorn ...
func MaskedHorrorShadowOfScorn(c *match.Card) {

	c.Name = "Masked Horror, Shadow of Scorn"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.OpponentDiscardsRandomCard))

}

// NightMasterShadowOfDecay ...
func NightMasterShadowOfDecay(c *match.Card) {

	c.Name = "Night Master, Shadow of Decay"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker)

}

// BlackFeatherShadowOfRage ...
func BlackFeatherShadowOfRage(c *match.Card) {

	c.Name = "Black Feather, Shadow of Rage"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			ctx.ScheduleAfter(func() {

				fx.Select(
					card.Player,
					ctx.Match,
					card.Player,
					match.BATTLEZONE,
					"Select one of your creatures that will be destroyed",
					1,
					1,
					false,
				).Map(func(x *match.Card) {
					ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				})

			})
		}
	})

}
