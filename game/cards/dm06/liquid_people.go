package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrystalJouster(c *match.Card) {

	c.Name = "Crystal Jouster"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.When(fx.WouldBeDestroyed, fx.ReturnToHand))

}

func AquaRider(c *match.Card) {

	c.Name = "Aqua Rider"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	isBlocker := false

	c.Use(fx.Creature,
		fx.When(ReceiveBlockerWhenOpponentPlaysCreatureOrSpell,
			func(c *match.Card, ctx *match.Context) { isBlocker = true }),
		fx.When(fx.EndOfTurn,
			func(c *match.Card, ctx *match.Context) {
				isBlocker = false
				c.RemoveConditionBySource(c.ID)
			}),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return isBlocker },
			func(c *match.Card, ctx *match.Context) {
				fx.ForceBlocker(c, ctx, c.ID)
			}),
	)
}
