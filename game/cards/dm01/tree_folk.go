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
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: You may put 1 creature from your graveyard to your mana zone", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			false,
		).Map(func(x *match.Card) {
			x.Tapped = false
			ctx.Match.BroadcastState()
			x.Player.MoveCard(x.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved from %s's graveyard to their manazone", x.Name, x.Player.Username()))
		})
	}))

}
