package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Called for cards that use Whenever this attacks
// If this function actions results in changes to the blockers, it should go ahead and modifiy it
func WheneverThisAttacks(f func(*match.Card, *match.Context)) match.HandlerFunc {
	return When(Attacking, func(c *match.Card, ctx2 *match.Context) {
		c.AddUniqueSourceCondition(cnd.WheneverThisAttacks, f, c.ID)
	})
}
