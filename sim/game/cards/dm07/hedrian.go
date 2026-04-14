package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

func BattleshipMutant(c *match.Card) {

	c.Name = "Battleship Mutant"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(card.Player, match.BATTLEZONE,
			func(c *match.Card) bool { return c.Civ == civ.Darkness },
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
			x.AddCondition(cnd.DoubleBreaker, true, card.ID)
			x.AddCondition(cnd.DestroyAfterBattle, true, card.ID)
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}

func PropellerMutant(c *match.Card) {

	c.Name = "Propeller Mutant"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		// Reimplementing here OpponentDiscardsRandomCard with a twist. Since the discard is activated when this creature is destroyed,
		// if it is destoryed by a spell, witouht this, there is a chance you would discard that spell from the hand since we only
		// move the card from hand after its effect was applied.
		// The proper solution would be to use a spellzone for when spells are active. This is a work around.

		event, ok := ctx.Event.(*match.CardMoved)
		if !ok {
			return
		}

		hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)
		if err != nil || len(hand) < 1 {
			return
		}

		randomcard := hand[rand.Intn(len(hand))]
		for randomcard.ID == event.Source {
			if len(hand) == 1 {
				return
			}
			randomcard = hand[rand.Intn(len(hand))]
		}

		discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(randomcard.ID, match.HAND, match.GRAVEYARD, card.ID)
		if err == nil {
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
		}
	}))

}
