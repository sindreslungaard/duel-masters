package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TwinCannonSkyterror ...
func TwinCannonSkyterror(c *match.Card) {

	c.Name = "Twin-Cannon Skyterror"
	c.Power = 7000
	c.Civ = civ.Fire
	c.Family = family.ArmoredWyvern
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.Doublebreaker)

}

// BladerushSkyterrorQ ...
func BladerushSkyterrorQ(c *match.Card) {

	c.Name = "Bladerush Skyterror Q"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = family.ArmoredWyvern
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(card2 *match.Card, ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return

			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			})

		})

	}))

}
