package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TwitchHornTheAggressor ...
func TwitchHornTheAggressor(c *match.Card) {
	c.Name = "Twitch Horn, the Aggressor"
	c.Power = 2000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if !attacking || c.Zone != match.BATTLEZONE {
			return 0
		}

		return len(fx.FindFilter(
			c.Player,
			match.MANAZONE,
			func(card *match.Card) bool { return card.Tapped },
		)) * 2000
	}

	c.Use(fx.Creature)
}

// AncientHornTheWatcher ...
func AncientHornTheWatcher(c *match.Card) {

	c.Name = "Ancient Horn, the Watcher"
	c.Civs = []string{civ.Nature}
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
