package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigaslug ...
func Gigaslug(c *match.Card) {

	c.Name = "Gigaslug"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker(), fx.Slayer, fx.CantAttackCreatures, fx.CantAttackPlayers)

}

// Gigabalza ...
func Gigabalza(c *match.Card) {

	c.Name = "Gigabalza"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, fx.OpponentDiscardsRandomCard))

}

// GigappiPonto ...
func GigappiPonto(c *match.Card) {

	c.Name = "Gigappi Ponto"
	c.Power = 4000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Chimera, family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}

// Gigarayze ...
func Gigarayze(c *match.Card) {

	c.Name = "Gigarayze"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: You may return a water or fire creature from your graveyard to your hand.", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool {
				return x.HasCondition(cnd.Creature) && (x.HasCiv(civ.Water) || x.HasCiv(civ.Fire))
			},
			true,
		).Map(func(chosen *match.Card) {
			moved, err := card.Player.MoveCard(chosen.ID, match.GRAVEYARD, match.HAND, card.ID)

			if err == nil && moved.Zone == match.HAND {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand by %s", chosen.Name, card.Player.Username(), card.Name))
			}
		})
	}))

}

// gigavrandDrawThreshold is how many cards the opponent has to have drawn during
// a turn for Gigavrand to empty their hand at the end of it.
const gigavrandDrawThreshold = 2

// Gigavrand ...
func Gigavrand(c *match.Card) {

	c.Name = "Gigavrand"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	// Nothing in the engine remembers how much was drawn, so Gigavrand counts
	// for itself. The tally runs whatever zone this card is in, because the
	// printed condition asks about the whole turn rather than about the part of
	// it Gigavrand was in play for; only the punishment needs the battle zone.
	drawnThisTurn := 0

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		switch ctx.Event.(type) {

		case *match.CardMoved:
			event := ctx.Event.(*match.CardMoved)

			// Player.DrawCards is the only thing that tags a move this way, so
			// searching a card out of the deck is not drawing it.
			if event.Source != "draw" || event.From != match.DECK || event.To != match.HAND {
				return
			}

			drawer := ctx.Match.Player1.Player
			if event.MatchPlayerID != 1 {
				drawer = ctx.Match.Player2.Player
			}

			if drawer == ctx.Match.Opponent(card.Player) {
				drawnThisTurn++
			}

		case *match.EndStep:
			if card.Zone != match.BATTLEZONE || drawnThisTurn < gigavrandDrawThreshold {
				return
			}

			fx.OpponentDiscardsHand(card, ctx)

		case *match.EndOfTurnStep:
			// Reset after the end step has had its look, so the next turn starts
			// counting from nothing.
			drawnThisTurn = 0
		}
	})

}
