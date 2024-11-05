package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func JusticeJamming(c *match.Card) {

	c.Name = "Justice Jamming"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		option := fx.MultipleChoiceQuestion(
			card.Player,
			ctx.Match,
			"Tap all creatures of a civ type",
			[]string{"All darkness", "All fire"},
		)
		var civToTap string
		if option == 0 {
			civToTap = civ.Darkness
		} else {
			civToTap = civ.Fire
		}

		fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {
			if x.Civ == civToTap && !x.Tapped {
				x.Tapped = true
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was tapped", x.Name))
			}
		})
		fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE).Map(func(x *match.Card) {
			if x.Civ == civToTap && !x.Tapped {
				x.Tapped = true
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was tapped", x.Name))
			}
		})
	}))

}

// ApocalypseVise ...
func ApocalypseVise(c *match.Card) {

	c.Name = "Apocalypse Vise"
	c.Civ = civ.Fire
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		correctSelection := false

		for !correctSelection {

			creatures := fx.Select(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Destroy any number of your opponent's creatures that have total power 8000 or less.",
				0,
				16,
				false,
			)

			totalPower := 0

			for _, creature := range creatures {
				totalPower += ctx.Match.GetPower(creature, false)
			}

			if totalPower <= 8000 {

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
				}

				correctSelection = true

			}

		}

	}))

}

func HopelessVortex(c *match.Card) {
	c.Name = "Hopeless Vortex"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.DestroyOpCreature))
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
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's manazone", x.Name, x.Player.Username()))
		})

	}))

}

func FruitOfEternity(c *match.Card) {

	c.Name = "Fruit of Eternity"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx *match.Context, exit func()) {

			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card.Player == card.Player {
				ctx.InterruptFlow()
				card.Player.MoveCard(event.Card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was destroyed but moved to manazone because of %s", event.Card.Name, card.Name))
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
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})
	}))
}

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
				c.AddCondition(cnd.CantBeBlocked, nil, card)

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

		fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE,
			"Select 1 creature from your battlezone that will gain \"Slayer\"", 1, 1, false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.Slayer, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given slayer", x.Name))
		})
	}))
}

func EnergyCharger(c *match.Card) {
	c.Name = "Energy Charger"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE,
			"Select 1 creature from your battlezone that will gain +2000 power", 1, 1, false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 2000, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given +2000", x.Name))
		})
	}))
}
