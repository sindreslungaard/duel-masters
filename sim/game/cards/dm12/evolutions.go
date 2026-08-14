package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// The multi-family evolutions of the set. Each carries both of the races it can
// be put on, which is what lets fx.Evolution accept either as a base, and each
// lifts its own kin by 2000 while it is in the battle zone.

// AgiraTheWarlordCrawler ...
func AgiraTheWarlordCrawler(c *match.Card) {

	c.Name = "Agira, the Warlord Crawler"
	c.Power = 5500
	c.Civs = []string{civ.Light, civ.Water}
	c.Family = []string{family.Gladiator, family.EarthEater}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light, civ.Water}

	kin := []string{family.Gladiator, family.EarthEater}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.PowerAmplifier, 2000, fx.OtherCreaturesSharingAFamily(kin))),
		func(card *match.Card, ctx *match.Context) {
			event, ok := ctx.Event.(*match.Battle)

			if !ok || card.Zone != match.BATTLEZONE || !event.Blocked {
				return
			}

			blocker := event.Defender

			if blocker.Player != card.Player || !blocker.SharesAFamily(kin) {
				return
			}

			if !fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s's effect: %s blocked. Do you want to draw a card?", card.Name, blocker.Name)) {
				return
			}

			card.Player.DrawCards(1)
		})

}

// CometEyeTheSpectralSpud ...
func CometEyeTheSpectralSpud(c *match.Card) {

	c.Name = "Comet Eye, The Spectral Spud"
	c.Power = 5500
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.WildVeggies, family.RainbowPhantom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	kin := []string{family.WildVeggies, family.RainbowPhantom}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.PowerAmplifier, 2000, fx.OtherCreaturesSharingAFamily(kin))),
		fx.When(fx.EndOfMyTurnCreatureBZ, func(card *match.Card, ctx *match.Context) {
			tapped := fx.FindFilter(card.Player, match.BATTLEZONE, func(x *match.Card) bool {
				return x.Tapped && x.SharesAFamily(kin)
			})

			if len(tapped) < 1 {
				return
			}

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may untap any number of your %s and %s.", card.Name, family.WildVeggies, family.RainbowPhantom),
				1,
				len(tapped),
				true,
				func(x *match.Card) bool { return x.Tapped && x.SharesAFamily(kin) },
				false,
			).Map(func(creature *match.Card) {
				creature.Tapped = false
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was untapped by %s", creature.Name, card.Name))
			})
		}))

}

// HydroozeTheMutantEmperor ...
func HydroozeTheMutantEmperor(c *match.Card) {

	c.Name = "Hydrooze, the Mutant Emperor"
	c.Power = 5000
	c.Civs = []string{civ.Water, civ.Darkness}
	c.Family = []string{family.CyberLord, family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	kin := []string{family.CyberLord, family.Hedrian}

	// "Your Cyber Lords or Hedrians can't be blocked" has no "other", so it
	// covers Hydrooze itself as well.
	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.PowerAmplifier, 2000, fx.OtherCreaturesSharingAFamily(kin))),
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.CantBeBlocked, nil, fx.CreaturesSharingAFamily(kin))))

}

// PhantomachTheGigatrooper ...
func PhantomachTheGigatrooper(c *match.Card) {

	c.Name = "Phantomach, the Gigatrooper"
	c.Power = 6000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Chimera, family.Armorloid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	kin := []string{family.Chimera, family.Armorloid}

	// The power bonus is for its "other" kin, but the double breaker is for
	// "each" of them, so Phantomach gets that one itself.
	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.PowerAmplifier, 2000, fx.OtherCreaturesSharingAFamily(kin))),
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.DoubleBreaker, true, fx.CreaturesSharingAFamily(kin))))

}

// NemonexBajulasRobomantis ...
func NemonexBajulasRobomantis(c *match.Card) {

	c.Name = "Nemonex, Bajula's Robomantis"
	c.Power = 5000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.Xenoparts, family.GiantInsect}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	kin := []string{family.Xenoparts, family.GiantInsect}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.PowerAmplifier, 2000, fx.OtherCreaturesSharingAFamily(kin))),
		func(card *match.Card, ctx *match.Context) {
			if card.Zone != match.BATTLEZONE {
				return
			}

			// "Attacking and isn't blocked" has two shapes in this engine: an
			// unblocked attack on a player reaches the shields, and an
			// unblocked attack on a creature produces a battle that was not
			// blocked. The two are mutually exclusive, so neither double fires.
			attacker := (*match.Card)(nil)

			switch event := ctx.Event.(type) {
			case *match.BreakShieldEvent:
				attacker = event.Source
			case *match.Battle:
				if !event.Blocked {
					attacker = event.Attacker
				}
			}

			if attacker == nil || attacker.Player != card.Player || !attacker.SharesAFamily(kin) {
				return
			}

			fx.OpponentChoosesManaBurn(card, ctx)
		})

}
