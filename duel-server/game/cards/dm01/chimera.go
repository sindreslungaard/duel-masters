package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigaberos ...
func Gigaberos(c *match.Card) {

	c.Name = "Gigaberos"
	c.Power = 8000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		otherCreatures := make([]*match.Card, 0)
		for _, creature := range creatures {
			if creature.ID != card.ID {
				otherCreatures = append(otherCreatures, creature)
			}
		}

		if len(otherCreatures) < 2 {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("You don't have 2 other creatures to destroy. Your %s will be destroyed instead.", card.Name))
			ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
			return
		}

		choiceIndex := fx.MultipleChoiceQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to destroy 2 of your other creatures, or destroy %s instead?", card.Name), []string{"Destroy 2 of my other creatures", fmt.Sprintf("Destroy %s", card.Name)})

		if choiceIndex == 0 {
			// Destroy 2 other creatures
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				"Select 2 of your other creatures to destroy",
				2,
				2,
				false,
				func(x *match.Card) bool {
					return x.ID != card.ID
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		} else {
			// Destroy this creature
			ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
		}
	}))

}

// Gigagiele ...
func Gigagiele(c *match.Card) {

	c.Name = "Gigagiele"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}

// Gigargon ...
func Gigargon(c *match.Card) {

	c.Name = "Gigargon"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: Select up to 2 creatures from your graveyard that will be returned to your hand", card.Name),
			1,
			2,
			true,
			func(x *match.Card) bool {
				return x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was returned to %s's hand from their graveyard", x.Name, x.Player.Username()))
		})
	}))

}
