package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// RimuelCloudbreakElemental ...
func RimuelCloudbreakElemental(c *match.Card) {

	c.Name = "Rimuel, Cloudbreak Elemental"
	c.Power = 6000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		lightMana := len(fx.FindFilter(
			card.Player,
			match.MANAZONE,
			func(x *match.Card) bool { return x.HasCiv(civ.Light) && !x.Tapped },
		))

		nrCreaturesOpp := len(fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return !x.Tapped },
		))

		toSelect := lightMana
		if toSelect > nrCreaturesOpp {
			toSelect = nrCreaturesOpp
		}

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("Rimuel, Cloudbreak Elemental: select %d of opponent's creatures and tap them", toSelect),
			toSelect,
			toSelect,
			false,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s's %s was tapped by Rimuel, Cloudbreak Elemental", x.Player.Username(), x.Name))
		})
	}))

}

// AlcadeiasLordOfSpirits ...
func AlcadeiasLordOfSpirits(c *match.Card) {

	c.Name = "Alcadeias, Lord of Spirits"
	c.Power = 12500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	installed := false

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.FilterShieldTriggers(ctx, func(x *match.Card) bool { return x.HasCiv(civ.Light) || !x.HasCondition(cnd.Spell) })

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

			p := ctx.Match.CurrentPlayer()

			playedCard, err := p.Player.GetCard(event.CardID, match.HAND)

			if err != nil || !playedCard.HasCondition(cnd.Spell) {
				return
			}
			if !playedCard.HasCiv(civ.Light) {
				ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "Only light spells may be cast while Alcadeias, Lord of Spirits is in the battle zone")
				ctx.InterruptFlow()
			}
		}

		// PlayCardEvent only covers a spell played from hand. A cast that skips
		// it entirely — a shield trigger, or a "Spell Steal" effect such as
		// Bluum Erkis, Flare Guardian forcing a cast of an opponent's spell —
		// only ever fires SpellCast. A plain handler reacting to SpellCast
		// would race the cast spell's own SpellCast handler (whichever card's
		// handlers happen to be iterated first), while a persistent effect
		// always runs before every card handler, so it reliably wins that race
		// regardless of handler order. (A normal, undisguised shield trigger
		// cast never reaches here as a non-light spell: FilterShieldTriggers
		// already keeps it off the offered list before it can be chosen.)
		if !installed {
			installed = true

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				if card.Zone != match.BATTLEZONE {
					installed = false
					exit()
					return
				}

				event, ok := ctx2.Event.(*match.SpellCast)
				if !ok {
					return
				}

				caster := ctx2.Match.Player1.Player
				if event.MatchPlayerID == 2 {
					caster = ctx2.Match.Player2.Player
				}

				castCard, err := caster.GetCard(event.CardID, match.HAND)
				if err != nil || castCard.HasCiv(civ.Light) {
					return
				}

				ctx2.Match.WarnPlayer(caster, "Only light spells may be cast while Alcadeias, Lord of Spirits is in the battle zone")
				ctx2.InterruptFlow()
			})
		}
	})

}

// AerisFlightElemental ...
func AerisFlightElemental(c *match.Card) {

	c.Name = "Aeris, Flight Elemental"
	c.Power = 9000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.AttackCreature)

		if !ok || event.CardID != card.ID {
			return
		}

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.HasCiv(civ.Darkness) && !x.Tapped
			},
		).Map(func(x *match.Card) {
			// don't add if already in the list of attackable creatures
			for _, creature := range event.AttackableCreatures {
				if creature.ID == x.ID {
					return
				}
			}

			event.AttackableCreatures = append(event.AttackableCreatures, x)
		})

	})

}
