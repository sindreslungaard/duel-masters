package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Psyshroom ...
func Psyshroom(c *match.Card) {

	c.Name = "Psyshroom"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BalloonMushroom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: You may choose a nature card from your graveyard to put into your mana zone", card.Name),
			0,
			1,
			true,
			func(x *match.Card) bool { return x.Civ == civ.Nature },
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to %s's manazone from their graveyard by %s's effect.", x.Name, card.Player.Username(), card.Name))
		})
	}))

}
