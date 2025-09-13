package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SnorkLaShrineGuardian ...
func SnorkLaShrineGuardian(c *match.Card) {

	c.Name = "Snork La, Shrine Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok && event.From == match.MANAZONE && event.To == match.GRAVEYARD && event.Source != "" {

			opp := ctx.Match.Opponent(card.Player)

			sourceCard1, _ := opp.GetCard(event.Source, match.BATTLEZONE)
			sourceCard2, _ := opp.GetCard(event.Source, match.MANAZONE)
			sourceCard3, _ := opp.GetCard(event.Source, match.GRAVEYARD)
			sourceCard4, _ := opp.GetCard(event.Source, match.SHIELDZONE)
			sourceCard5, _ := opp.GetCard(event.Source, match.HAND)

			if sourceCard1 == nil && sourceCard2 == nil && sourceCard3 == nil && sourceCard4 == nil && sourceCard5 == nil {
				return
			}

			manaBurned, _ := card.Player.GetCard(event.CardID, match.GRAVEYARD)

			if manaBurned != nil && fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s's effect: Do you want to return %s to your mana zone?", card.Name, manaBurned.Name)) {
				card.Player.MoveCard(event.CardID, match.GRAVEYARD, match.MANAZONE, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s prevented %s from being discarded from the manazone", card.Name, manaBurned.Name))
			}
		}

	})

}

// GalliaZohlIronGuardianQ ...
func GalliaZohlIronGuardianQ(c *match.Card) {

	c.Name = "Gallia Zohl, Iron Guardian Q"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian, family.Survivor}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				fx.ForceBlocker(x, ctx2, card.ID)
			})
		})
	}))

}
