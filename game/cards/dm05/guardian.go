package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SnorkLaShrineGuardian ...
func SnorkLaShrineGuardian(c *match.Card) {

	c.Name = "Snork La, Shrine Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok && event.From == match.MANAZONE && event.To == match.GRAVEYARD {

			card.Player.MoveCard(event.CardID, match.GRAVEYARD, match.MANAZONE)
			ctx.Match.Chat("Server", "Snork La, Shrine Guardian prevented card from being discarded from the manazone")

		}

	})

}

// GalliaZohlIronGuardianQ ...
func GalliaZohlIronGuardianQ(c *match.Card) {

	c.Name = "Gallia Zohl, Iron Guardian Q"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return

			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
			})

		})

	}))

}
