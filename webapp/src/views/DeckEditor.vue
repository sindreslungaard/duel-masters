<template>
  <main class="deck-editor">
    <Header />
    <CardList :cards="cards" />
    <DeckView />
  </main>
</template>

<script>
import { call } from "../remote";
import Header from "../components/Header.vue";
import CardList from "@/components/cards/CardList";
import DeckView from "@/components/decks/DeckView";
import BaseMixin from "@/mixins/BaseMixin";
import _ from "lodash";

const permissions = () => {
  let p = localStorage.getItem("permissions");
  if (!p) {
    return [];
  }
  return p;
};

export default {
  name: "DeckEditor",
  components: {
    Header,
    CardList,
    DeckView
  },
  mixins: [BaseMixin],
  data() {
    return {
      selectedDeck: null,
      selectedDeckUid: null,
      hasFinishedLoading: false,
      showWizard: false,
      decks: [],
      cards: [],
      deckCopy: null,
      previewCard: null
    };
  },
  async created() {
    try {
      let cardsResponse = await call({ path: "/cards", method: "GET" });
      let decksResponse = await call({ path: "/decks", method: "GET" });
      this.decks = decksResponse.data;
      this.cards = cardsResponse.data;

      if (this.decks.length < 1) {
        this.decks.push({
          name: "My first deck",
          cards: [],
          public: false
        });
      }
      this.selectedDeck = this.decks[0];
      this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
      this.selectedDeckUid = this.selectedDeck.uid;
      this.hasFinishedLoading = true;
    } catch (e) {
      console.log(e);
    }
  },
};
</script>

<style lang="scss" scoped>
main {
  display: grid;
  grid-gap: var(--spacing);
  padding: var(--spacing);
  grid-template-columns: minmax(auto, 50%) minmax(auto, 50%);
  grid-template-rows: auto minmax(0, 1fr);
  width: 100%;
  height: 100vh;

  @include tablet {
    height: auto;
    grid-template-columns: minmax(auto, 25%) minmax(auto, 75%);
    grid-template-rows: repeat(3, auto);
  }

  @include mobile {
    grid-template-rows: repeat(4, auto);
  }
}

.header {
  grid-column: 1 / 3;
  grid-row: 1;

  @include tablet {
    grid-column: 1 / 3;
  }
}

.card-list {
  grid-column: 1;
  grid-row: 2;
}

.deck-view {
  grid-column: 2;
  grid-row: 2;
}

.actions {
  display: flex;
  margin: 0 calc(-0.5 * var(--spacing));
  margin-bottom: var(--spacing);

  .button {
    margin: 0 calc(0.5 * var(--spacing));

    &.right {
      margin-left: auto;
    }
  }
}
</style>
