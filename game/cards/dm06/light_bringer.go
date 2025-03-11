package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

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

	c.Use(fx.Creature,
		fx.When(fx.IHaveCastASpell, func(card *match.Card, ctx *match.Context) { spellcast = true }),
		func(card *match.Card, ctx *match.Context) {

			// reset spellcast flag at the end of the turn
			_, ok := ctx.Event.(*match.EndStep)
			if ok {
				spellcast = false
			}

			if event, ok := ctx.Event.(*match.PlayCardEvent); ok {
				if event.CardID == card.ID {
					if !spellcast {
						ctx.InterruptFlow()
						ctx.Match.WarnPlayer(card.Player, "You can summon this creature only if you have cast a spell this turn")
						return
					}
				}

			}
		})
}

func AdomisTheOracle(c *match.Card) {

	c.Name = "Adomis, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)
		cards["Your shields"] = fx.Find(c.Player, match.SHIELDZONE)
		cards["Opponent's shields"] = fx.Find(ctx.Match.Opponent(c.Player), match.SHIELDZONE)

		fx.SelectMultipartBackside(
			card.Player,
			ctx.Match,
			cards,
			"Select 1 shield that will be shown to you",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.ShowCards(
				card.Player,
				"The shield is:",
				[]string{x.ImageID},
			)
		})

	}

	c.Use(fx.Creature, fx.TapAbility)
}
