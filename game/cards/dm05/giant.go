package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"github.com/sirupsen/logrus"
)

// AvalancheGiant ...
func AvalancheGiant(c *match.Card) {

	c.Name = "Avalanche Giant"
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = family.Giant
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.CantAttackCreatures, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if !event.Blocked || event.Attacker != card {
				return
			}
			opponent := ctx.Match.Opponent(card.Player)

			shieldzone, err := opponent.Container(match.SHIELDZONE)
	
			if err != nil {
				return
			}
			if len(shieldzone) > 0 {
				for {

					ctx.Match.NewBacksideAction(card.Player, shieldzone, 1, 1, "Avalanche Giant ability: select shield to break", false)
					action := <-card.Player.Action

					if len(action.Cards) != 1 || !match.AssertCardsIn(shieldzone, action.Cards[0]) {
						ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
						continue
					}
					shieldsToBreak := make([]*match.Card, 0)
					for _, cardID := range action.Cards {
						shield, err := opponent.GetCard(cardID, match.SHIELDZONE)
						if err != nil {
							logrus.Debug("Could not find specified shield in shieldzone")
							continue
						}
						shieldsToBreak = append(shieldsToBreak, shield)
					}
					ctx.Match.BreakShields(shieldsToBreak)
					ctx.Match.CloseAction(card.Player)
					break
				}
			}
		}
	})

}
