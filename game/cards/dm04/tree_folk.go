package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SupportingTulip ...
func SupportingTulip(c *match.Card) {

	c.Name = "Supporting Tulip"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			supportingTulipSpecial(card, ctx, event.CardID)
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			supportingTulipSpecial(card, ctx, event.CardID)
		}

	})

}

func supportingTulipSpecial(card *match.Card, ctx *match.Context, cardID string) {

	p := ctx.Match.CurrentPlayer()

	creature, err := p.Player.GetCard(cardID, match.BATTLEZONE)

	if err != nil {
		return
	}

	if !creature.HasFamily(family.AngelCommand) {
		return
	}

	if ctx.Match.IsPlayerTurn(card.Player) {
		creature.AddUniqueSourceCondition(cnd.PowerAttacker, 4000, card.ID)
	}
}

// ExplodingCactus ...
func ExplodingCactus(c *match.Card) {

	c.Name = "Exploding Cactus"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.Civ == civ.Light }) {
			return 2000
		}

		return 0
	}

	c.Use(fx.Creature)

}
