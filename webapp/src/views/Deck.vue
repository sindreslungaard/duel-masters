<template>
  <div>
    <div class="main">
      <h1 v-if="errorMsg">{{ errorMsg }}</h1>
      <h1>{{ name }}</h1>
      <hr />
      <div class="deck-container">
        <div v-for="(card, index) in cards" :key="index" class="deck-card">
          <p>x{{ card.count }}</p>
          <img :src="`https://scans.shobu.io/${card.uid}.jpg`" />
        </div>
      </div>
      <div class="button-container" v-if="isUserAuthenticated()">
        <Button class="btn" :disabled='apiCallInProgess' @click='saveCardsToDeck'>Copy to decklist</Button>
      </div>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";
import { showMessageDialog } from "../helpers/showDialog"

export default {

  name: "deck",

  data() {
    return {
      errorMsg: null,
      name: null,
      cards: [],
      apiCallInProgess: false,
    };
  },

  methods: {

    getUsername() {
      return localStorage.getItem("username");
    },

    getUserUid() {
      return localStorage.getItem("uid");
    },

    isUserAuthenticated() {
      return this.getUserUid() !== null
    },

    saveCardsToDeck() {
      let cardsUids = this.getUidListFromCards(this.cards);
      this.saveNewDeck(this.name, cardsUids);
    },

    getUidListFromCards(cards){
      let uids = [];
      cards.forEach(card => {
        for(let i=0;i<card.count;i++){
          uids.push(card.uid)
        }
      });
      return uids;
    },

    saveNewDeck(deckName, cardUids) {
      this.apiCallInProgess = true;
      call({
        path: "/decks",
        method: "POST",
        body: {
          name: deckName,
          cards: cardUids,
          public: false,
        }
      }).then(res => {
        showMessageDialog(this.$modal, "Successfully saved deck");
      }).catch(error => {
        showMessageDialog(this.$modal, "Invalid request. If you think this is an error please contact us via discord.");
      }).finally(() => {
        this.apiCallInProgess = false
      });
    },

    closeOverlay() {
      this.overlayMessage = null;
    },
  },

  async created() {
    try {
      let res = await call({
        path: `/deck/${this.$route.params.uid}`,
        method: "GET"
      });
      this.name = res.data.name + " by " + res.data.owner;

      for (let card of res.data.cards) {
        let exists = this.cards.find(x => x.uid === card);

        if (exists) {
          exists.count += 1;
        } else {
          this.cards.push({
            uid: card,
            count: 1
          });
        }
      }
    } catch (e) {
      this.errorMsg = e.response ? e.response.data.message : e;
    }
  }
};
</script>

<style scoped>

.button-container {
  text-align: center;
}

.deck-card > p {
  text-align: center;
  margin: 0;
  margin-bottom: 5px;
  color: #ccc;
}

.deck-container {
  margin-top: 20px;
  margin-bottom: 20px;
}

.deck-card {
  display: inline-block;
  width: 18%;
  margin-right: 2.5%;
  margin-bottom: 2%;
}

.deck-card:nth-of-type(5n) {
  margin-right: 0;
}

.deck-card > img {
  width: 100%;
  border-radius: 4%;
}

.main {
  margin: 100px 105px !important;
}

.main > h1 {
  margin: 0;
  padding: 0;
}

hr {
  border-color: #888;
}

.main {
  margin: 0 15px;
}

input,
textarea,
select {
  border: none;
  background: #484c52;
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

a {
  color: #7289da;
}

.btn {
  display: inline-block;
  background: #7289da;
  color: #e3e3e5;
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
  background: #677bc4;
}

.btn:active {
  background: #5b6eae;
}
</style>
