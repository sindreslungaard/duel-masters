package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func ArmoredTransportGaliacruse(c *match.Card) {

	c.Name = "Armored Transport Galiacruse"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(card *match.Card) bool { return card.Civ == civ.Fire },
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.AttackUntapped, nil, card)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s can attack untapped creatures this turn", x.Name))
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}

func OtherworldlyWarriorNaglu(c *match.Card) {

	c.Name = "Otherworldly Warrior Naglu"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker3000, fx.Doublebreaker, fx.CantBeAttacked)

}
