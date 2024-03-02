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

		card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND)
		ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to the %s's hand", c.Name, c.Player.Username()))
	}))
}

/*
// LavaWalkerExecuto
func LavaWalkerExecuto(c *match.Card) {

	c.Name = "Lava Walker Executo"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(
		fx.Creature,
		fx.Evolution,
		fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})

					exit()
					return
				}

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Fire },
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.TapAbility, true, card.ID)
				})

			})

		}),
		fx.When(
			func(card *match.Card, ctx *match.Context) bool {
				if _, ok := ctx.Event.(*match.TapAbility); ok {

					if card.Player != c.Player && card.Civ != civ.Fire {
						return false
					}

					if card.HasCondition(cnd.SummoningSickness) {
						ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't use tap ability because it has summoning sickness", card.Name))
						return false
					}

					if card.Tapped {
						ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't use tap ability because it is already tapped", card.Name))
						return false
					}

					return true

				}
				return false

			},

			func(card *match.Card, ctx *match.Context) {

				ctx.Match.Chat("Server", fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))
				creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain +3000 Power", 1, 1, false)
				for _, creature := range creatures {

					creature.AddCondition(cnd.PowerAmplifier, 3000, card.ID)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was given +3000 power by %s until end of turn", creature.Name, card.Name))

					card.Tapped = true
				}

			},
		),
	)
}
*/
