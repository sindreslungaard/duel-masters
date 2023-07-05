package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
    //"duel-masters/game/cnd"
)



func PyrofighterMagnus(c *match.Card) {
    
    c.Name = "Pyrofighter Magnus"
    c.Power = 3000
    c.Civ = civ.Fire
    c.Family = family.Dragonoid
    c.ManaCost = 3
    c.ManaRequirement = []string{civ.Fire}
    
    c.Use(fx.Creature, fx.SpeedAttacker, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
        
            
            ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to the %s's hand", c.Name, c.Player.Username()))
    }))
}
