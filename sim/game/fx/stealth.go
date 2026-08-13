package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

func Stealth(card *match.Card, ctx *match.Context, civ string) {
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Stealth, civ, card.ID)
	}
}

func FireStealth(card *match.Card, ctx *match.Context) {
	Stealth(card, ctx, civ.Fire)
}

func LightStealth(card *match.Card, ctx *match.Context) {
	Stealth(card, ctx, civ.Light)
}

func DarknessStealth(card *match.Card, ctx *match.Context) {
	Stealth(card, ctx, civ.Darkness)
}

func NatureStealth(card *match.Card, ctx *match.Context) {
	Stealth(card, ctx, civ.Nature)
}

func WaterStealth(card *match.Card, ctx *match.Context) {
	Stealth(card, ctx, civ.Water)
}
