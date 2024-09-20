package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// Called for cards that use Whenever this attacks
// If this function actions results in changes to the blockers, it should go ahead and modifiy it
func WheneverThisAttacks(f func(*match.Card, *match.Context)) match.HandlerFunc {
	return When(Attacking, func(c *match.Card, ctx2 *match.Context) {
		c.AddUniqueSourceCondition(cnd.WheneverThisAttacks, f, c.ID)
	})
}

func WheneverThisAttacksMayTapDorFCreature() match.HandlerFunc {
	return WheneverThisAttacks(func(c *match.Card, ctx *match.Context) {
		filter := func(x *match.Card) bool { return x.Civ == civ.Fire || x.Civ == civ.Darkness }
		cards := make(map[string][]*match.Card)
		cards["Your creatures"] = FindFilter(c.Player, match.BATTLEZONE, filter)
		cards["Opponent's creatures"] = FindFilter(ctx.Match.Opponent(c.Player), match.BATTLEZONE, filter)

		SelectMultipart(
			c.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Select a card to tap or close to cancel", c.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = true
			RemoveBlockerFromList(x, ctx)
		})

	})
}

// TODO: This is currently nerfed because we make shield selection to early in the attack.
// In theory with this power they should make a more informed attack
func WheneverThisAttacksMayLookAtOpShield() match.HandlerFunc {
	return WheneverThisAttacks(func(card *match.Card, ctx *match.Context) {
		SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's shields that will be shown to you", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.ShowCards(
				card.Player,
				"Your opponent's shield:",
				[]string{x.ImageID},
			)
		})
	})
}
