package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func PyrofighterMagnus(c *match.Card) {

	c.Name = "Pyrofighter Magnus"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
		ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to the %s's hand", c.Name, c.Player.Username()))
	}))
}

func Torchclencher(c *match.Card) {

	c.Name = "Torchclencher"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking && match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x != c && x.Civ == civ.Fire }) {
			power += 3000
		}

		return power
	}

}
