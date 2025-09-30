package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MessaBanhaExpanseGuardian ...
func MessaBanhaExpanseGuardian(c *match.Card) {

	c.Name = "Messa Banha, Expanse Guardian"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.Block); ok {
			if event.Attacker.Player == card.Player ||
				event.Attacker.Zone != match.BATTLEZONE ||
				card.Zone != match.BATTLEZONE {
				return
			}

			//@TODO test this
			if len(event.Blockers) > 0 && event.Attacker.IsBlockable(ctx) {
				for _, b := range event.Blockers {
					if b.ID == card.ID {
						// Force the battle between the attacker and this card
						ctx.InterruptFlow()
						ctx.Match.Battle(event.Attacker, card, true, len(event.ShieldsAttacked) > 0)
						return
					}
				}
			}
		}
	})

}

// PalaOlesisMorningGuardian ...
func PalaOlesisMorningGuardian(c *match.Card) {

	c.Name = "Pala Olesis, Morning Guardian"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})

					exit()
					return
				}

				if !ctx2.Match.IsPlayerTurn(card.Player) {
					fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.ID != card.ID
						},
					).Map(func(x *match.Card) {
						x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
					})
				} else {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})
				}
			})
		}))

}
