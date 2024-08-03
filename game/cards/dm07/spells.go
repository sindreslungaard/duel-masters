package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func HopelessVortex(c *match.Card) {
	c.Name = "Hopeless Vortex"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Destroy one of your opponent's creatures", 1, 1, false)

		for _, creature := range creatures {

			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

		}
	}))
}

func FreezingIcehammer(c *match.Card) {

	c.Name = "Freezing Icehammer"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Select 1 of your opponent's water or darkness creatures and put it in their manazone",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.Civ == civ.Water || x.Civ == civ.Darkness },
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", x.Name, x.Player.Username()))

		})

	}))

}

func FruitOfEternity(c *match.Card) {

	c.Name = "Fruit Of Eternity"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx *match.Context, exit func()) {

			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card.Player == card.Player {

				//could probably think about making an option to choose when the creature with ability to return to hand dies
				ctx.InterruptFlow()
				card.Player.MoveCard(event.Card.ID, match.BATTLEZONE, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed and moved to %s's manazone", event.Card.Name, event.Card.Player.Username()))
			}

			// remove persistent effect when turn ends
			_, ok := ctx.Event.(*match.EndStep)
			if ok {
				exit()
			}
		})
	}))
}

func VacuumGel(c *match.Card) {

	c.Name = "Vacuum Gel"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Destroy one of your opponent's untapped light or untapped nature creatures",
			1,
			1,
			false,
			func(x *match.Card) bool { return !x.Tapped && x.Civ == civ.Light || x.Civ == civ.Nature },
		).Map(func(x *match.Card) {

			ctx.Match.Destroy(x, card, match.DestroyedBySpell)

		})
	}))
}

// func JusticeJamming(c *match.Card) {

// 	c.Name = "Justice Jamming"
// 	c.Civ = civ.Light
// 	c.ManaCost = 3
// 	c.ManaRequirement = []string{civ.Light}

// 	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

// 		creatures := make(map[string][]*match.Card)

// 		creatures["Darkness creatures"] = getDarknessCreatures(card, ctx)
// 		creatures["Fire creatures"] = getFireCreatures(card, ctx)

// 		defer ctx.Match.CloseAction(card.Player)

// 		fx.SelectMultipart(card.Player, ctx.Match, creatures, "Select all fire creatures or all darkness creatures to tap:", len(creatures)+1, len(creatures)+1, false).Map(func(creatures *match.Card) {

// 			creatures.Tapped = true
// 			ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by Justice Jamming", creatures.Name))
// 		})

// 	}))

// }

// func getDarknessCreatures(card *match.Card, ctx *match.Context) fx.CardCollection {

// 	darknessCreatures := fx.FindFilter(
// 		card.Player,
// 		match.BATTLEZONE,
// 		func(x *match.Card) bool { return x.Civ == civ.Darkness },
// 	)

// 	darknessCreatures = append(darknessCreatures,

// 		fx.FindFilter(
// 			ctx.Match.Opponent(card.Player),
// 			match.BATTLEZONE,
// 			func(x *match.Card) bool { return x.Civ == civ.Darkness },
// 		)...,
// 	)

// 	return darknessCreatures
// }

// func getFireCreatures(card *match.Card, ctx *match.Context) fx.CardCollection {

// 	fireCreatures := fx.FindFilter(
// 		card.Player,
// 		match.BATTLEZONE,
// 		func(x *match.Card) bool { return x.Civ == civ.Fire },
// 	)

// 	fireCreatures = append(fireCreatures,

// 		fx.FindFilter(
// 			ctx.Match.Opponent(card.Player),
// 			match.BATTLEZONE,
// 			func(x *match.Card) bool { return x.Civ == civ.Fire },
// 		)...,
// 	)

// 	return fireCreatures
// }

func MiraclePortal(c *match.Card) {

	c.Name = "Miracle Portal"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Choose one of your creatures in the battle zone. This turn, it can't be blocked and you ignore any effects that would prevent that creature from attacking your opponent.",
			1,
			1,
			false,
		).Map(func(c *match.Card) {

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				if c.HasCondition(cnd.SummoningSickness) {
					c.RemoveCondition(cnd.SummoningSickness)
					c.AddCondition(cnd.CantAttackCreatures, nil, card.ID)
				}

				c.RemoveCondition(cnd.CantAttackPlayers)
				c.AddCondition(cnd.CantBeBlocked, true, card)

				if event, ok := ctx2.Event.(*match.AttackCreature); ok {

					// Is this event for me or someone else?
					if event.CardID != c.ID || !c.HasCondition(cnd.CantAttackCreatures) {
						return
					}

					ctx2.Match.WarnPlayer(c.Player, fmt.Sprintf("%s can't attack creatures", c.Name))

					ctx2.InterruptFlow()

				}

				// remove persistent effect when turn ends
				_, ok := ctx2.Event.(*match.EndStep)
				if ok {
					exit()
				}

			})
		})
	}))
}

func VenomCharger(c *match.Card) {
	c.Name = "Venom Charger"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Slayer\"", 1, 1, false)

		for _, creature := range creatures {

			creature.AddCondition(cnd.Slayer, true, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given slayer", creature.Name))

		}

	}))
}

func EnergyCharger(c *match.Card) {
	c.Name = "Energy Charger"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Power Attacker +2000\"", 1, 1, false)

		for _, creature := range creatures {

			creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000", creature.Name))

		}
	}))
}

func RiptideCharger(c *match.Card) {
	c.Name = "Riptide Charger"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		cards["Your creatures"] = fx.Find(card.Player, match.BATTLEZONE)
		cards["Opponent's creatures"] = fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose 1 creature in the battlezone that will be sent to its owner's hand", card.Name),
			1,
			1,
			true).Map(func(creature *match.Card) {
			creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
		})

	}))
}

func LightningCharger(c *match.Card) {
	c.Name = "Lightning Charger"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponents creatures that will be tapped", 1, 1, false)

		for _, creature := range creatures {

			creature.Tapped = true
			ctx.Match.Chat("Server", creature.Name+" was tapped")

		}

	}))
}

func MulchCharger(c *match.Card) {
	c.Name = "Mulch Charger"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 of your creatures and put it in your manazone", 1, 1, false)

		for _, creature := range creatures {

			creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.MANAZONE)
			creature.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", creature.Name, creature.Player.Username()))

		}

	}))
}
