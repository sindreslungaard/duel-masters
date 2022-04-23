package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// FullDefensor ...
func FullDefensor(c *match.Card) {

	c.Name = "Full Defensor"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				// on all events, add blocker to our creatures
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
				})

				// remove persistent effect on start of next turn
				_, ok := ctx2.Event.(*match.StartOfTurnStep)
				if ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
				}

			})

		}
	})
}

// CloneFactory ...
func CloneFactory(c *match.Card) {

	c.Name = "Clone Factory"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Clone Factory: Return up to 2 cards from your mana zone to your hand",
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = false
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their manazone by Clone Factory", x.Name, card.Player.Username()))
		})

	}))
}

// MegaDetonator ...
func MegaDetonator(c *match.Card) {

	c.Name = "Mega Detonator"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		hand, err := card.Player.Container(match.HAND)

		if err != nil {
			return
		}

		handLen := len(hand)

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			"Mega Detonator: Choose which cards to discard.",
			0,
			handLen,
			false,
		)

		nrSelected := len(selectedCards)

		selectedCards.Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's graveyard from their hand by Mega Detonator", x.Name, card.Player.Username()))
		})

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		nrCreatures := len(creatures)

		if nrCreatures < nrSelected {
			nrSelected = nrCreatures
		}

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("Mega Detonator: Choose %d creatures that will get double breaker.", nrSelected),
			nrSelected,
			nrSelected,
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.DoubleBreaker, true, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s got double breaker from Mega Detonator", x.Name))
		})

	}))
}

// SwordOfMalevolentDeath ...
func SwordOfMalevolentDeath(c *match.Card) {

	c.Name = "Sword of Malevolent Death"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		nrDarkCards := len(fx.FindFilter(
			card.Player,
			match.MANAZONE,
			func(x *match.Card) bool { return x.Civ == civ.Darkness },
		))

		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {

			x.AddCondition(cnd.PowerAttacker, nrDarkCards*1000, card.ID)
		})
	}))
}

// HydroHurricane ...
func HydroHurricane(c *match.Card) {

	c.Name = "Hydro Hurricane"
	c.Civ = civ.Water
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		nrLight := len(fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Light },
		))

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("Hydro Hurricane: Choose up to %d cards from your opponent's mana zone that will be returned to his hand.", nrLight),
			0,
			nrLight,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s got moved to %s hand from his mana zone by Hydro Hurricane", x.Name, x.Player.Username()))
		})

		nrDark := len(fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Darkness },
		))

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("Hydro Hurricane: Choose up to %d cards from your opponent's battle zone that will be returned to his hand.", nrDark),
			0,
			nrDark,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s got moved to %s hand from his battle zone by Hydro Hurricane", x.Name, x.Player.Username()))
		})
	}))
}

// MysticInscription ...
func MysticInscription(c *match.Card) {

	c.Name = "Mystic Inscription"
	c.Civ = civ.Nature
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		topCard := card.Player.PeekDeck(1)

		for _, c := range topCard {
			card.Player.MoveCard(c.ID, match.DECK, match.SHIELDZONE)
		}
	}))
}

// SwordOfBenevolentLife ...
func SwordOfBenevolentLife(c *match.Card) {

	c.Name = "Sword of Benevolent Life"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		nrLightCards := len(fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Light },
		))

		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {

			x.AddCondition(cnd.PowerAmplifier, nrLightCards*1000, card.ID)
		})

	}))
}

// ChainsOfSacrifice ...
func ChainsOfSacrifice(c *match.Card) {

	c.Name = "Chains of Sacrifice"
	c.Civ = civ.Darkness
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Chains of Sacrifice: Select one of your creatures and destroy it.",
			1,
			1,
			false,
		).Map(func(x *match.Card) {

			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
			ctx.Match.Chat("Server", fmt.Sprintf("%s's %s has been destroyed.", x.Player.Username(), x.Name))
		})

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Chains of Sacrifice: Select up to 2 of your opponent's creatures and destroy them.",
			0,
			2,
			false,
		).Map(func(x *match.Card) {

			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
			ctx.Match.Chat("Server", fmt.Sprintf("%s's %s has been destroyed.", x.Player.Username(), x.Name))
		})

	}))
}

// Darkpact ...
func Darkpact(c *match.Card) {

	c.Name = "Darkpact"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		manaZone, err := card.Player.Container(match.MANAZONE)

		if err != nil {
			return
		}

		nrMana := len(manaZone)

		selected := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Darkpact: Select any number of cards from the manazone and send them to graveyard.",
			0,
			nrMana,
			false,
		)

		nrSelected := len(selected)

		selected.Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's mana zone to his graveyard.", x.Name, x.Player.Username()))
		})

		card.Player.DrawCards(nrSelected)

	}))
}

// SoulGulp ...
func SoulGulp(c *match.Card) {

	c.Name = "Soul Gulp"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		nrLight := len(fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Light },
		))

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		nrDiscard := nrLight
		nrHand := len(fx.Find(ctx.Match.Opponent(card.Player), match.HAND))

		if nrDiscard > nrHand {
			nrDiscard = nrHand
		}

		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			"SoulGulp: Select %d cards from your hand that will be sent to your graveyard",
			nrDiscard,
			nrDiscard,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's hand to his graveyard.", x.Name, x.Player.Username()))
		})

	}))
}

// WhiskingWhirlwind ...
func WhiskingWhirlwind(c *match.Card) {

	c.Name = "Whisking Whirlwind"
	c.Civ = civ.Light
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			// on all events, add blocker to our creatures
			if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.Tapped = false
				})
			}

			// remove persistent effect on start of next turn
			_, ok := ctx2.Event.(*match.StartOfTurnStep)
			if ok && ctx2.Match.IsPlayerTurn(card.Player) {
				exit()
			}

		})
	}))
}

// ScreamingSunburst ...
func ScreamingSunburst(c *match.Card) {

	c.Name = "Screaming Sunburst"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		getNonLightCreatures(card, ctx).Map(func(x *match.Card) {
			x.Tapped = true
		})
	}))
}

func getNonLightCreatures(card *match.Card, ctx *match.Context) fx.CardCollection {

	creatures := fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool { return x.Civ != civ.Light },
	)

	creatures = append(creatures,

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ != civ.Light },
		)...,
	)

	return creatures
}
