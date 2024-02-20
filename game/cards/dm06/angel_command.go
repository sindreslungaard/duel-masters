package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func GarielElementalOfSunbeams(c *match.Card) {

	c.Name = "Gariel, Elemental Of Sunbeams"
	c.Power = 7500
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	spellcast := false

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if _, ok := ctx2.Event.(*match.SpellCast); ok {
				spellcast = true
			}

			// remove persistent effect when turn ends
			_, ok := ctx2.Event.(*match.EndStep)
			if ok {
				exit()
				spellcast = false
			}

		})

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {
			if event.CardID == card.ID {
				if !spellcast {
					ctx.InterruptFlow()
					ctx.Match.WarnPlayer(card.Player, "You can summon this creature only if you have cast a spell this round")
					return
				}
			}

		}

	})

}
