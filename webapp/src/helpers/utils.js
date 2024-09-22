export const compareCards = (card1, card2, sort) => {
  var cat1 = card1[sort.by], 
      cat2 = card2[sort.by];
  if (Array.isArray(cat1)) cat1 = cat1[0];
  if (Array.isArray(cat2)) cat2 = cat2[0];
  if (cat1 == null) cat1 = "";
  if (cat2 == null) cat2 = "";

  return cat1 === parseInt(cat1, 10) &&
         cat2 === parseInt(cat2, 10)
          ? sort.directionNum *
            (cat1 < cat2
              ? -1
              : cat1 > cat2
              ? 1
              : 0)
          : sort.directionNum *
            cat1.localeCompare(cat2);
}

export const getCardsForDeck = (cardUids, cardList) => {
  let cards = [];
  for (let uid of cardUids) {
    let card = cardList.find(x => x.uid === uid)
    if (card === undefined) return [];
    card = JSON.parse(JSON.stringify(card));

    let existingCard = cards.find(x => x.uid === card.uid);
    if (existingCard) {
      existingCard.count += 1;
    } else {
      card.count = 1;
      cards.push(card);
    }
  }

  cards.sort((c1, c2) => compareCards(c1, c2, {
    by: "manaCost",
    directionNum: 1
  }));

  return cards;
}

export const playSound = (sound) => {
  if(sound) {
    var audio = new Audio(sound);
    audio.volume = 0.2;
    audio.play();
  }
}