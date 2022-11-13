package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SyriusFirmamentElemental ...
func SyriusFirmamentElemental(c *match.Card) {

	c.Name = "Syrius, Firmament Elemental"
	c.Power = 12000
	c.Civ = civ.Light
	c.Family = family.AngelCommand
	c.ManaCost = 11
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.Triplebreaker)

}

// SyforceAuroraElemental ...
func SyforceAuroraElemental(c *match.Card) {

	c.Name = "Syforce, Aurora Elemental"
	c.Power = 7000
	c.Civ = civ.Light
	c.Family = family.AngelCommand
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		spells := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "You may select 1 spell from your mana zone that will be sent to your hand", 0, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

		for _, spell := range spells {

			card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", spell.Player.Username(), spell.Name))

		}
	}))

}
