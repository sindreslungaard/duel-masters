package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// EnchantedSoil ...
func EnchantedSoil(c *match.Card) {

	c.Name = "Enchanted Soil"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "Enchanted Soil: Select 2 creatures from your graveyard and put it in your manazone", 0, 2, true, func(x *match.Card) bool {
				return x.HasCondition(cnd.Creature)
			}).Map(func(c *match.Card) {
				card.Player.MoveCard(c.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone by %s", c.Name, c.Player.Username(), card.Name))
			})

		}
	})
}

// SchemingHands ...
func SchemingHands(c *match.Card) {

	c.Name = "Scheming Hands"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.HAND, "Scheming Hands: Discard a card from your opponent's hand", 0, 1, false).Map(func(c *match.Card) {
				c.Player.MoveCard(c.ID, match.HAND, match.GRAVEYARD, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's graveyard by %s", c.Name, c.Player.Username(), card.Name))
			})

		}
	})
}

// CyclonePanic ...
func CyclonePanic(c *match.Card) {

	c.Name = "Cyclone Panic"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		// p1
		cards1 := fx.Find(card.Player, match.HAND)
		n1 := 0

		for _, c1 := range cards1 {
			if c1.ID != card.ID {
				card.Player.MoveCard(c1.ID, match.HAND, match.DECK, card.ID)
				n1 += 1
			}
		}

		card.Player.ShuffleDeck()
		card.Player.DrawCards(n1)

		// p2
		cards2 := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)
		n2 := len(cards2)

		for _, c2 := range cards2 {
			ctx.Match.Opponent(card.Player).MoveCard(c2.ID, match.HAND, match.DECK, card.ID)
		}

		ctx.Match.Opponent(card.Player).ShuffleDeck()
		ctx.Match.Opponent(card.Player).DrawCards(n2)

		ctx.Match.Chat("Server", "Cyclone Panic: Both players shuffled their deck and replaced the cards in their hand with new ones")

	}))
}

// GlorySnow ...
func GlorySnow(c *match.Card) {

	c.Name = "Glory Snow"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			if len(fx.Find(card.Player, match.MANAZONE)) < len(fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)) {

				cards := card.Player.PeekDeck(2)

				for _, toMove := range cards {

					card.Player.MoveCard(toMove.ID, match.DECK, match.MANAZONE, card.ID)
					ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the manazone from the top of their deck", card.Player.Username(), toMove.Name))

				}
			}

		}

	})

}

// SlimeVeil ...
func SlimeVeil(c *match.Card) {

	c.Name = "Slime Veil"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				// remove persistent effect on start of next turn
				if _, ok := ctx2.Event.(*match.StartOfTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
				}

				// on all events, add force attack to opponent's creatures
				fx.Find(
					ctx2.Match.Opponent(card.Player),
					match.BATTLEZONE,
				).Map(func(c *match.Card) {

					if _, ok := ctx2.Event.(*match.EndTurnEvent); ok && c.Zone == match.BATTLEZONE {

						if ctx2.Match.IsPlayerTurn(c.Player) && !c.HasCondition(cnd.SummoningSickness) && !c.Tapped {

							if c.HasCondition(cnd.CantAttackPlayers) {

								if c.HasCondition(cnd.CantAttackCreatures) {
									return
								}

								attackableCreatures := fx.FindFilter(
									ctx2.Match.Opponent(c.Player),
									match.BATTLEZONE,
									func(x *match.Card) bool { return x.Tapped || c.HasCondition(cnd.AttackUntapped) })

								if len(attackableCreatures) == 0 {
									return
								}

							}

							ctx2.Match.WarnPlayer(c.Player, fmt.Sprintf("%s must attack before you can end your turn", c.Name))
							ctx2.InterruptFlow()

						}

					}

				})

			})

		}
	})
}

// BrutalCharge ...
func BrutalCharge(c *match.Card) {

	c.Name = "Brutal Charge"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	cardPlayed := false
	shieldsBroken := 0

	c.Use(
		fx.Spell,
		fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
			cardPlayed = true
		}),
		fx.When(fx.ShieldBroken, func(card *match.Card, ctx *match.Context) {

			if cardPlayed {
				shieldsBroken++
			}

		}),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

			if cardPlayed {

				fx.SelectFilterFullList(
					card.Player,
					ctx.Match,
					card.Player,
					match.DECK,
					fmt.Sprintf("Select %d creature from your deck that will be shown to your opponent and sent to your hand", shieldsBroken),
					shieldsBroken,
					shieldsBroken,
					true,
					func(c *match.Card) bool { return c.HasCondition(cnd.Creature) },
					true,
				).Map(func(c *match.Card) {
					c.Player.MoveCard(c.ID, match.DECK, match.HAND, card.ID)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand by %s", c.Name, c.Player.Username(), card.Name))
				})

				card.Player.ShuffleDeck()

				cardPlayed = false
				shieldsBroken = 0

			}

		}))

}

// MiracleQuest ...
func MiracleQuest(c *match.Card) {

	c.Name = "Miracle Quest"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	cardPlayed := false
	shieldsBroken := 0

	c.Use(
		fx.Spell,
		fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
			cardPlayed = true
		}),
		fx.When(fx.ShieldBroken, func(card *match.Card, ctx *match.Context) {

			if cardPlayed {
				shieldsBroken++
			}

		}),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

			if cardPlayed {

				card.Player.DrawCards(shieldsBroken * 2)

				cardPlayed = false
				shieldsBroken = 0

			}
		}))

}

// DivineRiptide ...
func DivineRiptide(c *match.Card) {

	c.Name = "Divine Riptide"
	c.Civ = civ.Water
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		for _, c := range append(fx.Find(card.Player, match.MANAZONE), fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)...) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their mana zone by %s", c.Name, c.Player.Username(), card.Name))
		}

	}))
}

// CataclysmicEruption ...
func CataclysmicEruption(c *match.Card) {

	c.Name = "Cataclysmic Eruption"
	c.Civ = civ.Fire
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		// No. of nature creatures
		n := len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return card.Civ == civ.Nature }))

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("%s: Select upto %d cards from your opponent's manazone and put it in their graveyard", card.Name, n),
			0,
			n,
			false,
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was put into %s's graveyard from their manazone by %s", c.Name, c.Player.Username(), card.Name))
		})

	}))

}

// ThunderNet ...
func ThunderNet(c *match.Card) {

	c.Name = "Thunder Net"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		// No. of water creatures
		n := len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return card.Civ == civ.Water }))

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: select %d of opponent's creatures and tap them", card.Name, n),
			0,
			n,
			false,
		).Map(func(c *match.Card) {
			c.Tapped = true
			ctx.Match.Chat("Server", fmt.Sprintf("%s's %s was tapped by %s", c.Player.Username(), c.Name, card.Name))
		})

	}))

}
