package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaRanger ...
func AquaRanger(c *match.Card) {

	c.Name = "Aqua Ranger"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.CantBeBlocked, fx.When(fx.WouldBeDestroyed, fx.ReturnToHand))

}

// AquaGrappler ...
func AquaGrappler(c *match.Card) {

	c.Name = "Aqua Grappler"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		numberOfMyOtherTappedCreatures := len(fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(card2 *match.Card) bool {
				return card2.Tapped && card2.ID != card.ID
			},
		))

		fx.DrawUpto(card, ctx, numberOfMyOtherTappedCreatures)
	}))

}
