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

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Choose 1 of your opponent's creature in the battlezone and tap it. If it has 'Dragon' in its race, add the top card of your deck to your shields face down.",
			1,
			1,
			false,
		)

		if len(selectedCards) > 0 {

			selectedCards[0].Tapped = true
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was tapped by %s", selectedCards[0].Name, card.Name))

			if selectedCards[0].SharesAFamily(family.Dragons) {
				fx.TopCardToShield(card, ctx)
			}

		}

	}))
}

// LaserWhip ...
func LaserWhip(c *match.Card) {

	c.Name = "Laser Whip"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Choose 1 of your opponent's creature in the battlezone and tap it.",
			1,
			1,
			false,
		)

		if len(selectedCards) > 0 {

			selectedCards[0].Tapped = true
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was tapped by %s", selectedCards[0].Name, card.Name))

			selectedCards = fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				"You may choose 1 of your creatures in the battlezone. If you do, it can't be blocked this turn.",
				1,
				1,
				true,
			)

			if len(selectedCards) > 0 {
				selectedCards[0].Use(func(card2 *match.Card, ctx *match.Context) {
					ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
						if selectedCards[0].Zone != match.BATTLEZONE {
							selectedCards[0].RemoveConditionBySource(selectedCards[0].ID)
							exit()
							return
						}

						selectedCards[0].AddUniqueSourceCondition(cnd.CantBeBlocked, true, selectedCards[0].ID)

						if _, ok := ctx2.Event.(*match.EndStep); ok {
							selectedCards[0].RemoveConditionBySource(selectedCards[0].ID)
							exit()
							return
						}
					})
				})
			}

		}

	}))
}

// LunarCharger ...
func LunarCharger(c *match.Card) {

	c.Name = "Lunar Charger"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Choose up to 2 of your creatures in the battlezone. At the end of the turn, you may untap them.",
			1,
			2,
			true,
		)

		for _, selCard := range selectedCards {

			selCard.Use(func(card2 *match.Card, ctx *match.Context) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if card2.Zone != match.BATTLEZONE {
						exit()
						return
					}

					if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
						ctx2.ScheduleAfter(func() {
							if card2.Tapped {
								// you may untap this creature
								if fx.BinaryQuestion(card2.Player, ctx.Match, fmt.Sprintf("%s's effect: Do you want to untap %s?", card.Name, card2.Name)) {
									card2.Tapped = false
								}
							}

							exit()
						})
					}
				})
			})

		}

	}))
}
