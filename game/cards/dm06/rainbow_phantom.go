package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func CosmogoldSpectralKnight(c *match.Card) {

	c.Name = "Cosmogold, Spectral Knight"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.SelectFilterSelectablesOnly(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Select 1 spell from your mana zone that will be sent to your hand",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
		).Map(func(spell *match.Card) {
			card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s retrieved %s from the mana zone to their hand using %s's tap ability", spell.Player.Username(), spell.Name, card.Name))
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}

func MoontearSpectralKnight(c *match.Card) {

	c.Name = "Moontear, Spectral Knight"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 2
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
						ctx.Match.WarnPlayer(card.Player, "You can summon this creature only if you have cast a spell this round")
						return
					}
				}

			}

		})

}
