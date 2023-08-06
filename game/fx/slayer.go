package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

type SlayerCondition func(target *match.Card) bool

// Slayer destroys the source card when the card is destroyed
func Slayer(card *match.Card, ctx *match.Context) {
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Slayer, nil, card.ID)
	}
}

func ConditionalSlayer(condition SlayerCondition) func(card *match.Card, ctx *match.Context) {
	return func(card *match.Card, ctx *match.Context) {
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddCondition(cnd.Slayer, condition, card.ID)
		}
	}
}

// Suicide destroys the card when it wins a battle
func Suicide(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.Battle); ok {
		if event.Attacker == card || event.Defender == card {
			ctx.ScheduleAfter(func() {
				// Still in the battlezone so it won the battle
				creature, err := card.Player.GetCard(card.ID, match.BATTLEZONE)
				if err != nil {
					return
				}
				ctx.Match.Destroy(creature, creature, match.DestroyedBySlayer)
			})
		}

	}

}
