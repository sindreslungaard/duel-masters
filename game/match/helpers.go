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

// Search prompts the user to select n cards from the specified container
func Search(p *Player, m *Match, containerName string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	cards, err := p.Container(containerName)

	if err != nil || len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, true)

	defer m.CloseAction(p)

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

			selectedCard, err := p.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)
		}

		break

	}

	return result

}

// SearchForCnd prompts the user to select n cards from the specified container that matches the given condition
func SearchForCnd(p *Player, m *Match, containerName string, condition string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	container, err := p.Container(containerName)

	if err != nil || len(container) < 1 {
		return result
	}

	cards := make([]*Card, 0)

	for _, c := range container {
		if c.HasCondition(condition) {
			cards = append(cards, c)
		}
	}

	if len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, true)

	defer m.CloseAction(p)

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

			selectedCard, err := p.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)
		}

		break

	}

	return result

}

// SearchForFamily prompts the user to select n cards from the specified container that matches the given family
func SearchForFamily(p *Player, m *Match, containerName string, family string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	container, err := p.Container(containerName)

	if err != nil || len(container) < 1 {
		return result
	}

	cards := make([]*Card, 0)

	for _, c := range container {
		if c.Family == family {
			cards = append(cards, c)
		}
	}

	if len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, true)

	defer m.CloseAction(p)

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

			selectedCard, err := p.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)
		}

		break

	}

	return result

}
