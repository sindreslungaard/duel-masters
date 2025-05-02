package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func MightyBanditAceOfThieves(c *match.Card) {

	c.Name = "Mighty Bandit, Ace of Thieves"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = MightyBanditAceOfThievesTapAbility

	c.Use(fx.Creature, fx.TapAbility)
}

func MightyBanditAceOfThievesTapAbility(card *match.Card, ctx *match.Context) {
	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))
	fx.Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Select 1 creature from your battlezone that will gain +5000 power", card.Name),
		1,
		1,
		false,
	).Map(func(x *match.Card) {
		x.AddCondition(cnd.PowerAmplifier, 5000, card.ID)
		ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given +5000 power by %s until end of turn", x.Name, card.Name))
	})

}

func InnocentHunterBladeOfAll(c *match.Card) {

	c.Name = "Innocent Hunter, Blade of All"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.EvolveIntoAnyFamily)

}
