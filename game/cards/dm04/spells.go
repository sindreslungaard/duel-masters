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

			ctx.Match.ApplyPersistentEffect(func(_ *match.Card, ctx2 *match.Context, exit func()) {

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
			
			x.PowerModifier = func(m *match.Match, attacking bool) int {

				if attacking {
					return nrDarkCards * 1000
				}

				return 0
			}
		})
	}))
}
