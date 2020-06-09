<template>
  <div>
    <div class="main">
      <h1 v-if="errorMsg">{{ errorMsg }}</h1>
      <h1>{{ name }}</h1>
      <hr />
      <div class="deck-container">
        <div v-for="(card, index) in cards" :key="index" class="deck-card">
          <p>{{ card.name }} [x{{ card.count }}]</p>
          <img :src="`/assets/cards/all/${card.uid}.jpg`" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { call } from "../remote";

export default {
  name: "deck",
  data() {
    return {
      errorMsg: "",
      name: "",
      cards: []
    };
  },
  methods: {},

  async created() {
    try {
      let res = await call({
        path: `/decks/${this.$route.params.uid}`,
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
            name: res.data.references.find(x => x.uid === card).name,
            count: 1
          });
        }
      }
    } catch (e) {
      this.errorMsg = e.response.data.message;
    }
  }
};
</script>

<style scoped>
.deck-card > p {
  text-align: center;
  margin: 0;
  margin-bottom: 5px;
  color: #ccc;
}

.count {
  position: absolute;
  width: 50px;
  height: 50px;
  background: red;
  border-radius: 50px;
  margin-left: 13.2%;
  margin-top: 0.2%;
}

.deck-container {
  margin-top: 20px;
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

.disabled {
  background: #7289da !important;
  opacity: 0.5;
}

.disabled:hover {
  cursor: not-allowed !important;
  background: #7289da !important;
}

.disabled:active {
  background: #7289da !important;
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
  background: #36393f;
  width: 250px;
  border-radius: 4px;
  color: #fff;
  border: 1px solid #666;
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

.wizard select {
  width: 220px;
  margin-top: 4px;
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
