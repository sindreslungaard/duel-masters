package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Klujadras ...
func Klujadras(c *match.Card) {

	c.Name = "Klujadras"
	c.Power = 4000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)

		// Each player counts their own wave strikers, so the totals are read
		// before either player draws and neither can affect the other.
		mine := fx.WaveStrikersInBattleZone(card.Player)
		theirs := fx.WaveStrikersInBattleZone(opponent)

		card.Player.DrawCards(mine)
		opponent.DrawCards(theirs)

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s: %s drew %d card(s) and %s drew %d card(s)", card.Name, card.Player.Username(), mine, opponent.Username(), theirs))
	})))

}
