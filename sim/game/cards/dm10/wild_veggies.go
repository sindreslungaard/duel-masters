package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ScowlingTomato ...
func ScowlingTomato(c *match.Card) {
	c.Name = "Scowling Tomato"
	c.Power = 2000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)
}

// ShamanBroccoli ...
func ShamanBroccoli(c *match.Card) {
	c.Name = "Shaman Broccoli"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.ReturnToMana))
}

// KaratePotato ...
func KaratePotato(c *match.Card) {

	c.Name = "Karate Potato"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: You may put up to 2 cards from your hand into your mana zone.", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was put from %s's hand to his mana zone.", x.Name, card.Player.Username()))
			}
		})
	}))

}

// HustleBerry ...
func HustleBerry(c *match.Card) {

	c.Name = "Hustle Berry"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.SilentSkill(fx.Draw1ToMana))

}
