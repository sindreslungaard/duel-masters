package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// RapidReincarnation ...
func RapidReincarnation(c *match.Card) {

	c.Name = "Rapid Reincarnation"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: You may destroy one of your creatures.\r\nIf you do, choose a creature in your hand that costs the same as or less than the number of cards in your mana zone\r\nand put it into the battle zone.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			manaZone, _ := card.Player.Container(match.MANAZONE)

			if manaZone != nil {
				ctx.Match.Destroy(x, card, match.DestroyedBySpell)

				manaCount := len(manaZone)

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					card.Player,
					match.HAND,
					fmt.Sprintf("%s's effect: Choose a creature in your hand that costs the same as or less than the number of cards in your mana zone\r\nand put it into the battle zone.", card.Name),
					1,
					1,
					false,
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature) && x.ManaCost <= manaCount && fx.CanBeSummoned(card.Player, x)
					},
					false,
				).Map(func(x *match.Card) {
					fx.ForcePutCreatureIntoBZ(ctx, x, match.HAND, card)
				})
			}
		})
	}))

}

// StaticWarp ...
func StaticWarp(c *match.Card) {

	c.Name = "Static Warp"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	textFmt := "%s's effect: Choose a creature in your battle zone. Tap the rest of the creatures in the battle zone."

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		players := []*match.Player{card.Player, ctx.Match.Opponent(card.Player)}

		for _, p := range players {
			fx.Select(
				p,
				ctx.Match,
				p,
				match.BATTLEZONE,
				fmt.Sprintf(textFmt, card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				fx.FindFilter(
					p,
					match.BATTLEZONE,
					func(y *match.Card) bool {
						return y.ID != x.ID
					},
				).Map(func(y *match.Card) {
					y.Tapped = true
				})
			})
		}

		ctx.Match.BroadcastState()
	}))

}

// SirenConcerto ...
func SirenConcerto(c *match.Card) {

	c.Name = "Siren Concerto"
	c.Civ = civ.Water
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Put a card from your mana zone into your hand.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			if err == nil {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into his hand from his mana zone.", card.Player.Username(), x.Name))
			}
		})

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: Put a card from your hand into your mana zone.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := card.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)
			if err == nil {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into his mana zone from his hand.", card.Player.Username(), x.Name))
			}
		})
	}))

}

// Transmogrify ...
func Transmogrify(c *match.Card) {

	c.Name = "Transmogrify"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		cards := make(map[string][]*match.Card)

		myCards, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		opponentCards, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: You may destroy a creature. If you do, its owner reveals cards from the top of this deck until he reveals a non-evolution creature. He puts that creature into the battle zone and puts the rest of those cards into his graveyard.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)

			for {
				topCards := x.Player.PeekDeck(1)

				if len(topCards) < 1 {
					return
				}

				topCard := topCards[0]

				if topCard.HasCondition(cnd.Creature) && !topCard.HasCondition(cnd.Evolution) {
					fx.ForcePutCreatureIntoBZ(ctx, topCard, match.DECK, card)
					return
				}

				ctx.Match.MoveCard(topCard, match.GRAVEYARD, card)
			}
		})
	}))

}

// InfernalCommand ...
func InfernalCommand(c *match.Card) {

	c.Name = "Infernal Command"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	//@TODO test this with the Nariel fix with Machai and ForceAttack!!
	// What happens if the opponent casts InfernalCommand or Slime Veil?
	// Or any other spell that gives ForceAttack to other creatures that are blocked from attacking by Nariel??

	// Perhaps modify Nariel effect to NOT remove CantAttack conditions on EndOfTurn event
	// so that creatures that are forced to attack by other effects will still be forced to attack
	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// remove persistent effect on start of next turn
			if _, ok := ctx2.Event.(*match.StartOfTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
				exit()
				return
			}

			// on all events, add force attack to opponent's creatures
			fx.Find(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(c *match.Card) {
				if _, ok := ctx2.Event.(*match.EndTurnEvent); ok && c.Zone == match.BATTLEZONE {
					if ctx2.Match.IsPlayerTurn(c.Player) && !fx.HasSummoningSickness(c) && !c.Tapped {
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
	}))

}

// Upheaval ...
func Upheaval(c *match.Card) {

	c.Name = "Upheaval"
	c.Civ = civ.Darkness
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		myManaCards := fx.Find(card.Player, match.MANAZONE)
		myHandCards := fx.FindFilter(card.Player, match.HAND, func(x *match.Card) bool { return x.ID != card.ID })
		myOppManaCards := fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)
		myOppHandCards := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)

		for _, x := range myManaCards {
			x.Player.MoveCard(x.ID, x.Zone, match.HAND, card.ID)
		}

		for _, x := range myHandCards {
			x.Player.MoveCard(x.ID, x.Zone, match.MANAZONE, card.ID)
			x.Tapped = true
		}

		for _, x := range myOppManaCards {
			x.Player.MoveCard(x.ID, x.Zone, match.HAND, card.ID)
		}

		for _, x := range myOppHandCards {
			x.Player.MoveCard(x.ID, x.Zone, match.MANAZONE, card.ID)
			x.Tapped = true
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: Both players moved their mana cards to their hand, and at the same time, their hand cards to their mana zone.", card.Name))
	}))

}

// ColossusBoost ...
func ColossusBoost(c *match.Card) {

	c.Name = "Colossus Boost"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: One of your creatures in the battle zone gets '+4000 Power' until the end of the turn.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				x.AddUniqueSourceCondition(cnd.PowerAmplifier, 4000, card.ID)
			})
		})
	}))

}

// ForcedFrenzy ...
func ForcedFrenzy(c *match.Card) {

	c.Name = "Forced Frenzy"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Find(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.StartOfTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
					return
				}

				fx.ForceAttack(x, ctx2)
			})
		})
	}))

}

// SupersonicJetpack ...
func SupersonicJetpack(c *match.Card) {

	c.Name = "Supersonic Jetpack"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: One of your creatures in the battlezone gets 'speed attacker' until end of turn.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
		})
	}))

}
