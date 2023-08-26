package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// VampireSilphy ...
func VampireSilphy(c *match.Card) {

	c.Name = "Vampire Silphy"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				opponent := ctx.Match.Opponent(card.Player)

				myCreatures, err := card.Player.Container(match.BATTLEZONE)
				if err != nil {
					return
				}

				opponentCreatures, err := opponent.Container(match.BATTLEZONE)
				if err != nil {
					return
				}

				toDestroy := make([]*match.Card, 0)

				for _, creature := range myCreatures {
					if ctx.Match.GetPower(creature, false) <= 3000 {
						toDestroy = append(toDestroy, creature)
					}
				}

				for _, creature := range opponentCreatures {
					if ctx.Match.GetPower(creature, false) <= 3000 {
						toDestroy = append(toDestroy, creature)
					}
				}

				for _, creature := range toDestroy {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

			}

		}

	})

}
