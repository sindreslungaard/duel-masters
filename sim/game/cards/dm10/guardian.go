package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MessaBahnaExpanseGuardian ...
func MessaBahnaExpanseGuardian(c *match.Card) {

	c.Name = "Messa Bahna, Expanse Guardian"
	c.Power = 5000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.BlockIfAbleWhenOppAttacks)

}

// PalaOlesisMorningGuardian ...
func PalaOlesisMorningGuardian(c *match.Card) {

	c.Name = "Pala Olesis, Morning Guardian"
	c.Power = 2500
	c.Civs = []string{civ.Light}
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
					fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.ID != card.ID
						},
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})
				}
			})
		}))

}

// LukiaLexPinnacleGuardian ...
func LukiaLexPinnacleGuardian(c *match.Card) {

	c.Name = "Lukia Lex, Pinnacle Guardian"
	c.Power = 2500
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.PowerAttacker3000, fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}

// BluumErkisFlareGuardian ...
func BluumErkisFlareGuardian(c *match.Card) {

	c.Name = "Bluum Erkis, Flare Guardian"
	c.Power = 8500
	c.Civs = []string{civ.Light, civ.Water}
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light, civ.Water}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Doublebreaker, fx.When(fx.BreakShield, bluumErkisSpellSteal))

}

// bluumErkisSpellSteal replaces the default handling for every shield this
// creature breaks: the shield is revealed to this creature's controller
// instead of quietly reaching its owner's hand, and a spell with "shield
// trigger" is not offered to its owner as an optional cast — this creature's
// controller casts it for no cost, mandatorily, as though it were their own
// spell, and it then lands in its real owner's graveyard. Official rulings
// call this "Spell Steal": "you"/"your opponent" in the stolen spell's own
// text mean this creature's controller and their opponent, so the effect can
// end up hurting the very player forced to cast it.
//
// A creature with "shield trigger" is unaffected: the ability only lets its
// controller cast spells. It goes to its owner's hand after being revealed
// the same way, and its owner still gets the normal optional shield-trigger
// offer to summon it for free.
func bluumErkisSpellSteal(card *match.Card, ctx *match.Context) {

	event, ok := ctx.Event.(*match.BreakShieldEvent)
	if !ok || event.Source != card {
		return
	}

	// ScheduleAfter so this sees the final set of shields after every other
	// card had a chance to protect or remove some of them from the event.
	ctx.ScheduleAfter(func() {

		shields := event.Cards
		if len(shields) < 1 {
			return
		}

		ctx.InterruptFlow()

		opponent := shields[0].Player
		var normalShieldTriggers []*match.Card

		for _, shield := range shields {

			isSpellTrigger := shield.HasCondition(cnd.ShieldTrigger) && shield.HasCondition(cnd.Spell)

			if isSpellTrigger {
				ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.BrokenShieldEvent{CardID: shield.ID, Source: card.ID}))

				ctx.Match.ShowCards(card.Player, fmt.Sprintf("%s's shield, shown to you by %s:", opponent.Username(), card.Name), []string{shield.ImageID})
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s forced %s to cast %s's shield trigger spell %s", card.Name, card.Player.Username(), opponent.Username(), shield.Name))

				fx.SpellStealCast(shield, card.Player, ctx)
				continue
			}

			moved, err := opponent.MoveCard(shield.ID, match.SHIELDZONE, match.HAND, card.ID)
			if err != nil {
				continue
			}

			ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.BrokenShieldEvent{CardID: moved.ID, Source: card.ID}))

			ctx.Match.ShowCards(card.Player, fmt.Sprintf("%s's shield, shown to you by %s:", opponent.Username(), card.Name), []string{moved.ImageID})

			if moved.HasCondition(cnd.ShieldTrigger) {
				normalShieldTriggers = append(normalShieldTriggers, moved)
			}
		}

		if len(normalShieldTriggers) > 0 {
			ctx.Match.ResolveShieldTriggers(normalShieldTriggers, card)
		}
	})

}
