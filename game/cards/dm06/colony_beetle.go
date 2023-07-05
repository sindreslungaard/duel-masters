package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
    "duel-masters/game/cnd"
)

func CarrierShell(c *match.Card) {
     
    c.Name = "Carrier Shell"
    c.Power = 2000
    c.Civ = civ.Nature
    c.Family = family.ColonyBeetle
    c.ManaCost = 3
    c.ManaRequirement = []string{civ.Nature}
    
    c.Use(fx.Creature, fx.PowerAttacker3000)
}


func SlumberShell(c *match.Card) {
    
    
    c.Name = "Slumber Shell"
    c.Power = 2000
    c.Civ = civ.Nature
    c.Family = family.ColonyBeetle
    c.ManaCost = 2
    c.ManaRequirement = []string{civ.Nature}
    
    c.Use(fx.Creature)
}