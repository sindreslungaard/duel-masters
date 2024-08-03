package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func VessTheOracle(c *match.Card) {

	c.Name = "Vess, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)
}

func YulukTheOracle(c *match.Card) {

	c.Name = "Yuluk, the Oracle"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	spellcast := false

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.SpellCast); ok {
			spellcast = true
		}

		_, ok := ctx.Event.(*match.EndStep)
		if ok {
			spellcast = false
		}

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {
			if event.CardID == card.ID {
				if !spellcast {
					ctx.InterruptFlow()
					ctx.Match.WarnPlayer(card.Player, "You can summon this creature only if you have cast a spell this round")
					return
				}
			}

		}
	})
}
