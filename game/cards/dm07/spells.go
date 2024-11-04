package dm07

import (
	"duel-masters/game/civ"
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
