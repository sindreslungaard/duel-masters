package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// VreemahFreakyMojoTotem ...
func VreemahFreakyMojoTotem(c *match.Card) {

	c.Name = "Vreemah, Freaky Mojo Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature,
		fx.When(fx.AnotherOwnCreatureSummoned, func(card *match.Card, ctx *match.Context) {
			beastFolks := fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family.BeastFolk) && !x.HasCondition(card.ID)
				},
			)

			beastFolks = append(beastFolks, fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family.BeastFolk) && !x.HasCondition(card.ID)
				},
			)...)

			beastFolks.Map(func(x *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if x.Zone != match.BATTLEZONE {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					x.AddUniqueSourceCondition(card.ID, nil, card.ID)
					x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
					x.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
				})
			})
		}))
}

// WhisperingTotem ...
func WhisperingTotem(c *match.Card) {

	c.Name = "Whispering Totem"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature,
		fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.DECK,
				fmt.Sprintf("%s's effect: Search your deck. You may take a Whispering Totem from your deck, show that creature to your opponent, and put it into your hand. Then shuffle your deck.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return x.Name == card.Name
				},
				true,
			).Map(func(x *match.Card) {
				card.Player.MoveCard(x.ID, match.DECK, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's hand from his deck due to %s's effect.", x.Name, card.Player.Username(), card.Name))
			})
		}))
}
