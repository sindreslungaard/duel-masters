<template>
  <div>

      <div v-show="warning || previewCard" class="overlay"></div>

      <div v-if="previewCard" class="card-preview">
        <img :src="`/assets/cards/all/${previewCard.uid}.jpg`">
        <div @click="previewCard = null" class="btn">Close</div>
      </div>

        <div v-show="warning" class="error">
            <p>{{ warning }}</p>
            <div @click="warning = ''" class="btn">Close</div>
        </div>

        <div v-if="showWizard" class="new-duel">
          <div class="backdrop"></div>
          <div class="wizard">
              <div class="spacer">
                <span class="headline">Edit your deck</span>
                <br><br>
                <form>
                    <span class="helper">Name</span>
                    <input v-model="selectedDeck.name" type="text" placeholder="Name">
                    <br><br>
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
                            Card Name
                            <input v-model="filterCardName" type="text" placeholder="Type to search">
                        </th>
                        <th>
                            Set
                            <select>
                                <option value="all">All</option>
                                <option value="dm-01">dm-01</option>
                            </select>
                        </th>
                        <th>
                            Civilization
                            <select v-model="filterCivilization">
                                <option value="all">All</option>
                                <option value="fire">Fire</option>
                                <option value="water">Water</option>
                                <option value="nature">Nature</option>
                                <option value="light">light</option>
                                <option value="darkness">Darkness</option>
                            </select>
                        </th>
                    </tr>
                    <tr @dblclick="previewCard = card" @contextmenu.prevent="previewCard = card" @click="selectedFromDeck = null; selected = card" v-for="(card, index) in cardsFiltered" :key="index" :class="[{ 'selected': selected === card }]">
                        <td>{{ card.name }}</td>
                        <td>{{ card.set }}</td>
                        <td>{{ card.civilization }}</td>
                    </tr>
                    </table>
              </div>
          </div>

          <div class="middle">
              <div v-show="!showWizard" @click="tryAddCard()" class="arrow-green"></div>
              <div v-show="!showWizard" @click="tryRemoveCard()" class="arrow-red"></div>
          </div>

          <div class="right">
              <select v-model="selectedDeckUid" class="fl" style="margin: 0;">
                <option v-for="(deck, index) in decks" :key="index" :value="deck.uid">{{ deck.name }}</option>
            </select>
            <img @click="showWizard = true" class="fl edit-ico" width="25px" src="/assets/images/edit_icon.png">
              <div class="right-btns">
                  <a :href="getShareUrl(selectedDeckUid)" v-if="selectedDeck.public" target="_blank"><img class="fl edit-ico share" width="25px" src="/assets/images/share_icon.png"></a>
                  <div @click="newDeck()" class="btn new">New Deck</div>
                  <template v-if="selectedDeck && deckCopy && !decksEqual(selectedDeck, deckCopy)">
                      <div @click="save()" class="btn save">Save</div>
                      <div @click="discard()" class="btn discard">Discard</div>
                  </template>
              </div>

              <div class="right-content bg">
                  <div class="cards-table">
                    <table>
                        <tr>
                            <th>Quantity</th>
                            <th>Card Name</th>
                            <th>Set</th>
                            <th>Civilization</th>
                        </tr>
                        <template v-if="selectedDeck">
                            <tr @dblclick="previewCard = card" @contextmenu.prevent="previewCard = card" @click="selected = null; selectedFromDeck = card" v-for="(card, index) in getCardsForDeck(selectedDeck.cards)" :key="index" :class="[{ 'selected': selectedFromDeck && selectedFromDeck.uid === card.uid }]">
                                <td>{{ card.count }}</td>
                                <td>{{ card.name }}</td>
                                <td>{{ card.set }}</td>
                                <td>{{ card.civilization }}</td>
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
import { call } from '../remote'
import Header from '../components/Header.vue'
import config from '../config'

export default {
  name: 'decks',
  components: {
      Header
  },
  computed: {
      username: () => localStorage.getItem('username')
  },
  data() {
      return {
          warning: "",
          showWizard: false,

          filterCardName: "",
          filterCivilization: "all",

          cards: [],
          cardsFiltered: [],
          selected: null,
          selectedFromDeck: null,

          decks: [],
          selectedDeck: null,
          selectedDeckUid: null,
          deckCopy: null,

          previewCard: null
      }
  },
  methods: {
      selectDeck(deck) {
          this.selectedDeck = deck
          this.deckCopy = JSON.parse(JSON.stringify(deck))
      },

      cardInfo(uid) {
          let card = this.cards.find(x => x.uid === uid)
          return card
      },

      getCardsForDeck(cardUids) {
          let cards = []
          for(let uid of cardUids) {
              let card = JSON.parse(JSON.stringify(this.cards.find(x => x.uid === uid)))

              let existingCard = cards.find(x => x.uid === card.uid)
              if(existingCard) {
                  existingCard.count += 1
              } else {
                  card.count = 1
                  cards.push(card)
              }
          }
          return cards
      },

      tryAddCard() {
          if(!this.selected) {
              return
          }
          this.selectedDeck.cards.push(this.selected.uid)
      },

      getShareUrl(uid) {
          return window.location.host + "/deck/" + uid
      },

      tryRemoveCard() {
          if(!this.selectedFromDeck) {
              return
          }

          let toSlice = -1
          for(let i = 0; i < this.selectedDeck.cards.length; i++) {
              if(this.selectedDeck.cards[i] === this.selectedFromDeck.uid) {
                  toSlice = i
              }
          }
          if(toSlice < 0) {
              return
          }

          this.selectedDeck.cards.splice(toSlice, 1)
      },

      newDeck() {
          if(!this.decksEqual(this.selectedDeck, this.deckCopy)) {
              this.warning = "Please save or discard the changes you've made before creating a new deck"
              return
          }
          this.decks.push({
              name: "Unnamed Deck",
              cards: [],
              public: false
          })
          this.deckCopy = JSON.parse(JSON.stringify(this.decks[this.decks.length - 1]))
          this.selectedDeck = this.decks[this.decks.length - 1]
          this.selectedDeckUid = this.selectedDeck.uid
          this.$nextTick(() => {
              this.deckCopy.name = "to.be.removed"
          })
      },

      async save() {
          try {
              let res = await call({
                  path: "/deck",
                  method: "POST",
                  body: this.selectedDeck
              })
              this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck))
              this.warning = "Successfully saved your deck"
          } catch(e) {
              this.warning = e.response.data.message
          }
      },

      discard() {
          if(this.deckCopy.name === "to.be.removed") {
              this.selectedDeck = this.decks[0]
              this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck))
              this.selectedDeckUid = this.selectedDeck.uid
              this.decks.pop()
              return
          }

          this.selectedDeck = JSON.parse(JSON.stringify(this.deckCopy))
      },

      decksEqual(deck1, deck2) {
          if(deck1.name !== deck2.name) {
              return false
          }
          if(deck1.public !== deck2.public) {
              return false
          }
          if(deck1.cards.length !== deck2.cards.length) {
              return false
          }
          for(let i = 0; i < deck1.cards.length; i++) {
              if(deck1.cards[i] !== deck2.cards[i]) {
                  return false
              }
          }
          return true
      },

      filter() {
        let filtered = this.cards.filter(x => x.name.toLowerCase().includes(this.filterCardName.toLowerCase()))
        if(this.filterCivilization !== "all") {
            filtered = filtered.filter(x => x.civilization === this.filterCivilization)
        }
        
          this.cardsFiltered = filtered
      }
  },
  async created() {
      try {
          let [cards, decks] = await Promise.all([
            call({ path: "/cards", method: "GET" }),
            call({ path: "/deck", method: "GET" })
          ])

          this.cards = cards.data
          this.cardsFiltered = this.cards
          this.decks = decks.data

          if(this.decks.length < 1) {
              this.decks.push({
                  name: "My first deck",
                  cards: [],
                  public: false,
              })  
          }
          this.selectedDeck = this.decks[0]
          this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck))
          this.selectedDeckUid = this.selectedDeck.uid
      }
      catch(e) {
          console.log(e)
      }
  },
  watch: {
    filterCardName: function(val) {
      this.filter()
    },
    filterCivilization: function(val) {
      this.filter()
    },
    selectedDeckUid: function(val) {
        if(!this.decksEqual(this.selectedDeck, this.deckCopy)) {
            this.warning = "You have unsaved changes in the currently selected deck. Save or discard before editing another deck."
            this.selectedDeckUid = this.selectedDeck.uid
            return
        }
        this.selectedDeck = this.decks.find(x => x.uid === val)
        this.deckCopy = JSON.parse(JSON.stringify(this.selectedDeck))
    }
  }
}
</script>

<style scoped>

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
    background: url('/assets/images/arrow_green_icon.png');
    background-size: cover;
    transform: rotate(0deg) scaleX(-1);
    opacity: 0.8;
}

.arrow-red {
    width: 40px;
    height: 40px;
    margin: 0 auto;
    background: url('/assets/images/arrow_red_icon.png');
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
    background: #3CA374 !important;
    display: inline-block;
    margin-right: 15px;
}

.save:hover {
    background: #3CA374 !important;
}

.discard {
    background: #FF4C4C !important;
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
  left: calc(50% - 300px / 2);
  top: calc(50vh - 480px / 2);
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
    background: #7289DA;
    border-color: #7289DA !important;
}

.right select {
    padding: 7px !important;
}

input, textarea, select {
    border: none;
    background: #484C52;
    padding: 5px !important;
    width: auto !important;
    margin-left: 5px;
    border-radius: 4px;
    color: #ccc;
    resize: none;
}
input:focus, textarea:focus, select:focus {
    outline: none
}
input:active, textarea:active, select:active {
    outline: none
}

.left::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    border-radius: 10px;
    background-color: #484C52;
  }

.left::-webkit-scrollbar {
  width: 6px;
  background-color: #484C52;
}

.left::-webkit-scrollbar-thumb {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    background-color: #222;
}

.right::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    border-radius: 10px;
    background-color: #484C52;
  }

.right::-webkit-scrollbar {
  width: 6px;
  background-color: #484C52;
}

.right::-webkit-scrollbar-thumb {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    background-color: #222;
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

td, th {
  border: 1px solid #4C4C4C;
  text-align: left;
  padding: 8px;
  cursor: pointer;
}

.bg {
    background: #2B2C31;
    border-radius: 4px;
}

.left, .right {
    display: inline-block;
    height: calc(100vh - 115px);
    overflow: auto;
}

.cards-table {
    margin: 10px;
    width: calc(100% - 10px * 2)
}

.left { 
    width: calc(50% - 100px / 2);
}

.middle {
    display: inline-block;
    width: 100px;
    text-align: center;
    overflow: auto;
    height: calc(52vh - 115px / 2)
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
    background: #36393F;
    width: 250px;
    border-radius: 4px;
    color: #fff;
    border: 1px solid #666;
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
    background: #FF4C4C;
    color: #fff;
}

.wizard .cancel:hover {
    background: #ed3e3e;
}

input, textarea, select {
    border: none;
    background: #484C52;
    padding: 10px;
    border-radius: 4px;
    width: 200px;
    color: #ccc;
    resize: none;
}
input:focus, textarea:focus, select:focus {
    outline: none
}
input:active, textarea:active, select:active {
    outline: none
}

.wizard select {
    width: 220px;
    margin-top: 4px;
}

.error p {
  padding: 5px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: #2B2E33 !important;
  border: 1px solid #222428;
}

.error {
  border: 1px solid #666;
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #36393F;
  z-index: 1000;
  left: calc(50% - 300px / 2);
  top: 40vh;
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
    cursor: default
}

.title {
    position: absolute;
    top: 16px;
    left: 16px;
}

.psa {
    margin: 16px;
    background: #2B2C31;
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
    color: #7289DA;
}

.btn {
  display: inline-block;
  background: #7289DA;
  color: #E3E3E5;
  font-size: 14px;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn:hover {
  cursor: pointer;
  background: #677BC4;
}

.btn:active {
  background: #5B6EAE !important;
}

</style>