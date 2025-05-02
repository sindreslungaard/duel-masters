package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// FortressShell ...
func FortressShell(c *match.Card) {

	c.Name = "Fortress Shell"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("%s: Select up to 2 cards in your opponent's mana zone that will be sent to their graveyard", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", x.Name, x.Player.Username(), card.Name))
		})
	}))

}
