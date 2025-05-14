package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BoomerangComet ...
func BoomerangComet(c *match.Card) {

	c.Name = "Boomerang Comet"
	c.Civ = civ.Light
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 card from your mana zone that will be sent to your hand", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s retrieved %s from the mana zone to their hand", x.Player.Username(), x.Name))
		})

		card.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, card.ID)
	}))

}

// LogicSphere ...
func LogicSphere(c *match.Card) {

	c.Name = "Logic Sphere"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 spell from your mana zone that will be sent to your hand", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s retrieved %s from the mana zone to their hand", x.Player.Username(), x.Name))
		})
	}))

}

// SundropArmor ...
func SundropArmor(c *match.Card) {

	c.Name = "Sundrop Armor"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s: Select 1 card from your hand that will be put as a shield", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.ID != card.ID },
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.HAND, match.SHIELDZONE, card.ID)
		})
	}))

}

// FloodValve ...
func FloodValve(c *match.Card) {

	c.Name = "Flood Valve"
	c.Civ = civ.Water
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.ReturnMyCardFromMZToHand))
}

// LiquidScope ...
func LiquidScope(c *match.Card) {

	c.Name = "Liquid Scope"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		shields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if err != nil {
			return
		}

		hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil {
			return
		}

		ids := make([]string, 0)
		cardsCount := len(hand)
		shieldsCount := len(shields)

		for _, card := range append(hand, shields...) {
			ids = append(ids, card.ImageID)
		}

		message := fmt.Sprintf("Your opponent's hand is shown first(%d), followed by their shields(%d):", cardsCount, shieldsCount)

		ctx.Match.ShowCards(
			card.Player,
			message,
			ids,
		)
	}))

}

// PsychicShaper ...
func PsychicShaper(c *match.Card) {

	c.Name = "Psychic Shaper"
	c.Civ = civ.Water
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		cards := card.Player.PeekDeck(4)

		for _, toMove := range cards {

			if toMove.Civ == civ.Water {
				card.Player.MoveCard(toMove.ID, match.DECK, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into the hand from the top of their deck", card.Player.Username(), toMove.Name))
			} else {
				card.Player.MoveCard(toMove.ID, match.DECK, match.GRAVEYARD, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into the graveyard from the top of their deck", card.Player.Username(), toMove.Name))
			}
		}
	}))
}

// EldritchPoison ...
func EldritchPoison(c *match.Card) {

	c.Name = "Eldritch Poison"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"%s: You may destroy one of your darkness creatures. If you do, return a creature from your mana zone to your hand.",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.Civ == civ.Darkness },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
			fx.ReturnCreatureFromManazoneToHand(card, ctx)
		})
	}))

}

// GhastlyDrain ...
func GhastlyDrain(c *match.Card) {

	c.Name = "Ghastly Drain"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		shields, err := card.Player.Container(match.SHIELDZONE)

		if err != nil {
			return
		}

		fx.SelectBackside(
			card.Player,
			ctx.Match, card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select any number of shields and move them into your hand", card.Name),
			1,
			len(shields),
			true,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.HAND, card)
		})
	}))

}

// SnakeAttack ...
func SnakeAttack(c *match.Card) {

	c.Name = "Snake Attack"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectBackside(
			card.Player,
			ctx.Match, card.Player,
			match.SHIELDZONE,
			"Select one shield and send it to graveyard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range creatures {
			creature.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given double breaker until the end of the turn", creature.Name))
		}

	}))
}

// BlazeCannon ...
func BlazeCannon(c *match.Card) {

	c.Name = "Blaze Cannon"
	c.Civ = civ.Fire
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Fire }) {
			return
		}

		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAttacker, 4000, card.ID)
			x.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given \"Power Attacker +4000\" and \"Double Breaker\" until the end of the turn", x.Name))
		})
	}))

}

// SearingWave ...
func SearingWave(c *match.Card) {

	c.Name = "Searing Wave"
	c.Civ = civ.Fire
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.ScheduleAfter(func() {
			fx.SelectBackside(
				card.Player,
				ctx.Match,
				card.Player,
				match.SHIELDZONE,
				fmt.Sprintf("%s: Select 1 shield and send it to your graveyard", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.GRAVEYARD, card)
			})

			fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 3000 },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedBySpell)
			})
		})
	}))

}

// VolcanicArrows ...
func VolcanicArrows(c *match.Card) {

	c.Name = "Volcanic Arrows"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select 1 shield and send it to graveyard", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's creatures with power 6000 or less and destroy it", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 6000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})
	}))
}

// AuroraOfReversal ...
func AuroraOfReversal(c *match.Card) {

	c.Name = "Aurora of Reversal"
	c.Civ = civ.Nature
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		shields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if err != nil {
			return
		}

		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select any number of shields that will be sent to your mana zone", card.Name),
			1,
			len(shields),
			true,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.MANAZONE, card)
		})
	}))

}

// ManaNexus ...
func ManaNexus(c *match.Card) {

	c.Name = "Mana Nexus"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select a card from the mana zone that will be put as a shield", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.SHIELDZONE, card)
		})
	}))

}

// RoarOfTheEarth ...
func RoarOfTheEarth(c *match.Card) {

	c.Name = "Roar of the Earth"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select a creature that costs 6 or more mana from your mana zone that will be put in your hand", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && x.ManaCost >= 6 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.HAND, card)
		})
	}))

}
