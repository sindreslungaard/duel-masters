package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// UnifiedResistance ...
func UnifiedResistance(c *match.Card) {

	c.Name = "Unified Resistance"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(
			card,
			ctx,
			fmt.Sprintf("%s's effect: Choose a race. Until the start of your next turn, each of your creatures in the battlezone of that race gets 'Blocker'", card.Name),
		)

		if family != "" {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("Until the start of your next turn, each of your '%s' creatures in the battlezone gets 'Blocker'", family))

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family)
					},
				).Map(func(x *match.Card) {
					_, ok := ctx2.Event.(*match.StartOfTurnStep)
					if ok && ctx2.Match.IsPlayerTurn(card.Player) {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					fx.ForceBlocker(x, ctx2, card.ID)
				})
			})
		}
	}))

}

// ImpossibleTunnel ...
func ImpossibleTunnel(c *match.Card) {

	c.Name = "Impossible Tunnel"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(
			card,
			ctx,
			fmt.Sprintf("%s's effect: Choose a race. Creatures of that race can't be blocked this turn.", card.Name),
		)

		if family != "" {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("'%s' creatures can't be blocked this turn.", family))

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family)
					},
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)

					_, ok := ctx2.Event.(*match.EndOfTurnStep)
					if ok && ctx2.Match.IsPlayerTurn(card.Player) {
						x.RemoveConditionBySource(card.ID)
						exit()
					}
				})
			})
		}
	}))

}

// ZombieCarnival ...
func ZombieCarnival(c *match.Card) {

	c.Name = "Zombie Carnival"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Return up to 3 creatures of that race from your graveyard to your hand.", card.Name))

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: Return up to 3 creatures of that race from your graveyard to your hand.", card.Name),
			1,
			3,
			true,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was return to your hand from your graveyard by %s's effect.", x.Name, card.Name))
		})
	}))

}

// DanceOfTheSproutlings ...
func DanceOfTheSproutlings(c *match.Card) {

	c.Name = "Dance of the Sproutlings"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. You may put any number of creatures of that race from your hand into your manazone.", card.Name))

		max := len(fx.FindFilter(
			card.Player,
			match.HAND,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
		))

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: You may put any number of '%s' creatures from your hand into your manazone.", card.Name, family),
			1,
			max,
			true,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's manazone from his hand due to %s's effect.", x.Name, card.Player.Username(), card.Name))
		})
	}))

}

// RelentlessBlitz ...
func RelentlessBlitz(c *match.Card) {

	c.Name = "Relentless Blitz"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. This turn, each creature of that race can attack untapped creatures and can't be blocked while attacking a creature.", card.Name))

		if family != "" {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("This turn, '%s' creatures can attack untapped creatures and can't be blocked while attacking a creature.", family))

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family) && x.HasCondition(cnd.Creature)
					},
				).Map(func(x *match.Card) {
					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					x.AddUniqueSourceCondition(cnd.AttackUntapped, true, card.ID)
					fx.CantBeBlockedWhileAttackingACreature(x, ctx2)
				})
			})
		}
	}))

}
