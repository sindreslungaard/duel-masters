package match

import (
	"errors"
)

// CardConstructor modifies a card with the appropriate information and abilities
type CardConstructor func(*Card)

var ctors = make(map[string]CardConstructor)

// AddCard adds a new card constructor to ctors
func AddCard(uid string, ctor CardConstructor) {
	ctors[uid] = ctor
}

// CardCtor returns a cardconstructor from card uid, or an error if it does not exist
func CardCtor(uid string) (CardConstructor, error) {
	if ctors[uid] == nil {
		return nil, errors.New("Card ctor does not exist for uid " + uid)
	}
	return ctors[uid], nil
}
