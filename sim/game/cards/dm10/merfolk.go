package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TidePatroller ...
func TidePatroller(c *match.Card) {

	c.Name = "Tide Patroller"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker())

}

// MysticMagician ...
func MysticMagician(c *match.Card) {

	c.Name = "Mystic Magician"
	c.Power = 3000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		// "Your creatures that have silent skill are put into the battle zone
		// tapped." Entering play is what matters, not how it happened, so this
		// covers summoning and any effect that puts one into play. A move
		// resets the tapped flag, so the post-move event is the point where
		// tapping it sticks.
		if event, ok := ctx.Event.(*match.CardMoved); ok && event.To == match.BATTLEZONE {
			arrival, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

			if err != nil || !arrival.HasCondition(cnd.SilentSkill) {
				return
			}

			card.Player.TapCard(arrival)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s entered the battle zone tapped because of %s", arrival.Name, card.Name))

			return
		}

		// "Whenever one of your creatures that has silent skill would be
		// destroyed, put it into your hand instead."
		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {
			victim := event.Card

			if victim.Player != card.Player ||
				victim.Zone != match.BATTLEZONE ||
				!victim.HasCondition(cnd.SilentSkill) {
				return
			}

			// Replacing destruction cancels the destruction itself, which is
			// what stops the creature reaching the graveyard.
			ctx.InterruptFlow()

			moved, err := victim.Player.MoveCard(victim.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand instead of being destroyed by %s's effect", victim.Name, card.Player.Username(), card.Name))
		}

	})

}
