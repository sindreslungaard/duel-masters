package match

// AssertCardsIn returns true or false based on if the specified card ids are present in the source []*Card
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

// ContainerHas returns true or false based on if the specified container includes a card that matches the given filter
func ContainerHas(p *Player, containerName string, filter func(*Card) bool) bool {

	cards, err := p.Container(containerName)

	if err != nil {
		return false
	}

	for _, card := range cards {

		if filter(card) {
			return true
		}

	}

	return false

}
