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
		"You may take a Dragon from your deck and put it into the battlezone.",
		1,
		1,
		true,
		func(c *match.Card) bool {
			return c.SharesAFamily(family.Dragons) &&
				fx.CanBeSummoned(card.Player, c)
		},
		true,
	).Map(func(selectedDragon *match.Card) {
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
		cardPlayedCtx := match.NewContext(ctx.Match, &match.CardPlayedEvent{
			CardID: selectedDragon.ID,
		})
		ctx.Match.HandleFx(cardPlayedCtx)

		if !cardPlayedCtx.Cancelled() {

			if !selectedDragon.HasCondition(cnd.Evolution) {
				selectedDragon.AddCondition(cnd.SummoningSickness, nil, nil)
			}

			card.Player.MoveCard(selectedDragon.ID, match.HAND, match.BATTLEZONE, selectedDragon.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to the battle zone by %s's effect", selectedDragon.Name, card.Name))

		}

		// then shuffle deck
		fx.ShuffleDeck(card, ctx, false)
	})

}
