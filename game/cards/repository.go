package cards

import (
	"duel-masters/game/cards/dm01"
	"duel-masters/game/match"
)

// Cards is a map with all the card id's in the game and corresponding CardConstructor
var Cards = map[string]match.CardConstructor{
	"57eeb3c3-2561-4841-a381-2e50d17533d1": dm01.AquaHulcus,
	"ecd1ae69-4f63-4e8d-a3f4-9a5c81f98a20": dm01.EmeraldGrass,
	"09b218fc-9c5a-48ef-9555-4908932271e9": dm01.AquaKnight,
	"4097a036-a775-4218-9a1d-f57ead85dda6": dm01.AquaSniper,
}
