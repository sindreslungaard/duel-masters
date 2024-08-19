package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ZorvazTheBonecrusher(c *match.Card) {

	c.Name = "Zorvaz, the Bonecrusher"
	c.Power = 8000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers, fx.Suicide)
}

func VileMulderWingOfTheVoid(c *match.Card) {

	c.Name = "Vile Mulder, Wing of the Void"
	c.Power = 7000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.CantAttackCreatures, fx.Doublebreaker, fx.Suicide)
}

func DaidalosGeneralOfFury(c *match.Card) {

	c.Name = "Daidalos, General of Fury"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Daidalos, General Of Fury: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

		for _, creature := range creatures {
			if creature.ID == card.ID {
				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				ctx.InterruptFlow()
			}
			if creature.ID != card.ID {
				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
			}
		}

	}))
}

func GnarvashMerchantOfBlood(c *match.Card) {

	c.Name = "Gnarvash, Merchant of Blood"
	c.Power = 8000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if len(fx.Find(card.Player, match.BATTLEZONE)) == 1 {

			ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
		}

	}))
}
