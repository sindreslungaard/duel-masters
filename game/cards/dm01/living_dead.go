package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BoneAssassin ...
func BoneAssassin(c *match.Card) {

	c.Name = "Bone Assassin, the Ripper"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}

// BoneSpider ...
func BoneSpider(c *match.Card) {

	c.Name = "Bone Spider"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Suicide)

}

// SkeletonSoldierTheDefiled ...
func SkeletonSoldierTheDefiled(c *match.Card) {

	c.Name = "Skeleton Soldier, the Defiled"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)

}

// WanderingBraineater ...
func WanderingBraineater(c *match.Card) {

	c.Name = "Wandering Braineater"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures)

}

// WrithingBoneGhoul ...
func WrithingBoneGhoul(c *match.Card) {

	c.Name = "Writhing Bone Ghoul"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)
}
