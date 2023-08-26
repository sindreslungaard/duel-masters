package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Draglide ...
func Draglide(c *match.Card) {

	c.Name = "Draglide"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack)

}

// GatlingSkyterror ...
func GatlingSkyterror(c *match.Card) {

	c.Name = "Gatling Skyterror"
	c.Power = 7000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.AttackUntapped)

}

// ScarletSkyterror ...
func ScarletSkyterror(c *match.Card) {

	c.Name = "Scarlet Skyterror"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.BATTLEZONE {
				return
			}

			blockers := make([]*match.Card, 0)

			myBattlezone, err := card.Player.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			opponentBattlezone, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			for _, creature := range myBattlezone {
				if creature.HasCondition(cnd.Blocker) {
					blockers = append(blockers, creature)
				}
			}

			for _, creature := range opponentBattlezone {
				if creature.HasCondition(cnd.Blocker) {
					blockers = append(blockers, creature)
				}
			}

			for _, blocker := range blockers {
				ctx.Match.Destroy(blocker, card, match.DestroyedByMiscAbility)
			}

		}

	})

}
