package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Gigazoul(c *match.Card) {

	c.Name = "Gigazoul"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.CantAttackIfOpHasNoShields)

}

func Gigakail(c *match.Card) {
	c.Name = "Gigakail"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ConditionalSlayer(func(target *match.Card) bool {
		return target.Civ == civ.Nature || target.Civ == civ.Light
	}))
}

func GigalingQ(c *match.Card) {

	c.Name = "Gigaling Q"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera, family.Survivor}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
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
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})
		})
	}))

}
