package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AzaghastTyrantOfShadows ...
func AzaghastTyrantOfShadows(c *match.Card) {

	c.Name = "Azaghast, Tyrant of Shadows"
	c.Power = 9000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DarkLord}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker,
		fx.When(fx.AnotherOwnGhostSummoned, func(card *match.Card, ctx *match.Context) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: you may destroy 1 of your opponent's untapped creatures.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return !x.Tapped
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		}))

}

// GabzagulWarlordOfPain ...
func GabzagulWarlordOfPain(c *match.Card) {

	c.Name = "Gabzagul, Warlord of Pain"
	c.Power = 5000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DarkLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			// ForceAttack only holds back the turn player's own creatures, so
			// sweeping both battle zones covers "each creature" on either turn.
			for _, player := range []*match.Player{card.Player, ctx2.Match.Opponent(card.Player)} {
				for _, creature := range fx.Find(player, match.BATTLEZONE) {
					if ctx2.Cancelled() {
						// HandleFx checks ctx.cancel between card handlers, but not
						// between persistent effects and never inside one. Once a
						// creature has blocked the end of the turn, warning again for
						// every other able creature would spam the player.
						return
					}

					fx.ForceAttack(creature, ctx2)
				}
			}
		})
	}))

}
