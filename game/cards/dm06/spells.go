package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func ProtectiveForce(c *match.Card) {

	c.Name = "Protective Force"
	c.Civ = civ.Light
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Select 1 creature with \"Blocker\" from your battlezone that will gain +4000 power",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given +4000 power", x.Name))
		})

	}))

}

func InvincibleAura(c *match.Card) {

	c.Name = "Invincible Aura"
	c.Civ = civ.Light
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		topCard := card.Player.PeekDeck(3)

		for _, c := range topCard {
			card.Player.MoveCard(c.ID, match.DECK, match.SHIELDZONE, card.ID)
		}
	}))
}

func InvincibleTechnology(c *match.Card) {

	c.Name = "Invincible Technology"
	c.Civ = civ.Water
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		selectedCards := fx.Select(card.Player, ctx.Match, card.Player, match.DECK, "Select any number of cards from your deck that will be sent to your hand", 0, 100, false)

		for _, selectedCard := range selectedCards {

			card.Player.MoveCard(selectedCard.ID, match.DECK, match.HAND, card.ID)
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

					card.Player.MoveCard(toMove.ID, match.DECK, match.SHIELDZONE, card.ID)
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
				card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
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
			x.Player.MoveCard(x.ID, match.DECK, match.GRAVEYARD, card.ID)
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

			ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
			defer ctx.Match.EndWait(card.Player)

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

		fx.SelectFilterFullList(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Mystic Treasure Chest: Put a non-nature card from your deck into your manazone",
			1,
			1,
			true,
			func(c *match.Card) bool { return c.Civ != civ.Nature },
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.MANAZONE, card.ID)
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

				if len(x.Attachments()) > 0 {
					baseCard := x.Attachments()[0]
					baseCard.Player.MoveCard(baseCard.ID, match.HIDDENZONE, match.BATTLEZONE, card.ID)
					if tapped && !baseCard.Tapped {
						baseCard.Tapped = true
					}
				}

				x.ClearAttachments()
				ctx.Match.MoveCard(x, match.MANAZONE, card)
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

				card.Player.MoveCard(toMove.ID, match.DECK, match.MANAZONE, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the manazone from the top of their deck", card.Player.Username(), toMove.Name))

			}

		}

	})
}

func BondsOfJustice(c *match.Card) {

	c.Name = "Bonds of Justice"
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

func RainOfArrows(c *match.Card) {

	c.Name = "Rain of Arrows"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		opponentHand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil {
			return
		}

		ids := make([]string, 0)

		for _, c := range opponentHand {
			ids = append(ids, c.ImageID)
		}

		ctx.Match.ShowCards(
			card.Player,
			"Your opponent's Hand:",
			ids,
		)

		for _, c := range opponentHand {
			if c.Civ == civ.Darkness && c.HasCondition(cnd.Spell) {
				c.Player.MoveCard(c.ID, match.HAND, match.GRAVEYARD, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's graveyard by %s", c.Name, c.Player.Username(), card.Name))
			}
		}

	}))

}

// CometMissile ...
func CometMissile(c *match.Card) {

	c.Name = "Comet Missile"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Destroy one of your opponent's creatures that has blocker and has power 6000 or less",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) && ctx.Match.GetPower(x, false) <= 6000 })

		for _, creature := range creatures {

			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

		}

	}))

}

// CrisisBoulder ...
func CrisisBoulder(c *match.Card) {

	c.Name = "Crisis Boulder"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		cards["Your creatures"] = fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)
		cards["Your mana"] = fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)

		fx.SelectMultipart(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose 1 of your creatures or a card in your mana zone it will be sent to your graveyard.", card.Name),
			1,
			1,
			false).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, c.Zone, match.GRAVEYARD, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's graveyard by %s", c.Name, c.Player.Username(), card.Name))
		})

	}))
}

// IntenseEvil ...
func IntenseEvil(c *match.Card) {

	c.Name = "Intense Evil"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		battleZone, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		nrCreature := len(battleZone)

		selected := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Intense Evil: Destroy any number of your creatures.",
			0,
			nrCreature,
			false,
		)

		nrSelected := len(selected)

		selected.Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})

		card.Player.DrawCards(nrSelected)

	}))
}

// ShockHurricane ...
func ShockHurricane(c *match.Card) {

	c.Name = "Shock Hurricane"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		battleZone, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		nrCreature := len(battleZone)

		selected := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Shock Hurricane: Return any number of your creatures from the battlezone to your hand.",
			0,
			nrCreature,
			false,
		)

		nrSelected := len(selected)

		selected.Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's battlezone to his hand.", x.Name, x.Player.Username()))
		})

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Choose %d creature(s) in the battlezone that will be sent to your opponent's hand", card.Name, nrSelected),
			0,
			nrSelected,
			false,
		).Map(func(creature *match.Card) {
			creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
		})

	}))

}
