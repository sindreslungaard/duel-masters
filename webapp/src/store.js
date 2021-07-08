import { createStore } from "vuex";

export default createStore({
  state: {
    deck: {
      cards: [],
    },
  },
  mutations: {
    setDeck(state, deck) {
      state.deck = deck;
    },
    updateDeckName(state, name) {
      state.deck.name = name;
    },
    updateDeckVisibility(state, visibility) {
      state.deck.public = visibility;
    },
    addCardToDeck(state, card) {
      if (
        state.deck.cards.filter(x => x == card.uid).length >= 4 &&
        !permissions().includes("admin")
      ) {
        return;
      }

      state.deck.cards.push(card.uid);
    },
    removeCardFromDeck(state, card) {
      let toSlice = -1;
      for (let i = 0; i < state.deck.cards.length; i++) {
        if (state.deck.cards[i] === card.uid) {
          toSlice = i;
        }
      }
      if (toSlice < 0) {
        return;
      }

      state.deck.cards.splice(toSlice, 1);
    },
  },
  getters: {
    deck(state) {
      return state.deck;
    },
    deckUid(state) {
      return state.deck.uid;
    },
  },
  actions: {},
});
