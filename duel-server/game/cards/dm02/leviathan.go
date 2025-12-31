package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingNautilus ...
func KingNautilus(c *match.Card) {

	c.Name = "King Nautilus"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		kingNautilusSpecial(card, ctx)

	})

}

func kingNautilusSpecial(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackConfirmed); ok {

		p := ctx.Match.CurrentPlayer()

		creature, err := p.Player.GetCard(event.CardID, match.BATTLEZONE)

		if err != nil {
			return
		}

		if !creature.HasFamily(family.LiquidPeople) {
			return
		}

		creature.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)

	}

}
