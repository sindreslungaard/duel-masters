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
