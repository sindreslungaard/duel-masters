package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SeaSlug ...
func SeaSlug(c *match.Card) {

	c.Name = "Sea Slug"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantBeBlocked)

}

// SplitHeadHydroturtleQ ...
func SplitHeadHydroturtleQ(c *match.Card) {

	c.Name = "Split-Head Hydroturtle Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish, family.Survivor}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		event, ok := ctx.Event.(*match.AttackConfirmed)

		if !ok {
			return
		}

		creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

		if err != nil {
			return
		}

		if creature.HasCondition(cnd.Survivor) {
			fx.MayDraw1(card, ctx)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s drew a card when %s attacked due to %s's survivor ability", card.Player.Username(), creature.Name, card.Name))
		}

	})

}

// LurkingEel ...
func LurkingEel(c *match.Card) {

	c.Name = "Lurking Eel"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ConditionalBlocker(func(target *match.Card) bool {
		return target.Civ == civ.Fire || target.Civ == civ.Nature
	}))

}
