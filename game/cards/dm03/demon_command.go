package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GirielGhastlyWarrior ...
func GirielGhastlyWarrior(c *match.Card) {

	c.Name = "Giriel, Ghastly Warrior"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// GamilKnightOfHatred ...
func GamilKnightOfHatred(c *match.Card) {

	c.Name = "Gamil, Knight of Hatred"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {
		
		ctx.ScheduleAfter(func() {
			creatures := match.Filter(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "Select 1 of your darkness creatures from the graveyard that will be returned to your hand", 0, 1, true, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && x.Civ == civ.Darkness })
					
			for _, creature := range creatures {
				card.Player.MoveCard(creature.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their  graveyard", creature.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
			}
		})
	}))

}