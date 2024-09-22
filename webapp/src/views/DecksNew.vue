<template>
  <div>
    <div v-show="warning" @click="closeOverlay()" class="overlay"></div>

    <CardPreviewPopup 
      :uid = previewCard?.uid
      :event = previewCardEvent
      :xPos = previewCardX
      :side = "'left'"
    />

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
      <div class="catalogue">
        <div class="catalogue-top-bar">
          <div class="catalogue-filters">
            <input
              v-model="filterCard"
              class="catalogue-filter"
              type="search"
              placeholder="Type to search"
            />
            <div 
              v-for="civ in ['light', 'darkness', 'nature', 'fire', 'water']" 
              class="civ-filter"
              :class="['civ-color-' + civ, filterCivilization[civ] ? 'civ-filter-selected' : '']"
              @click="filterCivilization[civ] = !filterCivilization[civ]">
            </div>
            <div style="margin-left: 5px;">
              <select 
                class="catalogue-filter"
                v-model="filterFamily">
                <option
                  class="family"
                  v-for="(family, index) in families"
                  :key="index"
                  :value="family"
                  >{{ family }}</option
                >
              </select>
            </div>
            <div 
              v-for="manaNr in ['1', '2', '3', '4', '5', '6', '7+']" 
              class="mana-filter"
              :class="{'mana-filter-selected': filterMana[manaNr]}"
              @click="filterMana[manaNr] = !filterMana[manaNr]"
            >
              {{manaNr}}
            </div>
            <select 
              v-model="filterSet"
              class="catalogue-filter"
            >
              <option
                class="set"
                v-for="(set, index) in sets"
                :key="index"
                :value="set"
                >{{ set }}</option
              >
            </select>
            <img
              class="reset-icon"
              src="/assets/images/reset-icon.svg"
              v-tooltip="'Reset all filters'"
              @click="resetFilters"
            />
          </div>
          <div class="zoom-icons">
            <img
              class="zoom-icon"
              @click="modifyCatalogueCardSize(10)"
              src="/assets/images/zoom-in-icon.svg"
              v-tooltip="'Increase card size'"
            />
            <img
              @click="modifyCatalogueCardSize(-10)"
              class="zoom-icon"
              src="/assets/images/zoom-out-icon.svg"
              v-tooltip="'Decrease card size'"
            />
          </div>
        </div>
        <div class="catalogue-cards-wrapper">
          <div class="catalogue-cards">
            <div
              v-for="card in filteredAndSortedCards"
              class="catalogue-card"
              :style="{ 'height': cardSize + 'px' }"
              @click="tryAddCard(card)"
              >
              <v-lazy-image
                class="max-w-full rounded-2xl"
                :src="`https://scans.shobu.io/${card.uid}.jpg`"
                src-placeholder="https://scans.shobu.io/backside.jpg"
                :alt="card.name"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="right" ref="rightSide">
        <span class="deck-card-total" v-if="selectedDeck">({{selectedDeck.cards.length}})</span>
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
            class="fl edit-ico"
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
              class="fl edit-ico share"
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
              class="fl edit-ico share"
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
              class="fl edit-ico share"
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
          <div
            class="deck-card-slot"
            v-for="(card, index) in getCardsForDeck(selectedDeck.cards)"
            :key="index"
            @mouseover="setPreviewCard($event, card)"
            @mouseleave="previewCard = null"
          >
            <div class="deck-card-counter">
              {{ card.count }}x
            </div>
            <div
              class="deck-card-oval"
              :class="'card-' + card.civilization.toLowerCase()"
              :key="index"
              @click="tryRemoveCard(card)"
            >
              <div class="deck-card-name">
                {{ card.name }}                  
              </div>
              <div class="deck-card-mana">
                {{ card.manaCost }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";
import Header from "../components/Header.vue";
import CardPreviewPopup from "../components/CardPreviewPopup.vue";
import VLazyImage from "v-lazy-image";

const ALL_FAMILIES = 'All Races';
const ALL_SETS = 'All Sets'

const permissions = () => {
  let p = localStorage.getItem("permissions");
  if (!p) {
    return [];
  }
  return p;
};

function compareCards(card1, card2, sort) {
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

function playSound (sound) {
  if(sound) {
    var audio = new Audio(sound);
    audio.volume = 0.2;
    audio.play();
  }
}

export default {
  name: "decks",
  components: {
    Header,
    CardPreviewPopup,
    VLazyImage,

  },
  computed: {
    username: () => localStorage.getItem("username")
  },
  data() {
    return {
      warning: "",
      showWizard: false,

      filterCard: "",
      filterFamily: ALL_FAMILIES,
      filterSet: ALL_SETS,
      families: [ALL_FAMILIES, "Spell"],
      filterCivilization: {
        "light": false,
        "darkness": false,
        "nature": false,
        "fire": false,
        "water": false,
      },
      filterMana: {
        '1': false,
        '2': false,
        '3': false,
        '4': false,
        '5': false,
        '6': false,
        '7+': false,
      },

      sets: [],
      cards: [],

      decks: [],
      selectedDeck: null,
      selectedDeckUid: null,
      deckCopy: null,

      previewCard: null,
      previewCardX: 0,
      previewCardEvent: null,

      cardSize: 360,
    };
  },
  methods: {
    setPreviewCard(event, card) {
      this.previewCardEvent = event;
      this.previewCardX = this.$refs.rightSide.getBoundingClientRect().x;
      this.previewCard = card;
    },

    selectDeck(deck) {
      this.selectedDeck = deck;
      this.deckCopy = JSON.parse(JSON.stringify(deck));
    },

    cardInfo(uid) {
      let card = this.cards.find(x => x.uid === uid);
      return card;
    },

    resetFilters() {
      Object.keys(this.filterCivilization).forEach(civ => this.filterCivilization[civ] = false);
      Object.keys(this.filterMana).forEach(manaCost => this.filterMana[manaCost] = false);
      this.filterFamily = ALL_FAMILIES;
      this.filterSet = ALL_SETS;
      this.filterCard = '';
    },

    getCardsForDeck(cardUids) {
      let cards = [];
      for (let uid of cardUids) {
        let card = this.cards.find(x => x.uid === uid)
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
    },

    tryAddCard(card) {
      if (
        this.selectedDeck.cards.filter(x => x == card.uid).length >= 4
      ) {
        if (!permissions().includes("admin")) {
          playSound("/assets/sounds/card-limit.wav");
          return;
        }
      }

      this.selectedDeck.cards.push(card.uid);
      playSound("/assets/sounds/card-added.mp3");
    },

    tryRemoveCard(card) {
      let uid = card.uid;
      let toSlice = this.selectedDeck.cards.indexOf(uid);
      if (toSlice < 0) return;

      playSound("/assets/sounds/card-removed.wav");
      this.selectedDeck.cards.splice(toSlice, 1);
      if (this.selectedDeck.cards.indexOf(uid) < 0 && this.previewCard.uid === uid ) {
        this.previewCard = null 
      }
    },

    modifyCatalogueCardSize(addition) {
      this.cardSize += addition
    },

    copyDeckList() {
      const cards = this.getCardsForDeck(this.selectedDeck.cards);

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

    displayFamily(family) {
      return family ? family.join(" / ") : "Spell"
    },
  },
  async created() {
    try {
      let [cards, decks] = await Promise.all([
        call({ path: "/cards", method: "GET" }),
        call({ path: "/decks", method: "GET" })
      ]);
      let sets = {};
      let cardsCiv = {
        "light": [],
        "darkness": [],
        "nature": [],
        "fire": [],
        "water": [],
      };
      for (let card of cards.data) {
        cardsCiv[card.civilization.toLowerCase()].push(card);
        if (!sets[card.set]) {
          sets[card.set] = true;
        }
      }
      this.sets = Object.keys(sets);
      this.sets.push(ALL_SETS);
      this.sets.sort();

      let sortedCards = [];
      Object.values(cardsCiv).forEach(
        civSet => sortedCards.push(...civSet.sort((c1, c2) => 
          compareCards(c1, c2, { by: "manaCost", directionNum: 1 })
        ))
      );

      this.cards = sortedCards;
      this.decks = decks.data;

      if (this.decks.length < 1) {
        this.decks.push({
          name: "My first deck",
          cards: [],
          public: false
        });
      }

      let families = [];
      for (let c of this.cards) {
        if (c.family) {
          for (let f of c.family) {
            if (!families.includes(f)) {
              families.push(f);
            }
          }
        }
      }
      families.sort();
      this.families.push(...families);

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
  computed: {
    filteredAndSortedCards() {
      let cards = this.cards.filter(card =>
        card.name.toLowerCase().includes(this.filterCard.toLowerCase()) ||
          card.text.toLowerCase().includes(this.filterCard.toLowerCase())
      );

      if (this.filterSet !== ALL_SETS) {
        cards = cards.filter(card => card.set === this.filterSet);
      }

      let filterCivilizationValues = Object.values(this.filterCivilization);
      if (!(filterCivilizationValues.every(v => v === filterCivilizationValues[0]))) {
        cards = cards.filter(
          card => this.filterCivilization[card.civilization] 
        );
      }

      if (this.filterFamily.toLowerCase() !== ALL_FAMILIES.toLowerCase()) {
        cards = cards.filter(
          card => 
            (this.filterFamily.toLowerCase() === "spell" && !card.family) || 
            (card.family && card.family.includes(this.filterFamily))
        );
      }

      let filterManaValues = Object.values(this.filterMana);
      if (!(filterManaValues.every(v => v === filterManaValues[0]))) {
        cards = cards.filter(
          card => {
            if (card.manaCost > 6)
              return this.filterMana['7+'];
            else
              return this.filterMana[card.manaCost.toString()];
          } 
        );
      }

      return cards;
    },
  }
};
</script>

<style scoped lang="scss">

$catalogue-main-color: #b8b7b7;
$catalogue-secondary-color: #e0dede;
$civ-nature-base-color: #118141;
$civ-fire-base-color: #D12027;
$civ-water-base-color: #47C6F2;
$civ-light-base-color: #FAD241;
$civ-darkness-base-color: #65696C;

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

.fl {
  float: left;
}

.zoom-icons {
  display: inline-flex;
  align-items: flex-end;
  align-self: end;
  padding: 5px 10px 5px 10px;
  margin-right: 10px;
  margin-left: 5px;
  justify-content: space-between;
  gap: 10px;
  background-color: $catalogue-main-color;
}

.zoom-icon {
  width: 25px;
  cursor: pointer;
}

.reset-icon {
  width: 35px;
  margin-left: 20px;
  background-color: lightsalmon;
  cursor: pointer;
}

.civ-filter {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-block;
  margin-left: 5px;
  cursor: pointer;
  border: 4px solid $catalogue-secondary-color;
}

.civ-filter-selected {
  border: 4px solid black;
}

.civ-color-fire {
  background-color: $civ-fire-base-color;
}

.civ-color-water {
  background-color: $civ-water-base-color;
}

.civ-color-light {
  background-color: $civ-light-base-color;
}

.civ-color-darkness {
  background-color: $civ-darkness-base-color;
}

.civ-color-nature {
  background-color: $civ-nature-base-color;
}

.mana-filter {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: 5px;
  cursor: pointer;
  color: white  ;
  background-color: black;
  font-weight: 600;
  border: 4px solid $catalogue-secondary-color;
}

.mana-filter-selected {
  border: 4px solid orange;
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
  height: calc(100% - 20px);
  margin-top: 10px;
  overflow-y: auto;
}

.card-preview {
  width: 300px;
  text-align: center;
  border-radius: 4px;
  height: 420px;
  z-index: 2005;
  position: absolute;
  left: calc(50% - 150px);
  top: calc(50vh - 240px);
}

.card-preview > img {
  width: 300px;
  border-radius: 15px;
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

.deck-card-slot {
  display: flex;
  margin-bottom: 4px;
  align-items: center;
  width: 100%;
}

.deck-card-oval {
  border-radius: 10px / 20px;
  color: #000;
  padding-left: 7px;
  padding-right: 3px;
  border: 3px solid;
  flex-grow: 1;
  display: inline-flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  min-width: 0;
}

.deck-card-name {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  font-size: 17px;
  font-family: "Perceval Bold", "Roboto", sans-serif;
}

.deck-card-mana {
  display: inline-block;
  background-color: black;
  color: white;
  font-weight: 600;
  border-radius: 50%;
  padding-left: 7px;
  padding-right: 7px;
  margin-top: 2px;
  margin-bottom: 2px;
}

.deck-card-counter {
  border-radius: 10px / 20px;
  color: black;
  padding-left: 6px;
  padding-right: 9px;
  padding-top: 4px;
  padding-bottom: 4px;
  background: radial-gradient(40% 100% at right, transparent 50%, lightgray 51%);
  display: inline-block;
  flex-grow: 0;
  margin-right: -3px;
}

.card-fire {
  border-color: $civ-fire-base-color;
  background: radial-gradient(#F0DDD5 0%, #E39E7F 80%, $civ-fire-base-color 100%);
}

.card-nature {
  border-color: $civ-nature-base-color;
  background: radial-gradient(#E1EBD0 0%, #84C188 80%, $civ-nature-base-color 100%);
}

.card-water {
  border-color: $civ-water-base-color;
  background: radial-gradient(#D4EDF8 0%, #93D4DE 80%, $civ-water-base-color 100%);
}

.card-darkness {
  border-color: $civ-darkness-base-color;
  background: radial-gradient(#CCD1D5 0%, #CACFD2 80%, $civ-darkness-base-color 100%);
}

.card-light {
  border-color: $civ-light-base-color;
  background: radial-gradient(#FCF7D0 0%, #FBF396 80%, $civ-light-base-color 100%);
}

.right {
  display: inline-block;
  height: calc(100vh - 110px);
  overflow-y: auto;
  width: calc(20% - 10px);
  max-width: 380px;
  flex-shrink: 0;
}

.catalogue {
  display: inline-block;
  height: calc(100vh - 115px);
  min-width: calc(80% - 10px);
  background-color: $catalogue-main-color;
  margin-right: 20px;
  flex-grow: 1;
}

.catalogue-cards-wrapper {
  overflow-y: auto;
  width: 100%;
  height: 100%;
  background-color: $catalogue-main-color;
}

.catalogue-cards {
  width: calc(100% - 10px);
  display: inline-flex;
  flex-wrap:wrap;
  gap: 5px;
  justify-content: space-evenly;
  
  padding: 5px;
}

.catalogue-top-bar {
  display: flex;
  justify-content: space-between;
  background-color: $catalogue-secondary-color;
}

.catalogue-filters {
  margin-left: 10px;
  flex-wrap: wrap;
  display: inline-flex;
  align-items: center;
  padding: 3px 0;
}

.catalogue-filter {
  background-color: black;
  color: white;
}

.catalogue-card {
  aspect-ratio: 264/367;
  margin-bottom: 5px;
  cursor: pointer;
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
