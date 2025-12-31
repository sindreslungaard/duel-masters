package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AncientHornTheWatcher ...
func AncientHornTheWatcher(c *match.Card) {

	c.Name = "Ancient Horn, the Watcher"
	c.Civ = civ.Nature
	c.Power = 5000
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		myShields, _ := card.Player.Container(match.SHIELDZONE)

		if len(myShields) >= 5 {
			fx.Find(
				card.Player,
				match.MANAZONE,
			).Map(func(x *match.Card) {
				x.Tapped = false
			})

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s untapped all his cards in his mana zone.", card.Name, card.Player.Username()))
			ctx.Match.BroadcastState()
		}
	}))

}
