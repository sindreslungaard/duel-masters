package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CoilingVines ...
func CoilingVines(c *match.Card) {

	c.Name = "Coiling Vines"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ReturnToMana)

}

// PoisonousDahlia ...
func PoisonousDahlia(c *match.Card) {

	c.Name = "Poisonous Dahlia"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantAttackPlayers)

}

// ThornyMandra ...
func ThornyMandra(c *match.Card) {

	c.Name = "Thorny Mandra"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		creatures := match.SearchForCnd(card.Player, ctx.Match, card.Player, match.GRAVEYARD, cnd.Creature, "Thorny Mandra: Select 1 creature from your battlezone that will be sent to your manazone", 1, 1, true)

		for _, creature := range creatures {
			creature.Tapped = false
			card.Player.MoveCard(creature.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved from %s's graveyard to their manazone", creature.Name, card.Player.Username()))
		}

	}))

}
