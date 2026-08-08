package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonCusdalf ...
func TerradragonCusdalf(c *match.Card) {
	c.Name = "Terradragon Cusdalf"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(
		fx.Creature,
		fx.Doublebreaker,
		fx.PowerAttacker4000,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.UntapManaEvent); ok &&
				card.Zone == match.BATTLEZONE &&
				event.Player == card.Player {
				ctx.InterruptFlow()
			}
		},
	)
}

// TerradragonDakmaBalgarow ...
func TerradragonDakmaBalgarow(c *match.Card) {
	c.Name = "Terradragon Dakma Balgarow"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		shields := len(fx.Find(c.Player, match.SHIELDZONE)) + len(fx.Find(m.Opponent(c.Player), match.SHIELDZONE))

		return shields * 2000
	}

	c.Use(fx.Creature, fx.PowerBreakerTiers(6000, 15000))
}

// UltimateDragon ...
func UltimateDragon(c *match.Card) {
	c.Name = "Ultimate Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	otherDragons := func() int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		return len(fx.FindFilter(
			c.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.ID != c.ID && x.SharesAFamily(family.Dragons)
			},
		))
	}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return otherDragons() * 5000
	}

	// Crew breaker-Dragon. The modifier is kept in sync with the battle zone on
	// every event instead of being added for the duration of an attack, so a
	// cancelled attack can never leave a stale value behind.
	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if _, calculatingPower := ctx.Event.(*match.GetPowerEvent); calculatingPower {
			return
		}

		wanted := otherDragons()
		current, has := ultimateDragonOwnShieldBreakModifier(card)

		if has && current == wanted {
			return
		}

		if has {
			card.RemoveSpecificConditionBySource(cnd.ShieldBreakModifier, card.ID)
		}

		if wanted > 0 {
			card.AddUniqueSourceCondition(cnd.ShieldBreakModifier, wanted, card.ID)
		}
	})
}

func ultimateDragonOwnShieldBreakModifier(card *match.Card) (int, bool) {
	for _, condition := range card.Conditions() {
		if condition.ID != cnd.ShieldBreakModifier || condition.Src != card.ID {
			continue
		}

		if val, ok := condition.Val.(int); ok {
			return val, true
		}
	}

	return 0, false
}
