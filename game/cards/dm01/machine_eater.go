package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ArtisanPicora ...
func ArtisanPicora(c *match.Card) {

	c.Name = "Artisan Picora"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.MachineEater}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Artisan Picora: Choose 1 card from your mana zone that will be sent to your graveyard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was sent to %s's graveyard by Artisan Picora", x.Name, x.Player.Username()))
		})

	}))

}

// NomadHeroGigio ...
func NomadHeroGigio(c *match.Card) {

	c.Name = "Nomad Hero Gigio"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.MachineEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.AttackUntapped)

}
