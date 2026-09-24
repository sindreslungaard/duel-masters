package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DeathPhoenixAvatarOfDoom ...
func DeathPhoenixAvatarOfDoom(c *match.Card) {

	c.Name = "Death Phoenix, Avatar of Doom"
	c.Power = 9000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Phoenix}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.Doublebreaker,
		fx.VortexEvolution(
			family.ZombieDragon, func(x *match.Card) bool { return x.HasFamily(family.ZombieDragon) },
			family.FireBird, func(x *match.Card) bool { return x.HasFamily(family.FireBird) },
		),
		fx.When(fx.BreakShield, fx.ShieldsToGraveyardInsteadOfHand),
		fx.When(fx.LeftBattlezone, fx.OpponentDiscardsHand),
	)

}

// soulPhoenixAvatarOfUnitySeparateBases implements "When this card would
// leave the battle zone, only the top card leaves the battle zone instead
// (separate the other cards into 2 creatures)."
//
// Registered before VortexEvolution in SoulPhoenixAvatarOfUnity's c.Use, so
// this runs first for the same CardMoved event: by clearing the attachments
// itself, the shared cascade that VortexEvolution runs afterward (which would
// otherwise carry both bases along to wherever this card is headed) finds
// nothing left to move.
func soulPhoenixAvatarOfUnitySeparateBases(card *match.Card, ctx *match.Context) {

	for _, base := range card.Attachments() {
		moved, err := card.Player.MoveCard(base.ID, match.HIDDENZONE, match.BATTLEZONE, card.ID)

		if err == nil && moved.Zone == match.BATTLEZONE {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s separated from %s and returned to the battle zone", moved.Name, card.Name))
		}
	}

	card.ClearAttachments()

}

// SoulPhoenixAvatarOfUnity ...
func SoulPhoenixAvatarOfUnity(c *match.Card) {

	c.Name = "Soul Phoenix, Avatar of Unity"
	c.Power = 13000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.Phoenix}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.Triplebreaker,
		fx.When(fx.LeftBattlezone, soulPhoenixAvatarOfUnitySeparateBases),
		fx.VortexEvolution(
			family.FireBird, func(x *match.Card) bool { return x.HasFamily(family.FireBird) },
			family.EarthDragon, func(x *match.Card) bool { return x.HasFamily(family.EarthDragon) },
		),
	)

}
