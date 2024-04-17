package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KizarBasikuTheOutrageous(c *match.Card) {

	c.Name = "Kizar Basiku, the Outrageous"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			if match.ContainerHas(ctx.Match.Opponent(card.Player), match.MANAZONE, func(x *match.Card) bool { return x.Civ == civ.Fire }) {
				card.AddCondition(cnd.CantBeBlocked, true, card.ID)
			}

		}
		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			if match.ContainerHas(ctx.Match.Opponent(card.Player), match.MANAZONE, func(x *match.Card) bool { return x.Civ == civ.Fire }) {
				card.AddCondition(cnd.CantBeBlocked, true, card.ID)
			}

		}

	})

}

func RomVizierOfTendrils(c *match.Card) {

	c.Name = "Rom, Vizier of Tendrils"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Rom, Vizier of Tendrils: Select 1 of your opponent's creature and tap it. Close to not tap any creatures.", 1, 1, true)

				for _, creature := range creatures {
					creature.Tapped = true
				}

			}

		}

	})

}
