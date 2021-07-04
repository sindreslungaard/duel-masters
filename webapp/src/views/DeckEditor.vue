<template>
  <main class="deck-editor">
    <!--
    <div v-if="previewCard" class="card-preview">
      <img :src="`/assets/cards/all/${previewCard.uid}.jpg`" />
      <div @click="previewCard = null" class="btn">Close</div>
    </div>


    <div v-if="showWizard" class="new-duel">
      <div class="backdrop"></div>
      <div class="wizard">
        <div class="spacer">
          <span class="headline">Edit your deck</span>
          <br /><br />
          <form>
            <span class="helper">Name</span>
            <input v-model="selectedDeck.name" type="text" placeholder="Name" />
            <br /><br />
            <span class="helper">Visibility</span>
            <select v-model="selectedDeck.public">
              <option :value="true">Public</option>
              <option :value="false">Private</option>
            </select>

            <div @click="showWizard = false" class="btn">
              Done
            </div>
          </form>
        </div>
      </div>
    </div>-->

    <Header />
    <CardList @leftClick="addCardToDeck" :cards="cards" />
    <Panel title="Deck" class="deck-view">
      <LoadingIndicator v-if="!hasFinishedLoading" />

      <template v-if="hasFinishedLoading">
        <div class="field deck-selection">
          <label for="decks">Decks</label>
          <select id="decks" v-model="selectedDeckUid">
            <option v-for="deck in decks" :key="deck.uid" :value="deck.uid">{{
              deck.name
            }}</option>
          </select>
        </div>

        <div class="actions">
          <template
            v-if="
              selectedDeck && deckCopy && !decksEqual(selectedDeck, deckCopy)
            "
          >
            <Button type="success" @click="save" text="Save changes"></Button>
            <Button
              type="error"
              @click="discard"
              text="Discard changes"
            ></Button>
          </template>
          <Button class="right" type="success" text="New Deck"></Button>
          <Button text="Edit Deck"></Button>
        </div>
        {{ selectedDeck.cards.length }} cards
        <table>
          <tr>
            <th>Quantity</th>
            <th>Card Name</th>
            <th>Set</th>
            <th>Civilization</th>
          </tr>
          <template v-if="selectedDeck">
            <tr
              @click.right.prevent="removeCardFromDeck(card)"
              @click.middle.prevent="addCardToDeck(card)"
              v-for="(card, index) in getCardsForDeck(selectedDeck.cards)"
              :key="index"
            >
              <td>{{ card.count }}</td>
              <td>{{ card.name }}</td>
              <td class="set">{{ card.set }}</td>
              <td class="civilization">{{ card.civilization }}</td>
            </tr>
          </template>
        </table>
      </template>
    </Panel>
  </main>
</template>

<script>
import { call } from "../remote";
import Header from "../components/Header.vue";
import CardList from "@/components/cards/CardList";
import Panel from "@/components/Panel";
import LoadingIndicator from "@/components/LoadingIndicator";
import Button from "@/components/buttons/Button";
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
    LoadingIndicator,
    Panel,
    Button
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
  methods: {
    addCardToDeck(card) {
      if (
        this.selectedDeck.cards.filter(x => x == card.uid).length >= 4 &&
        !permissions().includes("admin")
      ) {
        return;
      }

      this.selectedDeck.cards.push(card.uid);
    },
    removeCardFromDeck(card) {
      let toSlice = -1;
      for (let i = 0; i < this.selectedDeck.cards.length; i++) {
        if (this.selectedDeck.cards[i] === card.uid) {
          toSlice = i;
        }
      }
      if (toSlice < 0) {
        return;
      }

      this.selectedDeck.cards.splice(toSlice, 1);
    },
    async save() {
      try {
        let res = await call({
          path: "/decks",
          method: "POST",
          body: this.selectedDeck
        });
        this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
        this.showWarning("Successfully saved your deck");
      } catch (e) {
        this.showWarning(
          "Invalid request. Please ensure that the deck name is 1-30 characters and that you have between 40-50 cards in your deck."
        );
      }
    },

    discard() {
      if (this.deckCopy.name === "to.be.removed") {
        this.selectedDeck = this.decks[0];
        this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
        this.selectedDeckUid = this.selectedDeck.uid;
        this.decks.pop();
        return;
      }

      this.selectedDeck = JSON.parse(JSON.stringify(this.deckCopy));
    },

    getCardsForDeck(cardUids) {
      let cards = [];
      for (let uid of cardUids) {
        let card = JSON.parse(
          JSON.stringify(this.cards.find(x => x.uid === uid))
        );

        let existingCard = cards.find(x => x.uid === card.uid);
        if (existingCard) {
          existingCard.count += 1;
        } else {
          card.count = 1;
          cards.push(card);
        }
      }

      cards = _.sortBy(cards, "name");
      return cards;
    },

    getShareUrl(uid) {
      return window.location.host + "/decks/" + uid;
    },

    newDeck() {
      if (!this.decksEqual(this.selectedDeck, this.deckCopy)) {
        this.showWarning(
          "Please save or discard the changes you've made before creating a new deck"
        );
        return;
      }
      this.decks.push({
        name: "Unnamed Deck",
        cards: [],
        public: false
      });
      this.deckCopy = JSON.parse(
        JSON.stringify(this.decks[this.decks.length - 1])
      );
      this.selectedDeck = this.decks[this.decks.length - 1];
      this.selectedDeckUid = this.selectedDeck.uid;
      this.$nextTick(() => {
        this.deckCopy.name = "to.be.removed";
      });
    },

    decksEqual(deck1, deck2) {
      if (deck1.name !== deck2.name) {
        return false;
      }
      if (deck1.public !== deck2.public) {
        return false;
      }
      if (deck1.cards.length !== deck2.cards.length) {
        return false;
      }
      for (let i = 0; i < deck1.cards.length; i++) {
        if (deck1.cards[i] !== deck2.cards[i]) {
          return false;
        }
      }
      return true;
    }
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
  watch: {
    selectedDeckUid: function(val) {
      if (!this.decksEqual(this.selectedDeck, this.deckCopy)) {
        if (val != this.selectedDeck.uid) {
          this.showWarning(
            "You have unsaved changes in the currently selected deck. Save or discard before editing another deck."
          );
        }

        this.selectedDeckUid = this.selectedDeck.uid;
        return;
      }
      this.selectedDeck = this.decks.find(x => x.uid === val);
      this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
    }
  }
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
