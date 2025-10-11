package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SupersonicJetpack ...
func SupersonicJetpack(c *match.Card) {

	c.Name = "Supersonic Jetpack"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: One of your creatures in the battlezone gets 'speed attacker' until end of turn.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
		})
	}))

}
