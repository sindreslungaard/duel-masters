package match

// AssertCardsIn returns true or false based on if the specified card ids are present in the source []Card
func AssertCardsIn(src []*Card, test ...string) bool {

	for _, toTest := range test {

		ok := false

		for _, card := range src {
			if card.ID == toTest {
				ok = true
			}
		}

		if !ok {
			return false
		}

	}

	return true

}
