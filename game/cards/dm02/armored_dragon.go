package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BolzardDragon ...
func BolzardDragon(c *match.Card) {

	c.Name = "Bolzard Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 card from your opponent's mana zone that will be sent to their graveyard", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", x.Name, x.Player.Username(), card.Name))
		})
	}))

}
