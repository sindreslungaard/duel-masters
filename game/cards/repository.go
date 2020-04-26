package cards

import (
	"duel-masters/game/cards/dm01"
	"duel-masters/game/match"
)

// CardConstructor modifies a card with the appropriate information and abilities
type CardConstructor func(*match.Card)

// Cards is a map with all the card id's in the game and corresponding CardConstructor
var Cards = map[string]CardConstructor{
	"1d72eb3e-5185-449a-a16f-391bd2338343": dm01.AquaHulcus,
}
