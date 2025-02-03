package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DimensionSplitter ...
func DimensionSplitter(c *match.Card) {
	c.Name = "Dimension Splitter"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature,
		fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
			fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(creature *match.Card) bool { return creature.SharesAFamily(family.Dragons) },
			).Map(func(creature *match.Card) {
				creature.Player.MoveCard(creature.ID, match.GRAVEYARD, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their graveyard by %s", creature.Name, card.Player.Username(), card.Name))
			})
		}),
	)
}
