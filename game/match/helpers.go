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

// SearchDeckForCnd prompts the user to select n cards from their deck that matches the given condition
func SearchDeckForCnd(p *Player, m *Match, condition string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	deck, err := p.Container(DECK)

	if err != nil || len(deck) < 1 {
		return result
	}

	cards := make([]*Card, 0)

	for _, c := range deck {
		if c.HasCondition(condition) {
			cards = append(cards, c)
		}
	}

	if len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, true)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max || !AssertCardsIn(cards, action.Cards...) {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, c := range action.Cards {

			selectedCard, err := p.GetCard(c, DECK)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)
		}

		break

	}

	return result

}

// SearchDeckForFamily prompts the user to select n cards from their deck that matches the given family
func SearchDeckForFamily(p *Player, m *Match, family string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	deck, err := p.Container(DECK)

	if err != nil || len(deck) < 1 {
		return result
	}

	cards := make([]*Card, 0)

	for _, c := range deck {
		if c.Family == family {
			cards = append(cards, c)
		}
	}

	if len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, true)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max || !AssertCardsIn(cards, action.Cards...) {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, c := range cards {
			result = append(result, c)
		}

		break

	}

	return result

}
