package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ForbosSanctumGuardianQ(c *match.Card) {

	c.Name = "Forbos, Sanctum Guardian Q"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.MySurvivorSummoned, fx.SearchDeckTake1Spell))
}

func LuGilaSilverRiftGuardian(c *match.Card) {

	c.Name = "Lu Gila, Silver Rift Guardian"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					exit()
					return
				}

				// Can have a full refactoring to fx.When after CardMoved uses card pointer
				if event, ok := ctx2.Event.(*match.CardMoved); ok {
					if event.To != match.BATTLEZONE || event.From == match.HIDDENZONE {
						return
					}

					playedCard, err := ctx2.Match.Player1.Player.GetCard(event.CardID, match.BATTLEZONE)
					if err != nil {
						playedCard, err = ctx2.Match.Player2.Player.GetCard(event.CardID, match.BATTLEZONE)
						if err != nil {
							return
						}
					}

					if !playedCard.HasCondition(cnd.Evolution) {
						return
					}

					playedCard.Tapped = true
				}
			})
		}),
	)
}

func ArcBinetheAstounding(c *match.Card) {

	c.Name = "Arc Bine, the Astounding"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = fx.TapOpCreature

	c.Use(fx.Creature, fx.Evolution, fx.TapAbility,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			fx.GiveTapAbilityToAllies(
				card,
				ctx,
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Light },
				fx.TapOpCreature,
			)

		}),
	)

}
