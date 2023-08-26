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
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {

		ctx.Match.Chat("Server", fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))
		creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will get 'slayer'", 1, 1, false)
		for _, creature := range creatures {

			creature.AddCondition(cnd.Slayer, 1, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given 'slayer' by %s until end of turn", creature.Name, card.Name))

			card.Tapped = true
		}
	}))
}
