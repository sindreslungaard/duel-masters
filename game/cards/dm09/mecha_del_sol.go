package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// PetrovaChannelerOfSuns ...
func PetrovaChannelerOfSuns(c *match.Card) {

	c.Name = "Petrova, Channeler of Suns"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeSelectedByOpp,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					exit()
				}

				//TODO choose a race other than Mecha del Sol
				// each creature of that race has +4000 power (STATIC ABILITY)
				// with unique source condition with card.ID
			})
		}),
	)

}

// BruiserDragon ...
func BruiserDragon(c *match.Card) {

	c.Name = "Bruiser Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Put 1 of your shields into your graveyard.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: One of your shields was put into graveyard", card.Name))
		})
	}))

}
