<template>
  <div>
    <div v-show="warning" @click="closeOverlay()" class="overlay"></div>
    <div v-show="warning" class="error">
      <p class="whitespace-pre-wrap">{{ warning }}</p>
      <div @click="warning = ''" class="btn">Close</div>
    </div>

    <Header></Header>

    <div class="main">
      <CardsCatalogue 
        class="catalogue"
        :unsortedCards="cardPool"
        :eventMode="true"
        v-model:deck="deckCards"
      />

      <div class="right">  
        <!-- Add event name and deck size -->
        <span class="deck-card-total mr-4">({{deckCards.length}})</span>
        <div class="deck-crud-buttons">
          <div @click="save()" class="btn save">Save</div>
          <!-- <div @click="discard()" class="btn discard">Discard</div> -->
        </div>
        <div class="right-content">
          <DeckList 
            :cards="cardPool"
            v-model:deck="deckCards"
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
import DeckList from "../components/DeckList.vue";

export default {
  name: "event_deck",
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

      allCards: [],
      deckCards: [],
      cardPool: [],
    };
  },
  methods: {
    closeOverlay() {
      this.warning = null;
    },
    async save() {
      try {
        await call({
          path: `/event/${this.$route.params.id}/deck`,
          method: "PUT",
          body: {
            cards: this.deckCards
          },
        })
        this.warning = "Successfully saved your deck";
      } catch (e) {
        console.log(e)
        this.warning =
          "Invalid request. Please ensure that the deck has 40 cards.";
      }
    },

  },
  async created() {
    try {
      let [cards, deck] = await Promise.all([
        call({ path: "/cards", method: "GET" }),
        call({
          path: `/event/${this.$route.params.id}/deck`,
          method: "GET"
        })
      ]);
      // this.allCards = cards.data;
      this.deckCards = deck.data.cards;
      
      let cardPoolMap = {}
      deck.data.cardpool.forEach(id => {
        if (cardPoolMap[id] == null) {
          let card = cards.data.find(x => x.uid == id)
          cardPoolMap[card.uid] = card
          cardPoolMap[card.uid].count = 1
        }
        else {
          cardPoolMap[id].count += 1
        }
      })
      this.cardPool = Object.values(cardPoolMap)  
    } catch (e) {
      console.log(e);
    }
  },
};
</script>

<style scoped lang="scss">
.save {
  background: #3ca374 !important;
  display: inline-block;
  margin-right: 15px;
}

.save:hover {
  background: #3ca374 !important;
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

a {
  color: #7289da;
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

.deck-card-total {
  padding: 2px;
  background-color: lightgray;
  color: black;
}
</style>
