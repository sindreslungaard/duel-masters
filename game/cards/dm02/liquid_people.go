package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CrystalLancer ...
func CrystalLancer(c *match.Card) {

	c.Name = "Crystal Lancer"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature, fx.Evolution, fx.Doublebreaker)

}

// CrystalPaladin ...
func CrystalPaladin(c *match.Card) {

	c.Name = "Crystal Paladin"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s battle zone to their hand by Crystal Paladin", x.Name, x.Player.Username()))
		})

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s battle zone to their hand by Crystal Paladin", x.Name, x.Player.Username()))
		})

	}))

}

// AquaBouncer ...
func AquaBouncer(c *match.Card) {

	c.Name = "Aqua Bouncer"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		cards["Your creatures"] = fx.Find(card.Player, match.BATTLEZONE)
		cards["Opponent's creatures"] = fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			"Choose a card in the battle zone and return it to its owner's hand, or close to cancel",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s hand by %s", x.Name, x.Player.Username(), card.Name))
		})

	}))

}

// AquaShooter ...
func AquaShooter(c *match.Card) {

	c.Name = "Aqua Shooter"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker)

}
