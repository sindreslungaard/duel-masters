package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func FrostSpecterShadowOfAge(c *match.Card) {

	c.Name = "Frost Specter, Shadow of Age"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				// remove cards with current buffs
				getGhostCreatures(card, ctx).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return

			}

			getGhostCreatures(card, ctx).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})

		})

	}))

}

func getGhostCreatures(card *match.Card, ctx *match.Context) fx.CardCollection {

	ghostCreatures := fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool { return x.HasFamily(family.Ghost) },
	)

	ghostCreatures = append(ghostCreatures,

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasFamily(family.Ghost) },
		)...,
	)

	return ghostCreatures
}

func GrimSoulShadowOfReversal(c *match.Card) {

	c.Name = "Grim Soul, Shadow of Reversal"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Darkness },
		).Map(func(x *match.Card) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.GRAVEYARD,
				"Grim Soul, Shadow of Reversal: You may return a darkness creature from your graveyard to your hand",
				1,
				1,
				true,
				func(x *match.Card) bool { return x.Civ == civ.Darkness },
			).Map(func(x *match.Card) {
				card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their graveyard by Grim Soul, Shadow of Reversal", x.Name, card.Player.Username()))
				card.Tapped = true
			})
		})
	}))
}

func LoneTearShadowOfSolitude(c *match.Card) {

	c.Name = "LoneTearShadowOfSolitude"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		if len(fx.Find(card.Player, match.BATTLEZONE)) == 1 {
			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s exorcised itself", c.Name))
		}

	}))
}
