package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TauntingSkyterror ...
func TauntingSkyterror(c *match.Card) {

	c.Name = "Taunting Skyterror"
	c.Civ = civ.Fire
	c.Power = 3000
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if ctx2.Match.IsPlayerTurn(card.Player) || !card.Tapped {
				return
			}

			fx.Find(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				fx.ForceAttack(x, ctx2)
			})
		})
	}))

}
