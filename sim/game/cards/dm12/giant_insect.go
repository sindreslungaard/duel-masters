package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// pincerScarabPowerPerCard is what each card in the opponent's hand is worth to
// Pincer Scarab.
const pincerScarabPowerPerCard = 2000

// insectDoubleBreakerAt is the power both of this set's growing Giant Insects
// need before they break two shields.
const insectDoubleBreakerAt = 6000

// PincerScarab ...
func PincerScarab(c *match.Card) {

	c.Name = "Pincer Scarab"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		hand, err := m.Opponent(c.Player).Container(match.HAND)

		if err != nil {
			return 0
		}

		return len(hand) * pincerScarabPowerPerCard
	}

	// The breaker follows the power rather than the hand directly, so a bonus
	// from anywhere else counts towards it just the same.
	c.Use(fx.Creature, fx.PowerBreakerTiers(insectDoubleBreakerAt, 0))

}

// CopperLocust ...
func CopperLocust(c *match.Card) {

	c.Name = "Copper Locust"
	c.Power = 5000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	// "A player" is either player, so this watches every evolution rather than
	// only its controller's.
	//
	// Copper Locust being the base of that evolution is covered by the same
	// battle zone guard: it is already underneath the evolution by the time this
	// runs, so it has nothing to destroy.
	c.Use(fx.Creature, fx.When(fx.CreatureEvolved, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
	}))

}

// WingeyeMoth ...
func WingeyeMoth(c *match.Card) {

	c.Name = "Wingeye Moth"
	c.Power = 3000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.MyDrawStep, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		opponents := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

		strongestOpponent := 0
		for _, creature := range opponents {
			if power := ctx.Match.GetPower(creature, false); power > strongestOpponent {
				strongestOpponent = power
			}
		}

		// An empty battle zone on the other side has nothing to beat, so any
		// creature of yours clears the bar.
		outmuscles := false
		fx.Find(card.Player, match.BATTLEZONE).Map(func(creature *match.Card) {
			if len(opponents) < 1 || ctx.Match.GetPower(creature, false) > strongestOpponent {
				outmuscles = true
			}
		})

		if !outmuscles {
			return
		}

		fx.MayDraw1(card, ctx)
	}))

}
