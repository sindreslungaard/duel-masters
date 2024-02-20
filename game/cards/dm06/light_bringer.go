package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func AdomisTheOracle(c *match.Card) {

	c.Name = "Adomis, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {

		shields := fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			"Adomis, the Oracle: Select one of your shields that will be shown to you",
			1,
			1,
			true,
		)

		ids := make([]string, 0)

		for _, s := range shields {
			ids = append(ids, s.ImageID)
		}

		ctx.Match.ShowCards(
			card.Player,
			"Your shield:",
			ids,
		)

		card.Tapped = true

	}))
}

func VessTheOracle(c *match.Card) {

	c.Name = "Vess, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)
}

func YulukTheOracle(c *match.Card) {

	c.Name = "Yuluk, the Oracle"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	spellcast := false

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

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
