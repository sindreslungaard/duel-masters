package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"errors"
	"fmt"
)

func BlissTotemAvatarOfLuck(c *match.Card) {

	c.Name = "Bliss Totem, Avatar of Luck"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: Select up to 3 cards from your graveyard and put it in your manazone", card.Name),
			1,
			3,
			true,
		).Map(func(c *match.Card) {
			card.Player.MoveCard(c.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's manazone by %s", c.Name, c.Player.Username(), card.Name))
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}

func ClobberTotem(c *match.Card) {

	c.Name = "Clobber Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000, fx.Doublebreaker, fx.CantBeBlockedByPowerUpTo5000)

}

func ForbiddingTotem(c *match.Card) {

	c.Name = "Forbidding Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			attackableMysteryTotems, err := findAttackableMysteryTotems(card.Player, event.CardID, ctx)

			if err != nil || len(attackableMysteryTotems) == 0 {
				return
			}

			event.AttackableCreatures = attackableMysteryTotems
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			attackableMysteryTotems, err := findAttackableMysteryTotems(card.Player, event.CardID, ctx)

			if err != nil || len(attackableMysteryTotems) == 0 {
				return
			}

			ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "Your creature must attack a Mystery Totem if it attacks.")
			ctx.InterruptFlow()
		}
	})

}

func findAttackableMysteryTotems(player *match.Player, cardID string, ctx *match.Context) (fx.CardCollection, error) {
	opponentCreature, err := ctx.Match.Opponent(player).GetCard(cardID, match.BATTLEZONE)
	if err != nil || opponentCreature.HasCondition(cnd.CantAttackCreatures) {
		return nil, errors.New("not an opponent attacking creature")
	}
	opponentCreatureCanAU := opponentCreature.HasCondition(cnd.AttackUntapped)

	attackableMysteryTotems := fx.FindFilter(
		player,
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family.MysteryTotem) && (opponentCreatureCanAU || x.Tapped)
		},
	)

	return attackableMysteryTotems, nil
}
