package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LenaVizierOfBrilliance ...
func LenaVizierOfBrilliance(c *match.Card) {

	c.Name = "Lena, Vizier of Brilliance"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		spells := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "You may select 1 spell from your mana zone that will be sent to your hand", 0, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

		for _, spell := range spells {

			card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", spell.Player.Username(), spell.Name))

		}
	}))

}

// SiegBaliculaTheIntense ...
func SiegBaliculaTheIntense(c *match.Card) {

	c.Name = "Sieg Balicula, the Intense"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

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
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Light },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
			})

		})

	}))

}
