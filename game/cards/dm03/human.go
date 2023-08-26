package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ExplosiveDudeJoe ...
func ExplosiveDudeJoe(c *match.Card) {

	c.Name = "Explosive Dude Joe"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}

// MuramasaDukeOfBlades ...
func MuramasaDukeOfBlades(c *match.Card) {

	c.Name = "Muramasa, Duke of Blades"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		creatures := match.Filter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Muramasa, Duke of Blades: Select 1 of your opponent's creatures with power 2000 or less and destroy it",
			1,
			1,
			true,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
		)

		for _, creature := range creatures {
			ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
		}
	}))

}
