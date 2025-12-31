package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func JunkatzRabidDoll(c *match.Card) {

	c.Name = "Junkatz, Rabid Doll"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.DeathPuppet}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)
}

func LupaPoisonTippedDoll(c *match.Card) {

	c.Name = "Lupa, Poison-Tipped Doll"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.DeathPuppet}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 creature from your battlezone that will get 'Slayer'", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.Slayer, 1, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given 'slayer' by %s until end of turn", x.Name, card.Name))
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}
