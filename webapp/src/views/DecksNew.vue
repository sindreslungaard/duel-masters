<template>
  <div>
    <div v-show="warning" @click="closeOverlay()" class="overlay"></div>

    <div v-show="warning" class="error">
      <p class="text-block">{{ warning }}</p>
      <div @click="warning = ''" class="btn">Close</div>
    </div>

    <div v-if="showWizard" class="new-duel">
      <div class="backdrop" @click="showWizard = false"></div>
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
    </div>

    <Header></Header>

    <div class="main">
      <CardsCatalogue 
        class="catalogue"
        :unsortedCards="cards"
        v-model:deck="selectedDeck.cards"
      />

      <div class="right">
        <span class="deck-card-total" v-if="selectedDeck && selectedDeck.cards">({{selectedDeck.cards.length}})</span>
        <select v-model="selectedDeckUid" style="margin: 0;">
          <option
            v-for="(deck, index) in decks"
            :key="index"
            :value="deck.uid"
            >{{ deck.name }}</option
          >
        </select> 
        <div class="deck-secondary-buttons">
          <img
            @click="showWizard = true"
            class="float-left edit-ico"
            width="25px"
            src="/assets/images/edit_icon.png"
            v-tooltip="'Edit deck info'"
          />
          <a
            :href="'/deck/' + selectedDeckUid"
            v-if="selectedDeck && selectedDeck.public"
            target="_blank"
            v-tooltip="'Share deck'"
          >
            <img
              class="float-left edit-ico share"
              width="25px"
              src="/assets/images/share_icon.png"
            />
          </a>

          <a
            v-if="selectedDeck && selectedDeck.uid"
            @click="deleteDeck(selectedDeckUid)"
            target="_blank"
            v-tooltip="'Delete deck'"
          >
            <img
              class="float-left edit-ico share"
              width="25px"
              src="/assets/images/delete_icon.png"
            />
          </a>

          <a
            v-if="selectedDeck && selectedDeck.uid"
            @click="copyDeckList()"
            target="_blank"
            v-tooltip="'Copy deck list'"
          >
            <img
              class="float-left edit-ico share"
              width="25px"
              src="/assets/images/list_icon.png"
            />
          </a>
        </div>
        <div class="deck-crud-buttons">
          <div @click="newDeck()" class="btn new font-bold text-xs">NEW DECK</div>
          <template
            v-if="
              selectedDeck && deckCopy && !decksEqual(selectedDeck, deckCopy)
            "
          >
            <div @click="save()" class="btn save">Save</div>
            <div @click="discard()" class="btn discard">Discard</div>
          </template>
        </div>
        <div v-if="selectedDeck" class="right-content">
          <DeckList 
            :cards="cards"
            v-model:deck="selectedDeck.cards"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";
import Header from "../components/Header.vue";
import CardsCatalogue from "../components/CardsCatalogue.vue";
import { ref, reactive } from 'vue'
import { getCardsForDeck } from "../helpers/utils";
import DeckList from "../components/DeckList.vue";

export default {
  name: "decks",
  components: {
    CardsCatalogue,
    DeckList,
    Header,
  },
  computed: {
    username: () => localStorage.getItem("username")
  },
  data() {
    return {
      warning: "",
      showWizard: false,

      cards: [],

      decks: [],
      selectedDeck: { cards: [] },
      selectedDeckUid: null,
      deckCopy: null,
    };
  },
  methods: {
    copyDeckList() {
      const cards = getCardsForDeck(this.selectedDeck.cards, this.cards);

      const deckList = cards.map(card => `${card.count}x ${card.name}`);
      this.warning = `${deckList.join("\n")}\n\nTotal Cards: ${
        this.selectedDeck.cards.length
      }`;
    },

    newDeck() {
      if (!this.decksEqual(this.selectedDeck, this.deckCopy)) {
        this.warning =
          "Please save or discard the changes you've made before creating a new deck";
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

    closeOverlay() {
      this.warning = null;
    },

    async save() {
      try {
        let res = await call({
          path: "/decks",
          method: "POST",
          body: this.selectedDeck
        });
        this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
        this.warning = "Successfully saved your deck";
      } catch (e) {
        this.warning =
          "Invalid request. Please ensure that the deck name is 1-30 characters and that you have between 40-50 cards in your deck.";
      }
    },

    async deleteDeck() {
      try {
        let res = await call({
          path: "/deck/" + this.selectedDeckUid,
          method: "DELETE"
        });

        this.decks = this.decks.filter(x => x.uid !== this.selectedDeckUid);
        if (this.decks.length > 0) {
          this.selectedDeckUid = this.decks[0].uid;
          this.selectDeck(this.decks[0]);
        } else {
          this.newDeck();
        }

        this.warning = "Successfully deleted your deck";
      } catch (e) {
        this.warning = "Couldn't delete the deck you selected";
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
  },
  async created() {
    try {
      let [cards, decks] = await Promise.all([
        call({ path: "/cards", method: "GET" }),
        call({ path: "/decks", method: "GET" })
      ]);
      this.cards = cards.data;
      this.decks = decks.data;

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
    } catch (e) {
      console.log(e);
    }
  },
  watch: {
    selectedDeckUid: function(val) {
      if (!this.decksEqual(this.selectedDeck, this.deckCopy)) {
        this.warning =
          "You have unsaved changes in the currently selected deck. Save or discard before editing another deck.";
        this.selectedDeckUid = this.selectedDeck.uid;
        return;
      }
      this.selectedDeck = this.decks.find(x => x.uid === val);
      this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck));
    }
  },
};
</script>

<style scoped lang="scss">

.helper {
  display: block !important;
  margin-bottom: 5px;
}

.wizard input {
  margin-left: 0 !important;
}

.wizard select {
  margin-left: 0 !important;
  margin-bottom: 10px;
}

.wizard .btn {
  display: block;
}


.edit-ico {
  filter: invert(70%);
  margin: 2px 10px;
}
.edit-ico:hover {
  filter: invert(80%);
  cursor: pointer;
}

.new {
  display: inline-block;
  margin-right: 15px;
}

.save {
  background: #3ca374 !important;
  display: inline-block;
  margin-right: 15px;
}

.save:hover {
  background: #3ca374 !important;
}

.copy-deck-list {
  background: #dcb875 !important;
  display: inline-block;
  margin-right: 15px;
}

.copy-deck-list:hover {
  background: #d3b06e !important;
}

.discard {
  background: #ff4c4c !important;
  display: inline-block;
}

.discard:hover {
  background: #ed3e3e !important;
}

.right-content {
  margin-top: 10px;
  overflow-y: auto;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: #000;
  opacity: 0.5;
  z-index: 1000;
}

.share {
  margin-left: 0;
}

.selected {
  background: #7289da;
  border-color: #7289da !important;
}

.right select {
  padding: 7px !important;
}

input,
textarea,
select {
  border: none;
  background: #484c52;
  padding: 5px !important;
  width: auto !important;
  margin-left: 5px;
  border-radius: 4px;
  color: #ccc;
  resize: none;
}

select option {
  background: #111;
}

.deck-crud-buttons {
  display: inline-flex;
  align-items: center;
}

.deck-secondary-buttons {
  display: inline-flex;
}

.deck-card-total {
  padding: 2px;
  background-color: lightgray;
  color: black;
}

.right {
  display: inline-block;
  height: calc(100vh - 100px);
  overflow-y: auto;
  width: calc(20% - 10px);
  max-width: 380px;
  flex-shrink: 0;
}

.catalogue {
  display: inline-flex;
  height: calc(100vh - 100px);
  min-width: calc(80% - 10px);
  margin-right: 20px;
  flex-grow: 1;
}


.main {
  margin: 0 15px;
  display: flex;
  align-items: flex-start;
}

.new-duel .backdrop {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: #000;
  opacity: 0.5;
}

.new-duel .wizard {
  position: absolute;
  top: calc(50vh - 323px / 2);
  left: calc(50% - 250px / 2);
  width: 250px;
  border-radius: 4px;
  color: #fff;
  background: #0C0C0F;
  border: 1px solid #333;
  z-index: 5;
}

.wizard .headline {
  color: #ccc;
}

.wizard .spacer {
  margin: 15px;
}

.wizard .helper {
  color: #ccc;
  font-size: 13px;
}

.wizard .btn {
  margin: 0;
  width: 85px;
  text-align: center;
  margin-top: 15px;
}

.wizard .cancel {
  margin-left: 10px;
  background: #ff4c4c;
  color: #fff;
}

.wizard .cancel:hover {
  background: #ed3e3e;
}

input,
textarea,
select {
  border: none;
  background: url(/assets/images/overlay_30.png);
  padding: 10px;
  border-radius: 4px;
  width: 200px;
  color: #ccc;
  resize: none;
}
input:focus,
textarea:focus,
select:focus {
  outline: none;
}
input:active,
textarea:active,
select:active {
  outline: none;
}

.wizard select {
  width: 220px;
  margin-top: 4px;
}

.error p {
  padding: 10px 12px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: url(/assets/images/overlay_30.png) !important;
}

.error {
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #0C0C0F;
  border: 1px solid #333;
  z-index: 1000;
  left: calc(50% - 300px / 2);
  top: 50vh;
  transform:translateY(-50%);
  padding: 10px;
  font-size: 14px;
  color: #ccc;
}

.wizard .errorMsg {
  color: red;
  font-size: 14px;
  display: block;
  margin-top: 15px;
}

nav {
  text-align: right;
}

ul {
  list-style: none;
}

li {
  display: inline-block;
  padding-right: 10px;
  margin-right: 10px;
}

nav > ul > li:hover {
  cursor: pointer;
  color: #fff;
}

nav > ul > li.no-cursor:hover {
  cursor: default;
}

.title {
  position: absolute;
  top: 16px;
  left: 16px;
}

.psa {
  margin: 16px;
  background: #2b2c31;
  padding: 5px;
  min-height: 20px;
  border-radius: 4px;
  font-size: 14px;
  color: #ccc;
}

.psa > span {
  display: inline-block;
  vertical-align: middle;
  margin-left: 4px;
}

a {
  color: #7289da;
}

.text-block {
  white-space: pre-wrap;
}

.btn {
  display: inline-block;
  background: #5865F2;
  color: #e3e3e5;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn:hover {
  cursor: pointer;
  background: #515de2;
}

.btn:active {
  background: #4c58d3 !important;
}
</style>
