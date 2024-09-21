package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func BattleshipMutant(c *match.Card) {

	c.Name = "Battleship Mutant"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(card.Player, match.BATTLEZONE,
			func(c *match.Card) bool { return c.Civ == civ.Darkness },
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
			x.AddCondition(cnd.DoubleBreaker, true, card.ID)
			x.AddCondition(cnd.DestroyAfterBattle, true, card.ID)
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}
