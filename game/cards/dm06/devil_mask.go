package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func GrinningAxeTheMonstrosity(c *match.Card) {

	c.Name = "Grinning Axe, the Monstrosity"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.DevilMask}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)
}

func SkullcutterSwarmLeader(c *match.Card) {

	c.Name = "Skullcutter, Swarm Leader"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DevilMask}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		if len(fx.Find(card.Player, match.BATTLEZONE)) == 1 {
			if card.Zone != match.BATTLEZONE {
				return
			}
			ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
		}

	}))
}
