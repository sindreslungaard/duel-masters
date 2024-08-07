<template>
  <div>
    <div v-show="warning" @click="closeOverlay()" class="overlay"></div>

    <span v-if="previewCard">
      <div class="card-preview">
        <img
          :src="`https://scans.shobu.io/${previewCard.uid}.jpg`"
          alt="Full card"
        />
      </div>
    </span>

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
      <div class="left bg">
        <div class="cards-table">
          <table>
            <tr>
              <th>
                <div class="sort-btn" @click="toggleSort('name')">
                  Card Name
                  <img
                    class="sort-ico"
                    width="25px"
                    :src="`/assets/images/${sortIcons.name}.png`"
                  />
                </div>
                <input
                  v-model="filterCardName"
                  type="search"
                  placeholder="Type to search"
                />
              </th>
              <th>
                <div class="sort-btn" @click="toggleSort('set')">
                  Set
                  <img
                    class="sort-ico"
                    width="25px"
                    :src="`/assets/images/${sortIcons.set}.png`"
                  />
                </div>
                <select v-model="filterSet">
                  <option
                    class="set"
                    v-for="(set, index) in sets"
                    :key="index"
                    :value="set"
                    >{{ set }}</option
                  >
                </select>
              </th>
              <th>
                <div class="sort-btn" @click="toggleSort('civilization')">
                  Civilization
                  <img
                    class="sort-ico"
                    width="25px"
                    :src="`/assets/images/${sortIcons.civilization}.png`"
                  />
                </div>
                <select v-model="filterCivilization">
                  <option value="all">All</option>
                  <option value="fire">Fire</option>
                  <option value="water">Water</option>
                  <option value="nature">Nature</option>
                  <option value="light">Light</option>
                  <option value="darkness">Darkness</option>
                </select>
              </th>
              <th>
                <div class="sort-btn" @click="toggleSort('manaCost')">
                  Mana
                  <img
                    class="sort-ico"
                    width="25px"
                    :src="`/assets/images/${sortIcons.manaCost}.png`"
                  />
                </div>
              </th>
              <th>
                <div class="sort-btn" @click="toggleSort('family')">
                  Race
                  <img
                    class="sort-ico"
                    width="25px"
                    :src="`/assets/images/${sortIcons.family}.png`"
                  />
                </div>
                <select v-model="filterFamily">
                  <option
                    class="family"
                    v-for="(family, index) in families"
                    :key="index"
                    :value="family"
                    >{{ family }}</option
                  >
                </select>
              </th>
            </tr>
            <tr
              @contextmenu.prevent="previewCard = card"
              @click="
                selectedFromDeck = card;
                selected = card;
              "
              v-for="(card, index) in filteredAndSortedCards"
              :key="index"
              :class="[{ selected: selected && selected.uid === card.uid }]"
            >
              <td
                @mouseover="previewCard = card"
                @mouseleave="previewCard = null"
              >
                {{ card.name }}
              </td>
              <td class="set">{{ card.set }}</td>
              <td class="civilization">{{ card.civilization }}</td>
              <td class="manaCost">{{ card.manaCost }}</td>
              <td class="family">{{ displayFamily(card.family) }}</td>
            </tr>
          </table>
        </div>
      </div>

      <div class="middle">
        <div
          v-show="!showWizard"
          @click="tryAddCard()"
          class="arrow-green"
        ></div>
        <div
          v-show="!showWizard"
          @click="tryRemoveCard()"
          class="arrow-red"
        ></div>
      </div>

      <div class="right">
        <select v-model="selectedDeckUid" class="fl" style="margin: 0;">
          <option
            v-for="(deck, index) in decks"
            :key="index"
            :value="deck.uid"
            >{{ deck.name }}</option
          >
        </select>
        <img
          @click="showWizard = true"
          class="fl edit-ico"
          width="25px"
          src="/assets/images/edit_icon.png"
        />
        <div class="right-btns">
          <a
            :href="'/deck/' + selectedDeckUid"
            v-if="selectedDeck && selectedDeck.public"
            target="_blank"
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
          >
            <img
              class="fl edit-ico share"
              width="25px"
              src="/assets/images/list_icon.png"
            />
          </a>

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

        <div v-if="selectedDeck" class="right-content bg">
          <div class="cards-table">
            <table>
              <tr>
                <th>
                  <div class="sort-btn" @click="toggleSortForDeck('count')">
                    Quantity ({{ selectedDeck.cards.length }})
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.count}.png`"
                    />
                  </div>
                </th>

                <th>
                  <div class="sort-btn" @click="toggleSortForDeck('name')">
                    Card Name
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.name}.png`"
                    />
                  </div>
                </th>

                <th>
                  <div class="sort-btn" @click="toggleSortForDeck('set')">
                    Set
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.set}.png`"
                    />
                  </div>
                </th>

                <th>
                  <div
                    class="sort-btn"
                    @click="toggleSortForDeck('civilization')"
                  >
                    Civilization
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.civilization}.png`"
                    />
                  </div>
                </th>

                <th>
                  <div class="sort-btn" @click="toggleSortForDeck('manaCost')">
                    Mana
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.manaCost}.png`"
                    />
                  </div>
                </th>

                <th>
                  <div class="sort-btn" @click="toggleSortForDeck('family')">
                    Race
                    <img
                      class="sort-ico"
                      width="25px"
                      :src="`/assets/images/${sortDeckIcons.family}.png`"
                    />
                  </div>
                </th>
              </tr>
              <template v-if="selectedDeck">
                <tr
                  @contextmenu.prevent="previewCard = card"
                  @click="
                    selected = card;
                    selectedFromDeck = card;
                  "
                  v-for="(card, index) in getCardsForDeck(selectedDeck.cards)"
                  :key="index"
                  :class="[
                    {
                      selected:
                        selectedFromDeck && selectedFromDeck.uid === card.uid
                    }
                  ]"
                >
                  <td>{{ card.count }}</td>
                  <td
                    @mouseover="previewCard = card"
                    @mouseleave="previewCard = null"
                  >
                    {{ card.name }}
                  </td>
                  <td class="set">{{ card.set }}</td>
                  <td class="civilization">{{ card.civilization }}</td>
                  <td class="manaCost">{{ card.manaCost }}</td>
                  <td class="family">{{ displayFamily(card.family) }}</td>
                </tr>
              </template>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";
import Header from "../components/Header.vue";
import config from "../config";

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

export default {
  name: "decks",
  components: {
    Header
  },
  computed: {
    username: () => localStorage.getItem("username")
  },
  data() {
    return {
      warning: "",
      showWizard: false,

      filterCardName: "",
      filterCivilization: "all",
      filterFamily: "All",
      filterSet: "All",
      families: ["All", "Spell"],
      sort: {
        by: "name",
        directionNum: 1
      },
      sortForDeck: {
        by: "name",
        directionNum: 1
      },

      sets: [],
      cards: [],
      selected: null,
      selectedFromDeck: null,

      decks: [],
      selectedDeck: null,
      selectedDeckUid: null,
      deckCopy: null,

      previewCard: null
    };
  },
  methods: {
    selectDeck(deck) {
      this.selectedDeck = deck;
      this.deckCopy = JSON.parse(JSON.stringify(deck));
    },

    cardInfo(uid) {
      let card = this.cards.find(x => x.uid === uid);
      return card;
    },

    toggleSort(by) {
      this.sort = {
        directionNum: this.sort.by === by ? -this.sort.directionNum : 1,
        by
      };
    },

    toggleSortForDeck(by) {
      this.sortForDeck = {
        directionNum:
          this.sortForDeck.by === by ? -this.sortForDeck.directionNum : 1,
        by
      };
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

      cards.sort((c1, c2) => compareCards(c1, c2, this.sortForDeck));

      return cards;
    },

    tryAddCard() {
      if (!this.selected) {
        return;
      }

      if (
        this.selectedDeck.cards.filter(x => x == this.selected.uid).length >= 4
      ) {
        if (!permissions().includes("admin")) {
          return;
        }
      }

      this.selectedDeck.cards.push(this.selected.uid);
    },

    tryRemoveCard() {
      if (!this.selectedFromDeck) {
        return;
      }

      let toSlice = -1;
      for (let i = 0; i < this.selectedDeck.cards.length; i++) {
        if (this.selectedDeck.cards[i] === this.selectedFromDeck.uid) {
          toSlice = i;
        }
      }
      if (toSlice < 0) {
        return;
      }

      this.selectedDeck.cards.splice(toSlice, 1);
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
      for (let card of cards.data) {
        if (!sets[card.set]) {
          sets[card.set] = true;
        }
      }
      this.sets = Object.keys(sets);
      this.sets.push("All");
      this.sets.sort();

      this.cards = cards.data;
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
        card.name.toLowerCase().includes(this.filterCardName.toLowerCase())
      );

      if (this.filterSet !== "All") {
        cards = cards.filter(card => card.set === this.filterSet);
      }

      if (this.filterCivilization !== "all") {
        cards = cards.filter(
          card => card.civilization === this.filterCivilization
        );
      }

      if (this.filterFamily.toLowerCase() !== "all") {
        cards = cards.filter(
          card => 
            (this.filterFamily.toLowerCase() === "spell" && !card.family) || 
            (card.family && card.family.includes(this.filterFamily))
        );
      }

      cards.sort((c1, c2) => compareCards(c1, c2, this.sort));

      return cards;
    },

    sortIcons() {
      const result = {
        name: "arrow_up_down",
        set: "arrow_up_down",
        civilization: "arrow_up_down",
        manaCost: "arrow_up_down",
        family: "arrow_up_down"
      };

      result[this.sort.by] =
        this.sort.directionNum === 1 ? "arrow_down" : "arrow_up";

      return result;
    },

    sortDeckIcons() {
      const result = {
        count: "arrow_up_down",
        name: "arrow_up_down",
        set: "arrow_up_down",
        civilization: "arrow_up_down",
        manaCost: "arrow_up_down",
        family: "arrow_up_down"
      };

      result[this.sortForDeck.by] =
        this.sortForDeck.directionNum === 1 ? "arrow_down" : "arrow_up";

      return result;
    }
  }
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

.arrow-green {
  width: 40px;
  height: 40px;
  margin: 0 auto;
  background: url("/assets/images/arrow_green_icon.png");
  background-size: cover;
  transform: rotate(0deg) scaleX(-1);
  opacity: 0.8;
}

.arrow-red {
  width: 40px;
  height: 40px;
  margin: 0 auto;
  background: url("/assets/images/arrow_red_icon.png");
  background-size: cover;
  transform: rotate(0deg) scaleX(-1);
  opacity: 0.8;
}

.arrow-green:hover {
  cursor: pointer;
  opacity: 1;
}

.arrow-red:hover {
  cursor: pointer;
  opacity: 1;
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

.right-btns {
  text-align: right;
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
  width: 100%;
  height: calc(100% - 40px);
  margin-top: 10px;
  overflow: auto;
}

.right-content table {
  padding-top: 10px !important;
}

.card-preview {
  width: 300px;
  text-align: center;
  border-radius: 4px;
  height: 480px;
  z-index: 2005;
  position: absolute;
  left: calc(50% - 150px);
  top: calc(50vh - 240px);
}

.card-preview > img {
  width: 300px;
  border-radius: 15px;
  margin-bottom: 10px;
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

th {
  input,
  select {
    margin: 0;
    width: 100% !important;
  }
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

table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
  font-size: 14px;
  user-select: none;
  max-height: calc(100vh - 115px);
  overflow: scroll;
}

td,
th {
  border: 1px solid #1D202A;
  text-align: left;
  padding: 8px;
  cursor: pointer;
}

.bg {
  background: url(/assets/images/overlay_30.png);
  border-radius: 4px;
}

.left,
.right {
  display: inline-block;
  height: calc(100vh - 115px);
  overflow: auto;
}

.cards-table {
  margin: 10px;
  width: calc(100% - 10px * 2);
}

.left {
  width: calc(50% - 100px / 2);
}

.middle {
  display: inline-block;
  width: 100px;
  text-align: center;
  overflow: auto;
  height: calc(52vh - 115px / 2);
}

.right {
  width: calc(50% - 100px / 2);
}

.main {
  margin: 0 15px;
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

.sort-btn {
  white-space: nowrap;
  margin-bottom: 10px;
  display: flex;
  align-items: center;

  .sort-ico {
    margin-left: 5px;
    filter: invert(70%);

    &:hover {
      filter: invert(80%);
      cursor: pointer;
    }
  }
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

.set {
  text-transform: uppercase;
}

.family {
}

.civilization {
  text-transform: capitalize;
}
</style>
