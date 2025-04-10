package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// EarthstompGiant ...
func EarthstompGiant(c *match.Card) {

	c.Name = "Earthstomp Giant"
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		manaZone, err := card.Player.Container(match.MANAZONE)

		if err != nil {
			return
		}

		for _, manaZoneCard := range manaZone {
			if manaZoneCard.HasCondition(cnd.Creature) {
				manaZoneCard.Player.MoveCard(manaZoneCard.ID, match.MANAZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to %s's hand from manazone by %s's effect", manaZoneCard.Name, manaZoneCard.Player.Username(), card.Name))
			}
		}
	}))
}

// DawnGiant ...
func DawnGiant(c *match.Card) {

	c.Name = "Dawn Giant"
	c.Power = 11000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.CantAttackCreatures)
}
