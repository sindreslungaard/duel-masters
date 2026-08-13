package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ArdentLunatron ...
func ArdentLunatron(c *match.Card) {

	c.Name = "Ardent Lunatron"
	c.Power = 5000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.CantAttackCreatures, fx.BlockIfAbleWhenOppAttacks)

}

// HawkeyeLunatron ...
func HawkeyeLunatron(c *match.Card) {

	c.Name = "Hawkeye Lunatron"
	c.Power = 6000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	// The searched card is not shown to the opponent.
	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, fx.SearchDeckTakeXCardsWithoutShowing(1)))

}

// PinpointLunatron ...
func PinpointLunatron(c *match.Card) {

	c.Name = "Pinpoint Lunatron"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.SilentSkill(func(card *match.Card, ctx *match.Context) {

		opponent := ctx.Match.Opponent(card.Player)

		// One choice across four groups. Every creature in the battle zone is
		// eligible, including the card's controller's own and Pinpoint Lunatron
		// itself, and mana is eligible on both sides.
		cards := map[string][]*match.Card{
			"Your creatures":            fx.Find(card.Player, match.BATTLEZONE),
			"Your opponent's creatures": fx.Find(opponent, match.BATTLEZONE),
			"Your mana zone":            fx.Find(card.Player, match.MANAZONE),
			"Your opponent's mana zone": fx.Find(opponent, match.MANAZONE),
		}

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: Choose a creature in the battle zone or a card in either player's mana zone and return it to its owner's hand.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			from := x.Zone

			moved, err := x.Player.MoveCard(x.ID, from, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was returned to %s's hand by %s's effect", x.Name, x.Player.Username(), card.Name))
		})
	}))

}
