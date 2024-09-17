package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Gigagriff(c *match.Card) {

	c.Name = "Gigagriff"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.Slayer, fx.CantAttackPlayers, fx.CantAttackCreatures)
}

func PhantasmalHorrorGigazald(c *match.Card) {

	c.Name = "Phantasmal Horror Gigazald"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = fx.OpponentDiscardsRandomCard

	c.Use(fx.Creature, fx.Evolution, fx.TapAbility,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			fx.GiveTapAbilityToAllies(
				card,
				ctx,
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Darkness },
				fx.OpponentDiscardsRandomCard,
			)

		}),
	)
}
