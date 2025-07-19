package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
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
	fx.SelectFilter(
		card.Player,
		ctx.Match,
		card.Player,
		match.DECK,
		fmt.Sprintf("%s: You may take a Dragon from your deck and put it into the battlezone, then shuffle your deck.", card.Name),
		1,
		1,
		true,
		func(x *match.Card) bool {
			return x.SharesAFamily(family.Dragons) &&
				fx.CanBeSummoned(card.Player, x)
		},
		true,
	).Map(func(selDragon *match.Card) {
		// put into the Battlezone
		fx.ForcePutCreatureIntoBZ(ctx, selDragon, match.DECK, card)

		// that creature has "speed attacker" + at the end of the turn, destroy that creature
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if selDragon.Zone != match.BATTLEZONE {
				selDragon.RemoveConditionBySource(card.ID)
				exit()
				return
			}

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				selDragon.RemoveConditionBySource(card.ID)
				ctx2.Match.Destroy(selDragon, selDragon, match.DestroyedByMiscAbility)
				exit()
				return
			}

			selDragon.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
		})

		// then shuffle deck
		fx.ShuffleDeck(card, ctx, false)
	})

}
