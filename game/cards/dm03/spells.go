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

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			cards := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "Select 1 card from your mana zone that will be sent to your hand", 1, 1, false, func(x *match.Card) bool { return true })

			for _, card := range cards {

				card.Player.MoveCard(card.ID, match.MANAZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", card.Player.Username(), card.Name))

			}

			card.Player.MoveCard(card.ID, match.HAND, match.MANAZONE)
		}
	})
}

// LogicSphere ...
func LogicSphere(c *match.Card) {

	c.Name = "Logic Sphere"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			spells := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "Select 1 spell from your mana zone that will be sent to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

			for _, spell := range spells {

				card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", spell.Player.Username(), spell.Name))

			}
		}
	})
}

// SundropArmor ...
func SundropArmor(c *match.Card) {

	c.Name = "Sundrop Armor"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.HAND,
				"Select 1 card from your hand that will be put as a shield",
				1,
				1,
				false,
				func(c *match.Card) bool { return c.ID != card.ID },
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.HAND, match.SHIELDZONE)
			})

		}
	})
}

// FloodValve ...
func FloodValve(c *match.Card) {

	c.Name = "Flood Valve"
	c.Civ = civ.Water
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.MANAZONE,
				"Select 1 card from your mana zone that will be returned to your hand",
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.HAND, card)
			})

		}
	})
}

// LiquidScope ...
func LiquidScope(c *match.Card) {

	c.Name = "Liquid Scope"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			shields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

			if err != nil {
				return
			}

			hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

			if err != nil {
				return
			}

			ids := make([]string, 0)

			for _, s := range shields {
				ids = append(ids, s.ImageID)
			}

			ctx.Match.ShowCards(
				card.Player,
				"Your opponent's shields:",
				ids,
			)

			ids = make([]string, 0)

			for _, s := range hand {
				ids = append(ids, s.ImageID)
			}

			ctx.Match.ShowCards(
				card.Player,
				"Your opponent's hand:",
				ids,
			)
		}
	})
}

// PsychicShaper ...
func PsychicShaper(c *match.Card) {

	c.Name = "Psychic Shaper"
	c.Civ = civ.Water
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			cards := card.Player.PeekDeck(4)

			for _, toMove := range cards {

				if toMove.Civ == civ.Water {
					card.Player.MoveCard(toMove.ID, match.DECK, match.HAND)
					ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the hand from the top of their deck", card.Player.Username(), toMove.Name))
				} else {
					card.Player.MoveCard(toMove.ID, match.DECK, match.GRAVEYARD)
					ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the graveyard from the top of their deck", card.Player.Username(), toMove.Name))
				}
			}
		}
	})
}

// EldritchPoison ...
func EldritchPoison(c *match.Card) {

	c.Name = "Eldritch Poison"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 of your darkness creatures that will be destroyed", 0, 1, true, func(x *match.Card) bool { return x.Civ == civ.Darkness })

			if len(creatures) > 0 {

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
				}

				creatures = match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "Select 1 of your creatures from your mana zone that will be returned to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) })

				for _, creature := range creatures {
					card.Player.MoveCard(creature.ID, match.MANAZONE, match.HAND)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their mana zone", creature.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				}
			}
		}
	})
}

// GhastlyDrain ...
func GhastlyDrain(c *match.Card) {

	c.Name = "Ghastly Drain"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			shields, err := card.Player.Container(match.SHIELDZONE)

			if err != nil {
				return
			}
			nrShields := len(shields)

			fx.SelectBackside(
				card.Player,
				ctx.Match, card.Player,
				match.SHIELDZONE,
				"Select any number of shields and move them into your hand",
				0,
				nrShields,
				true,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.HAND, card)
			})
		}
	})
}

// SnakeAttack ...
func SnakeAttack(c *match.Card) {

	c.Name = "Snake Attack"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectBackside(
				card.Player,
				ctx.Match, card.Player,
				match.SHIELDZONE,
				"Select one shield and send it to graveyard",
				1,
				1,
				true,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.GRAVEYARD, card)
			})

			creatures, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			for _, creature := range creatures {

				creature.AddCondition(cnd.DoubleBreaker, nil, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given double breaker until the end of the turn", creature.Name))
			}
		}
	})
}

// BlazeCannon ...
func BlazeCannon(c *match.Card) {

	c.Name = "Blaze Cannon"
	c.Civ = civ.Fire
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Fire }) {
				return
			}

			creatures, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			for _, creature := range creatures {

				creature.AddCondition(cnd.PowerAttacker, 4000, card.ID)
				creature.AddCondition(cnd.DoubleBreaker, nil, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +4000 and double breaker until the end of the turn", creature.Name))

			}
		}
	})
}

// SearingWave ...
func SearingWave(c *match.Card) {

	c.Name = "Searing Wave"
	c.Civ = civ.Fire
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			ctx.ScheduleAfter(func() {

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

				creatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

				if err != nil {
					return
				}

				toDestroy := make([]*match.Card, 0)

				for _, creature := range creatures {

					if ctx.Match.GetPower(creature, false) <= 3000 {
						toDestroy = append(toDestroy, creature)
					}
				}

				for _, creature := range toDestroy {
					ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
				}

			})
		}
	})
}

// VolcanicArrows ...
func VolcanicArrows(c *match.Card) {

	c.Name = "Volcanic Arrows"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

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

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Volcanic arrows: Select 1 of your opponent's creatures with power 6000 or less and destroy it",
				0,
				1,
				false,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 6000 },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedBySpell)
			})
		}
	})
}

// AuroraOfReversal ...
func AuroraOfReversal(c *match.Card) {

	c.Name = "Aurora of Reversal"
	c.Civ = civ.Nature
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			shields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

			if err != nil {
				return
			}

			fx.SelectBackside(
				card.Player,
				ctx.Match,
				card.Player,
				match.SHIELDZONE,
				"Select any number of shields that will be sent to mana zone",
				0,
				len(shields),
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.MANAZONE, card)
			})
		}
	})
}

// ManaNexus ...
func ManaNexus(c *match.Card) {

	c.Name = "Mana Nexus"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.MANAZONE,
				"Select a card from the mana zone that will be put as a shield",
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.SHIELDZONE, card)
			})
		}
	})
}

// RoarOfTheEarth ...
func RoarOfTheEarth(c *match.Card) {

	c.Name = "Roar of the Earth"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.MANAZONE,
				"Select a creature that costs 6 ore more mana from the mana zone that will be put in your hand",
				1,
				1,
				false,
				func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && x.ManaCost >= 6 },
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.HAND, card)
			})
		}
	})
}
