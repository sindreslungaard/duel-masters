package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func ProtectiveForce(c *match.Card) {

	c.Name = "Protective Force"
	c.Civ = civ.Light
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Select 1 creature with \"Blocker\" from your battlezone that will gain +4000 power",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given +4000 power by %s's effect", x.Name, card.Name))
		})

	}))

}

func InvincibleAura(c *match.Card) {

	c.Name = "Invincible Aura"
	c.Civ = civ.Light
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		topCard := card.Player.PeekDeck(3)

		for _, c := range topCard {
			card.Player.MoveCard(c.ID, match.DECK, match.SHIELDZONE, card.ID)
		}
	}))
}

func InvincibleTechnology(c *match.Card) {

	c.Name = "Invincible Technology"
	c.Civ = civ.Water
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SearchDeckTakeCards(
			card,
			ctx,
			100,
			func(x *match.Card) bool { return true },
			"cards",
		)

	}))
}

func InvincibleAbyss(c *match.Card) {

	c.Name = "Invincible Abyss"
	c.Civ = civ.Darkness
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range creatures {
			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
		}
	}))
}

func InvincibleCataclysm(c *match.Card) {

	c.Name = "Invincible Cataclysm"
	c.Civ = civ.Fire
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			fmt.Sprintf("%s: Select up to 3 of your opponent's shields and send them to graveyard", card.Name),
			1,
			3,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})
	}))

}

func InvincibleUnity(c *match.Card) {

	c.Name = "Invincible Unity"
	c.Civ = civ.Nature
	c.ManaCost = 13
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range creatures {

			creature.AddCondition(cnd.PowerAmplifier, 8000, card.ID)
			creature.AddCondition(cnd.TripleBreaker, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given +8000 power and triple breaker until the end of the turn", creature.Name))

		}
	}))
}

func SphereOfWonder(c *match.Card) {

	c.Name = "Sphere of Wonder"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		if len(fx.Find(card.Player, match.SHIELDZONE)) < len(fx.Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)) {
			cards := card.Player.PeekDeck(1)

			for _, toMove := range cards {
				card.Player.MoveCard(toMove.ID, match.DECK, match.SHIELDZONE, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put a shield into the shieldzone from the top of their deck by %s's ability", card.Player.Username(), card.Name))
			}
		}
	}))

}

func MysticDreamscape(c *match.Card) {

	c.Name = "Mystic Dreamscape"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Return up to 3 cards from your mana zone to your hand", card.Name),
			1,
			3,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = false
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their manazone by %s", x.Name, card.Player.Username(), card.Name))
		})
	}))

}

func FutureSlash(c *match.Card) {

	c.Name = "Future Slash"
	c.Civ = civ.Darkness
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		opponent := ctx.Match.Opponent(card.Player)

		fx.Select(
			card.Player,
			ctx.Match,
			opponent,
			match.DECK,
			fmt.Sprintf("%s: Send up to 2 cards from your opponent's deck to his graveyard", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s in graveyard from their opponent's deck", x.Player.Username(), x.Name))
		})

		fx.ShuffleDeck(card, ctx, true)

	}))

}

func ProclamationOfDeath(c *match.Card) {

	c.Name = "Proclamation of Death"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 creature from your battlezone that will be sent to your graveyard", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}

func PhantomDragonsFlame(c *match.Card) {

	c.Name = "Phantom Dragon's Flame"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.DestroyBySpellOpCreature2000OrLess))
}

func SpasticMissile(c *match.Card) {

	c.Name = "Spastic Missile"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's creatures that will be destroyed", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 3000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})
	}))
}

func MysticTreasureChest(c *match.Card) {

	c.Name = "Mystic Treasure Chest"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SearchDeckPutIntoManazone(
			card,
			ctx,
			1,
			func(x *match.Card) bool { return x.Civ != civ.Nature },
			"non-nature card",
		)

	}))
}

func PangaeasWill(c *match.Card) {

	c.Name = "Pangaea's Will"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: You may select an evolution card from your opponent's battle zone and send the top card to their mana zone", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Evolution) },
			false,
		).Map(func(x *match.Card) {
			tapped := x.Tapped

			for _, baseCard := range x.Attachments() {
				baseCard.Player.MoveCard(baseCard.ID, match.HIDDENZONE, match.BATTLEZONE, card.ID)

				if tapped {
					baseCard.Tapped = true
				}
			}

			ctx.Match.BroadcastState()
			x.ClearAttachments()

			ctx.Match.MoveCard(x, match.MANAZONE, card)
		})
	}))

}

func FaerieLife(c *match.Card) {

	c.Name = "Faerie Life"
	c.Civ = civ.Nature
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.Draw1ToMana))
}

func BondsOfJustice(c *match.Card) {

	c.Name = "Bonds of Justice"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		notblockers := make([]*match.Card, 0)

		opponentBattlezone, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		myBattlezone, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		for _, creature := range myBattlezone {
			if !creature.HasCondition(cnd.Blocker) {
				notblockers = append(notblockers, creature)
			}
		}

		for _, creature := range opponentBattlezone {
			if !creature.HasCondition(cnd.Blocker) {
				notblockers = append(notblockers, creature)
			}
		}
		for _, notblocker := range notblockers {
			notblocker.Tapped = true
			ctx.Match.ReportActionInChat(notblocker.Player, fmt.Sprintf("%s was tapped by %s", notblocker.Name, card.Name))
		}

	}))
}

func EnergyStream(c *match.Card) {

	c.Name = "Energy Stream"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		card.Player.DrawCards(2)
	}))

}

func RainOfArrows(c *match.Card) {

	c.Name = "Rain of Arrows"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		opponentHand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil {
			return
		}

		ids := make([]string, 0)

		for _, c := range opponentHand {
			ids = append(ids, c.ImageID)
		}

		ctx.Match.ShowCards(
			card.Player,
			"Your opponent's Hand:",
			ids,
		)

		for _, c := range opponentHand {
			if c.Civ == civ.Darkness && c.HasCondition(cnd.Spell) {
				c.Player.MoveCard(c.ID, match.HAND, match.GRAVEYARD, card.ID)
				ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's graveyard by %s", c.Name, c.Player.Username(), card.Name))
			}
		}

	}))

}

// CometMissile ...
func CometMissile(c *match.Card) {

	c.Name = "Comet Missile"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Destroy one of your opponent's creatures that has blocker and has power 6000 or less",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) && ctx.Match.GetPower(x, false) <= 6000 },
			false,
		)

		for _, creature := range creatures {

			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)

		}

	}))

}

// CrisisBoulder ...
func CrisisBoulder(c *match.Card) {

	c.Name = "Crisis Boulder"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		cards["Your creatures"] = fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)
		cards["Your mana"] = fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)

		fx.SelectMultipart(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose 1 of your creatures or a card in your mana zone it will be sent to your graveyard.", card.Name),
			1,
			1,
			false,
		).Map(func(selectedCard *match.Card) {
			if selectedCard.Zone == match.BATTLEZONE {
				ctx.Match.Destroy(selectedCard, card, match.DestroyedBySpell)
				return
			}
			selectedCard.Player.MoveCard(selectedCard.ID, selectedCard.Zone, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's graveyard by %s", selectedCard.Name, selectedCard.Player.Username(), card.Name))
		})

	}))
}

// IntenseEvil ...
func IntenseEvil(c *match.Card) {

	c.Name = "Intense Evil"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		battleZone, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		nrCreature := len(battleZone)

		selected := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Intense Evil: Destroy any number of your creatures.",
			0,
			nrCreature,
			false,
		)

		nrSelected := len(selected)

		selected.Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})

		card.Player.DrawCards(nrSelected)

	}))
}

// ShockHurricane ...
func ShockHurricane(c *match.Card) {

	c.Name = "Shock Hurricane"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		battleZone, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		nrCreature := len(battleZone)

		selected := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Shock Hurricane: Return any number of your creatures from the battlezone to your hand.",
			0,
			nrCreature,
			false,
		)

		nrSelected := len(selected)

		selected.Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved from %s's battlezone to his hand.", x.Name, x.Player.Username()))
		})

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Choose %d creature(s) in the battlezone that will be sent to your opponent's hand", card.Name, nrSelected),
			0,
			nrSelected,
			false,
		).Map(func(creature *match.Card) {
			creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
		})

	}))

}
