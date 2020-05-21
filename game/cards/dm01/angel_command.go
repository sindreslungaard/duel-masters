package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HanusaRadianceElemental ...
func HanusaRadianceElemental(c *match.Card) {

	c.Name = "Hanusa, Radiance Elemental"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = family.AngelCommand
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// IocantTheOracle ...
func IocantTheOracle(c *match.Card) {

	c.Name = "Iocant, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = family.AngelCommand
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {
			if ctx.Match.IsPlayerTurn(card.Player) && (event.From == match.BATTLEZONE || event.To == match.BATTLEZONE) {
				iocantTheOracleSpecial(card)
			}
		}

	})

}

// +2000 power amplifier if my battlezone has angel command
func iocantTheOracleSpecial(card *match.Card) {

	card.RemoveConditionBySource(card.ID + "-custom")

	battlezone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return
	}

	hasAngelCommand := false

	for _, creature := range battlezone {
		if creature.Family == family.AngelCommand {
			hasAngelCommand = true
		}
	}

	if hasAngelCommand {
		card.AddCondition(cnd.PowerAmplifier, 2000, card.ID+"-custom")
	}

}

// UrthPurifyingElemental ...
func UrthPurifyingElemental(c *match.Card) {

	c.Name = "Urth, Purifying Elemental"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = family.AngelCommand
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
			card.Tapped = false
		}

	})

}
