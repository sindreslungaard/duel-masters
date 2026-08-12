package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DiamondiaTheBlizzardRider ...
func DiamondiaTheBlizzardRider(c *match.Card) {

	c.Name = "Diamondia, the Blizzard Rider"
	c.Power = 5000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		isSnowFaerie := func(x *match.Card) bool { return x.HasFamily(family.SnowFaerie) }

		// Both zones are snapshotted before anything moves, so a card leaving
		// one of them cannot disturb the walk over the other.
		gathered := make([]*match.Card, 0)
		for _, zone := range []string{match.GRAVEYARD, match.MANAZONE} {
			fx.FindFilter(card.Player, zone, isSnowFaerie).Map(func(x *match.Card) {
				gathered = append(gathered, x)
			})
		}

		returned := 0
		for _, faerie := range gathered {
			moved, err := card.Player.MoveCard(faerie.ID, faerie.Zone, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				continue
			}

			returned++
		}

		if returned > 0 {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s returned %d Snow Faerie(s) to %s's hand", card.Name, returned, card.Player.Username()))
		}
	}))

}
