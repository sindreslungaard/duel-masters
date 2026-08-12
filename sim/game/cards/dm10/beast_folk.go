package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AdventureBoar ...
func AdventureBoar(c *match.Card) {

	c.Name = "Adventure Boar"
	c.Civs = []string{civ.Nature}
	c.Power = 1000
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// EarthRipperTalonOfRage ...
func EarthRipperTalonOfRage(c *match.Card) {

	c.Name = "Earth Ripper, Talon of Rage"
	c.Power = 6000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.MANAZONE,
			func(x *match.Card) bool {
				return x.Tapped
			},
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved from %s's mana zone to his hand due to %s's effect.", x.Name, card.Player.Username(), card.Name))
			}
		})
	}))

}

// SkyswordTheSavageVizier ...
func SkyswordTheSavageVizier(c *match.Card) {

	c.Name = "Skysword, the Savage Vizier"
	c.Power = 2000
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.BeastFolk, family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	// Two separate top cards: the first goes to mana, the next to the shields.
	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Draw1ToMana(card, ctx)
		fx.TopCardToShield(card, ctx)
	}))

}

// SanfistTheSavageVizier ...
func SanfistTheSavageVizier(c *match.Card) {

	c.Name = "Sanfist, the Savage Vizier"
	c.Power = 3000
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.BeastFolk, family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Blocker(), func(card *match.Card, ctx *match.Context) {

		// "When this creature would be discarded from your hand during your
		// opponent's turn, you may put it into the battle zone instead."
		event, ok := ctx.Event.(*match.MoveCard)

		if !ok ||
			event.CardID != card.ID ||
			event.From != match.HAND ||
			event.To != match.GRAVEYARD ||
			ctx.Match.IsPlayerTurn(card.Player) {
			return
		}

		if !fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to put %s into the battle zone instead of discarding it?", card.Name)) {
			return
		}

		// Cancelling the pending move is what replaces the discard; the card is
		// then put into play from the hand it never left.
		ctx.InterruptFlow()

		fx.ForcePutCreatureIntoBZ(ctx, card, match.HAND, card)
	})

}
