package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ColossusBoost ...
func ColossusBoost(c *match.Card) {

	c.Name = "Colossus Boost"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: One of your creatures in the battle zone gets '+4000 Power' until the end of the turn.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				x.AddUniqueSourceCondition(cnd.PowerAmplifier, 4000, card.ID)
			})
		})
	}))

}

// ForcedFrenzy ...
func ForcedFrenzy(c *match.Card) {

	c.Name = "Forced Frenzy"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Find(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.StartOfTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
					return
				}

				fx.ForceAttack(x, ctx2)
			})
		})
	}))

}
