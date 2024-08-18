package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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

func LavaWalkerExecuto(c *match.Card) {

	c.Name = "Lava Walker Executo"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = lavaWalkerExecutoTapAbility

	c.Use(fx.Creature, fx.Evolution, fx.TapAbility,
		fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

			fx.GiveTapAbilityToAllies(
				card,
				ctx,
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Fire },
				lavaWalkerExecutoTapAbility,
			)

		}),
	)

}

func lavaWalkerExecutoTapAbility(card *match.Card, ctx *match.Context) {
	creatures := fx.SelectFilter(
		card.Player,
		ctx.Match,
		card.Player,
		match.BATTLEZONE,
		"Select 1 fire creature from your battlezone that will gain +3000 Power",
		1,
		1,
		false,
		func(x *match.Card) bool { return x.Civ == civ.Fire },
	)
	for _, creature := range creatures {
		creature.AddCondition(cnd.PowerAmplifier, 3000, card.ID)
		ctx.Match.Chat("Server", fmt.Sprintf("%s was given +3000 power by %s until end of turn", creature.Name, card.Name))
	}
}
