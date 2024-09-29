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

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE && card.Zone != match.HIDDENZONE {
			return
		}
    
		// Can have a full refactoring to fx.When after CardMoved uses card pointer
		if event, ok := ctx.Event.(*match.CardMoved); ok {
			if event.To != match.BATTLEZONE || event.From == match.HIDDENZONE {
				return
			}

			playedCard, err := ctx.Match.Player1.Player.GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				playedCard, err = ctx.Match.Player2.Player.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil {
					return
				}
			}

			if !playedCard.HasCondition(cnd.Evolution) {
				return
			}

			// If Lu Gila is the card evolved, its evolution also becomes tapped according to duel master
			// rulings (https://duelmasters.fandom.com/wiki/Lu_Gila,_Silver_Rift_Guardian/Rulings).
			at := playedCard.Attachments()
			if card.Zone == match.HIDDENZONE && (len(at) == 0 || at[len(at)-1] != card) {
				return
			}

			playedCard.Tapped = true
		}

	})

}

func ArcBinetheAstounding(c *match.Card) {

	c.Name = "Arc Bine, the Astounding"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = arcBinetheAstoundingSpecialAbility

	c.Use(fx.Creature, fx.Evolution, fx.TapAbility,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			fx.GiveTapAbilityToAllies(
				card,
				ctx,
				func(x *match.Card) bool { return x.ID != card.ID && x.Civ == civ.Light },
				arcBinetheAstoundingSpecialAbility,
			)

		}),
	)

}

func arcBinetheAstoundingSpecialAbility(card *match.Card, ctx *match.Context) {
	creatures := fx.Select(
		card.Player,
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		"Select 1 of your opponent's creature and tap it.",
		1,
		1,
		false)

	for _, creature := range creatures {
		creature.Tapped = true
	}
}
