package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// func ProtectiveForce(c *match.Card) {

// 	c.Name = "Protective Force"
// 	c.Civ = civ.Light
// 	c.ManaCost = 1
// 	c.ManaRequirement = []string{civ.Light}

// 	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

// 		fx.SelectFilter(
// 			card.Player,
// 			ctx.Match,
// 			card.Player,
// 			match.BATTLEZONE,
// 			"Protective Force: Select a blocker to give +4000 power to until the end of the turn",
// 			1,
// 			1,
// 			false,
// 			func(x *match.Card) bool {
// 				return x.HasCondition(cnd.Blocker)
// 			},
// 		).Map(func(x *match.Card) {
// 			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
// 		})

// 		// TODO: make sure the condition is removed at the end of opponent's turn if shield trigger

// 	}))
// }

func InvincibleAura(c *match.Card) {

	c.Name = "Invincible Aura"
	c.Civ = civ.Light
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		topCard := card.Player.PeekDeck(3)

		for _, c := range topCard {
			card.Player.MoveCard(c.ID, match.DECK, match.SHIELDZONE)
		}
	}))
}

func InvincibleTechnology(c *match.Card) {

	c.Name = "Invincible Technology"
	c.Civ = civ.Water
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		selectedCards := match.Search(card.Player, ctx.Match, card.Player, match.DECK, "Select any number of cards from your deck that will be sent to your hand", 1, 100, false)

		for _, selectedCard := range selectedCards {

			card.Player.MoveCard(selectedCard.ID, match.DECK, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the deck to their hand", card.Player.Username(), selectedCard.Name))
		}

		card.Player.ShuffleDeck()

	}))
}

func InvincibleAbyss(c *match.Card) {

	c.Name = "Invincible Abyss"
	c.Civ = civ.Darkness
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range creatures {
			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
		}
	}))
}

func InvincibleCataclysm(c *match.Card) {

	c.Name = "Invincible Cataclysm"
	c.Civ = civ.Fire
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			"Select up to 3 of your opponent's shields and send them to graveyard",
			1,
			3,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})

	}))
}

func InvincibleUnity(c *match.Card) {

	c.Name = "Invincible Unity"
	c.Civ = civ.Nature
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range creatures {

			creature.AddCondition(cnd.PowerAmplifier, 8000, card.ID)
			creature.AddCondition(cnd.TripleBreaker, nil, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given +8000 power and triple breaker until the end of the turn", creature.Name))

		}
	}))
}

func SphereOfWonder(c *match.Card) {

	c.Name = "Sphere of Wonder"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			if len(fx.Find(card.Player, match.SHIELDZONE)) < len(fx.Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)) {

				cards := card.Player.PeekDeck(1)

				for _, toMove := range cards {

					card.Player.MoveCard(toMove.ID, match.DECK, match.SHIELDZONE)
					ctx.Match.Chat("Server", fmt.Sprintf("%s put a shield into the shieldzone from the top of their deck by Sphere of Wonder's ability", card.Player.Username()))

				}
			}

		}

	})

}

func MysticDreamscape(c *match.Card) {

	c.Name = "Mystic Dreamscape"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.MANAZONE,
				"Mystic Dreamscape: Return up to 3 cards from your mana zone to your hand",
				1,
				3,
				true,
			).Map(func(x *match.Card) {
				x.Tapped = false
				card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their manazone by Mystic Dreamscape", x.Name, card.Player.Username()))
			})

		}
	})

}

func FutureSlash(c *match.Card) {

	c.Name = "Future Slash"
	c.Civ = civ.Darkness
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		opponent := ctx.Match.Opponent(card.Player)

		fx.Select(
			card.Player,
			ctx.Match,
			opponent,
			match.DECK,
			"Future Slash: Send up to 2 cards from your opponent's deck to his graveyard",
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s put %s in graveyard from their opponent's deck", x.Player.Username(), x.Name))
		})

		opponent.ShuffleDeck()

	}))

}

func ProclamationOfDeath(c *match.Card) {

	c.Name = "Proclamation of Death"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {
		if match.AmICasted(card, ctx) {
			creatures := fx.Select(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Proclamation of Death: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

			for _, creature := range creatures {
				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
			}
		}
	})
}

func PhantomDragonsFlame(c *match.Card) {

	c.Name = "Phantom Dragon's Flame"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures that will be destroyed", 1, 1, false, func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 })

			for _, creature := range creatures {

				ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

			}

		}

	})
}

func SpasticMissile(c *match.Card) {

	c.Name = "Spastic Missile"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures that will be destroyed", 1, 1, false, func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 3000 })

			for _, creature := range creatures {

				ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

			}

		}

	})
}

func MysticTreasureChest(c *match.Card) {

	c.Name = "Mystic Treasure Chest"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Mystic Treasure Chest: Put a non-nature card from your deck into your manazone",
			1,
			1,
			true,
			func(c *match.Card) bool { return c.Civ != civ.Nature }).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.MANAZONE)
			ctx.Match.Chat("Server", fmt.Sprintf("%s put %s in their manazone from their deck", x.Player.Username(), x.Name))
		})

		card.Player.ShuffleDeck()

	}))
}

func PangaeasWill(c *match.Card) {

	c.Name = "Pangaea's Will"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Pangaea's Will: Select an evolution card from your opponent's battle zone and send the top card to their mana zone",
				0,
				1,
				true,
				func(x *match.Card) bool { return x.HasCondition(cnd.Evolution) },
			).Map(func(x *match.Card) {
				tapped := x.Tapped
				baseCard := x.Attachments()[0]
				x.ClearAttachments()
				ctx.Match.MoveCard(x, match.MANAZONE, card)
				baseCard.Player.MoveCard(baseCard.ID, match.HIDDENZONE, match.BATTLEZONE)
				if tapped && !baseCard.Tapped {
					baseCard.Tapped = true
				}
			})
		}

	})
}

func FaerieLife(c *match.Card) {

	c.Name = "Faerie Life"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			cards := card.Player.PeekDeck(1)

			for _, toMove := range cards {

				card.Player.MoveCard(toMove.ID, match.DECK, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the manazone from the top of their deck", card.Player.Username(), toMove.Name))

			}

		}

	})
}

func BondsOfJustice(c *match.Card) {

	c.Name = "Bonds of Jusitce"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			notblockers := make([]*match.Card, 0)

			opponentBattlezone, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			myBattlezone, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			for _, creature := range myBattlezone {
				if !creature.HasCondition(cnd.Blocker) {
					notblockers = append(notblockers, creature)
				}
			}

			for _, creature := range opponentBattlezone {
				if !creature.HasCondition(cnd.Blocker) {
					notblockers = append(notblockers, creature)
				}
			}
			for _, notblocker := range notblockers {
				notblocker.Tapped = true
				ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by Bonds of Justice", notblocker.Name))
			}
		}

	})
}

func EnergyStream(c *match.Card) {

	c.Name = "Energy Stream"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			card.Player.DrawCards(2)

		}
	})

}
