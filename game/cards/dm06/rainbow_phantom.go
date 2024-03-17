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
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
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
			card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand using %s's tap ability", spell.Player.Username(), spell.Name, card.Name))
			card.Tapped = true
		})
	}))

}

func MoontearSpectralKnight(c *match.Card) {

	c.Name = "Moontear, Spectral Knight"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 2
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
