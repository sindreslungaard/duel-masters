package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonArqueDelacerna ...
func TerradragonArqueDelacerna(c *match.Card) {

	c.Name = "Terradragon Arque Delacerna"
	c.Power = 6000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PutIntoBattleZoneInsteadOfDiscard)

}
