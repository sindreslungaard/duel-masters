package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CorpseCharger ...
func CorpseCharger(c *match.Card) {
	c.Name = "Corpse Charger"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, fx.ReturnXCreaturesFromGraveToHand(1)))
}

// CraniumClamp ...
func CraniumClamp(c *match.Card) {
	c.Name = "Cranium Clamp"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.OpDiscardsXCards(2)))
}

// VolcanoCharger ...
func VolcanoCharger(c *match.Card) {

	c.Name = "Volcano Charger"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, fx.DestroyBySpellOpCreature2000OrLess))
}

// EurekaCharger ...
func EurekaCharger(c *match.Card) {

	c.Name = "Eureka Charger"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, fx.Draw1))
}

// MuscleCharger ...
func MuscleCharger(c *match.Card) {

	c.Name = "Muscle Charger"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Find(card.Player, match.BATTLEZONE).
			Map(func(creature *match.Card) {
				creature.AddCondition(cnd.PowerAmplifier, 3000, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given +3000 power until the end of the turn", creature.Name))
			})

	}))
}

// Dracobarrier ...
func Dracobarrier(c *match.Card) {

	c.Name = "Dracobarrier"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose 1 of your opponent's creature in the battlezone and tap it. If it has 'Dragon' in its race, add the top card of your deck to your shields face down.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was tapped by %s", x.Name, card.Name))

			if x.SharesAFamily(family.Dragons) {
				fx.TopCardToShield(card, ctx)
			}
		})

	}))
}

// LaserWhip ...
func LaserWhip(c *match.Card) {

	c.Name = "Laser Whip"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose 1 of your opponent's creature in the battlezone and tap it.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.BroadcastState()
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was tapped by %s", x.Name, card.Name))

			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may choose 1 of your creatures in the battlezone. If you do, it can't be blocked this turn.", card.Name),
				1,
				1,
				true,
			).Map(func(y *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if y.Zone != match.BATTLEZONE {
						y.RemoveConditionBySource(y.ID)
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						y.RemoveConditionBySource(y.ID)
						exit()
						return
					}

					y.AddUniqueSourceCondition(cnd.CantBeBlocked, true, y.ID)
				})
			})
		})
	}))

}

// LunarCharger ...
func LunarCharger(c *match.Card) {

	c.Name = "Lunar Charger"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose up to 2 of your creatures in the battlezone. At the end of the turn, you may untap them.", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					ctx2.ScheduleAfter(func() {
						if x.Tapped {
							// you may untap this creature
							if fx.BinaryQuestion(x.Player, ctx2.Match, fmt.Sprintf("%s's effect: Do you want to untap %s?", card.Name, x.Name)) {
								x.Tapped = false
								ctx2.Match.BroadcastState()
							}
						}

						exit()
					})
				}
			})
		})
	}))

}

// RootCharger ...
func RootCharger(c *match.Card) {

	c.Name = "Root Charger"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if event, ok := ctx2.Event.(*match.CreatureDestroyed); ok {
				if event.Card.Player == card.Player {
					card.Player.MoveCard(event.Card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
					ctx2.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s was moved to your manazone instead of being destroyed.", card.Name, event.Card.Name))
				}
			}

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				ctx2.ScheduleAfter(func() {
					exit()
				})
			}
		})
	}))

}

// MarineScramble ...
func MarineScramble(c *match.Card) {

	c.Name = "Marine Scramble"
	c.Civ = civ.Water
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			creatures := fx.Find(card.Player, match.BATTLEZONE)

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				creatures.Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})
				exit()
				return
			}

			creatures.Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)
			})
		})
	}))

}

// WaveLance ...
func WaveLance(c *match.Card) {

	c.Name = "Wave Lance"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards := fx.Find(
			card.Player,
			match.BATTLEZONE,
		)

		opponentCards := fx.Find(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
		)

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose a creature in the battlezone and return it to its owner's hand. If it has 'Dragon' in its race, you may draw a card.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.HAND, card)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to its owner's hand by %s", x.Name, card.Name))

			if x.SharesAFamily(family.Dragons) {
				fx.MayDraw1(card, ctx)
			}
		})

	}))

}

// FuriousOnslaught ...
func FuriousOnslaught(c *match.Card) {

	c.Name = "Furious Onslaught"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(creature *match.Card) bool { return creature.HasFamily(family.Dragonoid) },
		).Map(func(creature *match.Card) {
			creature.AddCondition(cnd.AddFamily, family.ArmoredDragon, card)
			creature.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
			creature.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given +4000 power, double breaker, and is now an Armored Dragon until the end of the turn", creature.Name))
		})
	}))

}

// DimensionSplitter ...
func DimensionSplitter(c *match.Card) {

	c.Name = "Dimension Splitter"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature,
		fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
			fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(x *match.Card) bool { return x.SharesAFamily(family.Dragons) },
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their graveyard by %s", x.Name, card.Player.Username(), card.Name))
			})
		}))

}
