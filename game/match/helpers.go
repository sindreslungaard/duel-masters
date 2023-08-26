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
func Search(p *Player, m *Match, containerOwner *Player, containerName string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	cards, err := containerOwner.Container(containerName)

	if err != nil || len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, cancellable)

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

			selectedCard, err := containerOwner.GetCard(c, containerName)

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
func SearchForCnd(p *Player, m *Match, containerOwner *Player, containerName string, condition string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	container, err := containerOwner.Container(containerName)

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

	m.NewAction(p, cards, min, max, text, cancellable)

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

			selectedCard, err := containerOwner.GetCard(c, containerName)

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
func SearchForFamily(p *Player, m *Match, containerOwner *Player, containerName string, family string, text string, min int, max int, cancellable bool) []*Card {

	result := make([]*Card, 0)

	container, err := containerOwner.Container(containerName)

	if err != nil || len(container) < 1 {
		return result
	}

	cards := make([]*Card, 0)

	for _, c := range container {
		if c.HasFamily(family) {
			cards = append(cards, c)
		}
	}

	if len(cards) < 1 {
		return result
	}

	m.NewAction(p, cards, min, max, text, cancellable)

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

			selectedCard, err := containerOwner.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)
		}

		break

	}

	return result

}

// Filter prompts the user to select n cards from the specified container that matches the given filter
func Filter(p *Player, m *Match, containerOwner *Player, containerName string, text string, min int, max int, cancellable bool, filter func(*Card) bool) []*Card {

	result := make([]*Card, 0)

	cards, err := containerOwner.Container(containerName)

	if err != nil || len(cards) < 1 {
		return result
	}

	filtered := make([]*Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			filtered = append(filtered, mCard)
		}
	}

	if len(filtered) < 1 {
		return result
	}

	m.NewAction(p, filtered, min, max, text, cancellable)

	defer m.CloseAction(p)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max || !AssertCardsIn(filtered, action.Cards...) {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, c := range action.Cards {

			selectedCard, err := containerOwner.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)

		}

		break

	}

	return result

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

// AmISummoned returns true or falsed based on if the card is played
func AmISummoned(card *Card, ctx *Context) bool {

	if event, ok := ctx.Event.(*CardMoved); ok {

		if event.CardID == card.ID && event.To == BATTLEZONE {
			return true
		}

	}

	return false
}

// AmICasted returns true or false based on if the card is casted as a spell
func AmICasted(card *Card, ctx *Context) bool {

	if event, ok := ctx.Event.(*SpellCast); ok {

		if event.CardID == card.ID {
			return true
		}

	}

	return false

}

// AmIDestroyed returns true or false based on if the card is destroyed
func AmIDestroyed(card *Card, ctx *Context) bool {

	if event, ok := ctx.Event.(*CreatureDestroyed); ok {

		if event.Card.ID == card.ID {
			return true
		}

	}

	return false
}
