package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func NeonCluster(c *match.Card) {

	c.Name = "Neon Cluster"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability to draw 2 cards", card.Player.Username(), card.Name))
		card.Player.DrawCards(2)
	}

	c.Use(fx.Creature, fx.TapAbility)
}

func OverloadCluster(c *match.Card) {

	c.Name = "Overload Cluster"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	isBlocker := false

	c.Use(fx.Creature,
		fx.When(ReceiveBlockerWhenOpponentPlaysCreatureOrSpell,
			func(c *match.Card, ctx *match.Context) { isBlocker = true }),
		fx.When(fx.EndOfTurn,
			func(c *match.Card, ctx *match.Context) {
				isBlocker = false
				c.RemoveConditionBySource(c.ID)
			}),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return isBlocker },
			func(c *match.Card, ctx *match.Context) {
				fx.ForceBlocker(c, ctx, c.ID)
			}),
	)
}

func ReceiveBlockerWhenOpponentPlaysCreatureOrSpell(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.SpellCast); ok {
		return shouldReceiveBlockerWhenOpponentPlaysCard(card, ctx, event.MatchPlayerID)
	}

	if fx.AnotherCreatureSummoned(card, ctx) {
		// TODO: This check can be removed once the card in CardMoved is passed as pointer
		// And MatchPlayerID is removed
		if event, ok := ctx.Event.(*match.CardMoved); ok {
			return shouldReceiveBlockerWhenOpponentPlaysCard(card, ctx, event.MatchPlayerID)
		}

	}

	return false

}

func shouldReceiveBlockerWhenOpponentPlaysCard(card *match.Card, ctx *match.Context, playedCardPlayerId byte) bool {

	if card.Zone != match.BATTLEZONE || playedCardPlayerId == 0 {
		return false
	}

	// Return if it's not the opponent that plays the card
	var playedCardPlayer *match.Player
	if playedCardPlayerId == 1 {
		playedCardPlayer = ctx.Match.Player1.Player
	} else {
		playedCardPlayer = ctx.Match.Player2.Player
	}

	return card.Player != playedCardPlayer

}

func FortMegacluster(c *match.Card) {

	c.Name = "Fort Megacluster"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = fortMegaclusterTapAbility

	c.Use(fx.Creature, fx.Evolution, fx.TapAbility,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			fx.GiveTapAbilityToAllies(
				card,
				ctx,
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Water },
				fortMegaclusterTapAbility,
			)

		}),
	)
}

func fortMegaclusterTapAbility(card *match.Card, ctx *match.Context) {
	card.Player.DrawCards(1)
	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability to draw 1 card", card.Player.Username(), card.Name))
}
