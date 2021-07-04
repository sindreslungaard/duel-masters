<template>
  <Panel title="Deck" class="deck-view">
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <div class="view-mode-switcher">
        <div class="field">List
          <input v-model="viewMode" value="list" type="radio" />
        </div>
      </div>
      <div class="field deck-selection">
        <label for="decks">Decks</label>
        <select
          id="decks"
          v-model="selectedDeckUid"
        >
          <option v-for="deck in decks" :key="deck.uid" :value="deck.uid">{{
            deck.name
          }}</option>
        </select>
      </div>
      <div class="field" v-if="$store.state.deck.public">
        <label for="share-url">Share URL</label>
        <input disabled id="share-url" :value="getShareUrl($store.state.deck.uid)" />
      </div>
      <div class="actions">
        <div class="actions--left"

          >
          <template v-if="
              $store.getters.deck && deckCopy && !decksEqual($store.getters.deck, deckCopy)
            ">
            <Button type="success" @click="saveDeck" text="Save changes"></Button>
            <Button
              type="error"
              @click="discardDeck"
              text="Discard changes"
            ></Button>
            </template>
        </div>
        <div class="actions--right">
          <Button type="success" @click="newDeck" text="New Deck"></Button>
          <Button text="Edit Deck" @click="openEditDeckDialog"></Button>
          <Button type="error"  @click="deleteDeck" text="Delete Deck"></Button>
        </div>

      </div>

      <template v-if="$store.getters.deck">
          {{ $store.getters.deck.cards.length }} cards

          <table class="list-view" v-if="viewMode === 'list'">
        <thead>
          <tr>
            <th>Count</th>
            <th>Name</th>
            <th>Set</th>
            <th>Civilization</th>
          </tr>
        </thead>
        <tbody>
          <tr
                @click.right.prevent="removeCardFromDeck(card)"
                @click.middle.prevent="addCardToDeck(card)"
                v-for="(card, index) in getCardsForDeck($store.getters.deck.cards)"
                :key="index"
          >
            <td>{{ card.count }}</td>
            <td class="fill">{{ card.name }}</td>
            <td class="set no-break">{{ card.set }}</td>
            <td class="civilization no-break">{{ card.civilization }}</td>
          </tr>
        </tbody>
      </table>

        </template>
    </template>
  </Panel>
</template>

<script>
import Panel from "@/components/Panel";
import LoadingIndicator from "@/components/LoadingIndicator";
import Button from "@/components/buttons/Button";
import DeckEditDialog from "@/components/dialogs/DeckEditDialog";
import { call } from "@/remote";
import _ from "lodash";
import BaseMixin from "@/mixins/BaseMixin";

export default {
  name: "DeckView",
  components: {
    Panel,
    LoadingIndicator,
    Button
  },
  mixins: [BaseMixin],
  props: {
    isLoading: {
      type: Boolean,
      required: false,
      default: false
    },
  },
  data() {
    return {
      cards: [],
      decks: [],
      selectedDeckUid: null,
      deckCopy: null,
      viewMode: "list",
    };
  },
  async created() {
    await this.fetchDecks();
    let cardsResponse = await call({ path: "/cards", method: "GET" });
    this.cards = cardsResponse.data;

    if(this.decks.length === 0) {
      this.newDeck();
    }

    this.$store.commit("setDeck", this.decks[0]);
    this.deckCopy = JSON.parse(JSON.stringify(this.$store.getters.deck));
    this.selectedDeckUid = this.$store.getters.deckUid;
  },
  methods: {
    async fetchDecks() {
      let decksResponse = await call({ path: "/decks", method: "GET" });
      this.decks = decksResponse.data;
    },
    async saveDeck() {
      try {
        let res = await call({
          path: "/decks",
          method: "POST",
          body: this.$store.getters.deck
        });
        this.deckCopy = JSON.parse(JSON.stringify(this.$store.getters.deck));
        this.showWarning("Successfully saved your deck");
      } catch (e) {
        this.showWarning(
          "Invalid request. Please ensure that the deck name is 1-30 characters and that you have between 40-50 cards in your deck."
        );
      }
    },
    discardDeck() {
      if (this.deckCopy.name === "to.be.removed") {
        this.$store.commit("setDeck", this.decks[0]);
        this.deckCopy = JSON.parse(JSON.stringify(this.$store.getters.deck));
        this.selectedDeckUid = this.$store.getters.deckUid;
        this.decks.pop();

        if(this.decks.length === 0) {
          this.newDeck();
        }

        return;
      }

      this.$store.commit("setDeck", JSON.parse(JSON.stringify(this.deckCopy)));
    },
    async deleteDeck() {
      try {
        let res = await call({
          path: "/deck/" + this.selectedDeckUid,
          method: "DELETE",
        });

        this.fetchDecks();

        // TODO: Check confirmation dialog
        if(this.decks.length > 0) {
          //this.selectedDeckUid = this.decks[0].uid;
          //this.selectDeck(this.decks[0]);
        } else {
          this.newDeck();
        }

        this.showWarning("Successfully deleted your deck");
      } catch (e) {
        this.showWarning("Couldn't delete the deck you selected");
      }
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
    removeCardFromDeck(card) {
      this.$store.commit("removeCardFromDeck", card);
    },
    addCardToDeck(card) {
      this.$store.commit("addCardToDeck", card);
    },
    openEditDeckDialog() {
      this.$modal.show(DeckEditDialog);
    },
    newDeck() {
      if (this.decks.length > 0 && !this.decksEqual(this.$store.state.deck, this.deckCopy)) {
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

      this.$store.commit("setDeck", this.decks[this.decks.length - 1]);
      this.deckCopy = JSON.parse(JSON.stringify(this.$store.state.deck));


      this.selectedDeckUid = this.$store.state.deck.uid;
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
    },
    getShareUrl(deckUid) {
      return window.location.origin + "/deck/" + deckUid;
    },
  },
  watch: {
      selectedDeckUid: function(value) {
      if (!this.decksEqual(this.$store.getters.deck, this.deckCopy)) {
        if (value != this.$store.getters.deckUid) {
          this.showWarning(
            "You have unsaved changes in the currently selected deck. Save or discard before editing another deck."
          );
        }

        this.selectedDeckUid = this.$store.getters.deckUid;
        return;
      }
      this.$store.commit("setDeck", this.decks.find(x => x.uid === value));
      console.log(this.$store.getters.deck);
      this.deckCopy = JSON.parse(JSON.stringify(this.$store.getters.deck));

    },
  },
  computed: {
    /**
     * Whether the component is ready to be displayed.
     */
    hasFinishedLoading() {
      return !this.isLoading && this.decks.length > 0;
    }
  }
};
</script>
<style lang="scss" scoped>
table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid var(--color-foreground-border);
  padding: calc(0.5 * var(--spacing));
  text-align: left;
}

.view-mode-switcher {
  margin-bottom: var(--spacing);
}

.actions {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing);

  &--left,
  &--right {
    display: flex;
    gap: var(--spacing);
  }
}
</style>
