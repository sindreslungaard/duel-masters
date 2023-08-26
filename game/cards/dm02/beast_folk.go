package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BarkwhipTheSmasher ...
func BarkwhipTheSmasher(c *match.Card) {

	c.Name = "Barkwhip, the Smasher"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {

			if card.Zone != match.BATTLEZONE || !card.Tapped || event.Card == card || event.Card.Player != card.Player {
				return
			}

			if event.Card.HasFamily(family.BeastFolk) {
				event.Power += 2000
			}

		}

	})

}

// FighterDualFang ...
func FighterDualFang(c *match.Card) {

	c.Name = "Fighter Dual Fang"
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.DrawToMana, fx.DrawToMana)

}

// SilverAxe ...
func SilverAxe(c *match.Card) {

	c.Name = "Silver Axe"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		cards := card.Player.PeekDeck(1)

		if len(cards) < 1 {
			return
		}

		c, err := card.Player.MoveCard(cards[0].ID, match.DECK, match.MANAZONE)

		if err != nil {
			return
		}

		ctx.Match.Chat("Server", fmt.Sprintf("%s was added to %s's manazone from the top of their deck", c.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))

	}))

}

// SilverFist ...
func SilverFist(c *match.Card) {

	c.Name = "Silver Fist"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
