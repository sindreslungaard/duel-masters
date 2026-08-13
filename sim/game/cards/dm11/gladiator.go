package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// WarlordAilzonius ...
func WarlordAilzonius(c *match.Card) {

	c.Name = "Warlord Ailzonius"
	c.Power = 8000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.CantBeSelectedByOpp)

}

// BelixTheExplorer ...
func BelixTheExplorer(c *match.Card) {

	c.Name = "Belix, the Explorer"
	c.Power = 3000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Return a spell from your mana zone to your hand.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
			false,
		).Map(func(spell *match.Card) {
			moved, err := card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s returned %s from their mana zone to their hand", card.Name, spell.Name))
		})
	}))

}

// BaraidTheExplorer ...
func BaraidTheExplorer(c *match.Card) {

	c.Name = "Baraid, the Explorer"
	c.Power = 5000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.SilentSkill(fx.GrantConditionToOwnCreatures(
		cnd.CantBeBlocked,
		nil,
		func(x *match.Card) bool { return x.HasCiv(civ.Light) },
		"can't be blocked",
	)))

}
