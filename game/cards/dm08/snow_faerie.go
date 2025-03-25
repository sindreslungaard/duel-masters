package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KachuaKeeperOfTheIcegate ...
func KachuaKeeperOfTheIcegate(c *match.Card) {

	c.Name = "Kachua, Keeper of the Icegate"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = kachuaKeeperOfTheIcegateTapAbility

	c.Use(fx.Creature, fx.TapAbility)
}

// Currently this is bugged:
// https://duelmasters.fandom.com/wiki/Kachua,_Keeper_of_the_Icegate/Rulings
// Bazagazeal case: you should be able to choose from this creature's effect or
// the card itself's effect (Bazagazeal)
func kachuaKeeperOfTheIcegateTapAbility(card *match.Card, ctx *match.Context) {

	// Search your deck. You MAY take a Dragon from your deck
	myDragons := fx.SelectFilter(
		card.Player,
		ctx.Match,
		card.Player,
		match.DECK,
		"You may take a Dragon from your deck and put it into the battlezone.",
		1,
		1,
		true,
		func(c *match.Card) bool {
			return c.SharesAFamily(family.Dragons) &&
				(c.HasCondition(cnd.Creature) &&
					((!c.HasCondition(cnd.Evolution) &&
						!c.HasCondition(cnd.EvolveIntoAnyFamily)) ||
						(c.HasCondition(cnd.EvolveIntoAnyFamily) ||
							(c.HasCondition(cnd.Evolution) &&
								len(fx.FindFilter(
									card.Player,
									match.BATTLEZONE,
									func(c2 *match.Card) bool {
										return c2.SharesAFamily(c.Family)
									},
								)) > 0))))
		},
		true,
	)

	if len(myDragons) > 0 {

		selectedDragon := myDragons[0]

		// that creature has "speed attacker" + at the end of the turn, destroy that creature
		selectedDragon.Use(fx.When(fx.InTheBattlezone, func(selDragon *match.Card, ctx2 *match.Context) {
			ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit func()) {
				if selDragon.Zone != match.BATTLEZONE {
					exit()
					return
				}

				if _, ok := ctx3.Event.(*match.EndStep); ok {
					ctx3.Match.Destroy(selDragon, selDragon, match.DestroyedByMiscAbility)
					exit()
					return
				}

				// TODO will be refactored when PR #314 will be merged (https://github.com/sindreslungaard/duel-masters/pull/314)
				selDragon.RemoveCondition(cnd.SummoningSickness)
			})
		}))

		// and put into the Battlezone
		ctx.Match.MoveCard(selectedDragon, match.BATTLEZONE, card)

		// then shuffle deck
		fx.ShuffleDeck(card, ctx, false)

	}

}
