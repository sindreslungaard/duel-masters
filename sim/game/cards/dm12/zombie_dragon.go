package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NecrodragonJagraveen ...
func NecrodragonJagraveen(c *match.Card) {

	c.Name = "Necrodragon Jagraveen"
	c.Power = 6000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker(), fx.Doublebreaker, fx.When(fx.Blocks, func(card *match.Card, ctx *match.Context) {
		// Only blocking costs it its life. The condition is read by the battle
		// that is already resolving, so setting it here is in time.
		card.AddUniqueSourceCondition(cnd.DestroyAfterBattle, true, card.ID)
	}))

}

// NecrodragonZalva ...
func NecrodragonZalva(c *match.Card) {

	c.Name = "Necrodragon Zalva"
	c.Power = 5000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		ctx.Match.Opponent(card.Player).DrawCards(1)
	}))

}
