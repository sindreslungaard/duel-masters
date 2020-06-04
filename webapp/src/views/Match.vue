<template>
  <div>
    <div v-show="wait || previewCard || previewCards || errorMessage || warning || action" class="overlay"></div>
    
    <div v-show="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
      <div @click="redirect('overview')" class="btn">Back to overview</div>
    </div>

    <div v-show="warning" class="error warn">
      <p>{{ warning }}</p>
      <div @click="warning = ''" class="btn">Close</div>
    </div>

    <div v-show="wait" class="error">
      <p>{{ wait }}<span class="dots">{{ loadingDots }}</span></p>
    </div>

    <div v-if="previewCard" class="card-preview">
      <img :src="`/assets/cards/all/${previewCard.uid}.jpg`">
      <div @click="dismissLarge()" class="btn">Close</div>
    </div>

    <div v-if="previewCards" class="cards-preview">
      <h1>{{ previewCardsText }}</h1>
      <img  @contextmenu.prevent="previewCards = null; previewCardsText = null; previewCard = card" v-for="(card, index) in previewCards" :key="index" :src="`/assets/cards/all/${card.uid}.jpg`">
      <br><br>
      <div @click="previewCards = null; previewCardsText = null" class="btn">Close</div>
    </div>

    <!-- action (card selection) -->
    <div v-if="action" class="action">
      <span>{{ action.text }}</span>
      <template v-if="actionObject">
        <select class="action-select" v-model="actionDrowdownSelection">
          <option v-for="(option, index) in actionObject" :key="index" :value="index">{{ index }}</option>
        </select>
      </template>
      <div v-if="!actionObject" class="action-cards">
        <div v-for="(card, index) in action.cards" :key="index" class="card"><img @click="actionSelect(card)" :class="[actionSelects.includes(card) ? 'glow-' + card.civilization : '']" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
      </div>
      <div v-if="actionObject" class="action-cards">
        <div v-for="(card, index) in actionObject[actionDrowdownSelection]" :key="index" class="card"><img @click="actionSelect(card)" :class="[actionSelects.includes(card) ? 'glow-' + card.civilization : '']" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
        <p v-if="actionObject[actionDrowdownSelection].length < 1">There's no cards in this category. Use the dropdown above to switch category.</p>
      </div>
      <div @click="chooseAction()" class="btn">Choose</div>
      <div @click="cancelAction()" v-if="action.cancellable" class="btn">Close</div>
      <span style="color: red">{{ actionError }}</span>
    </div>

    <!-- Lobby -->
    <div v-if="decks.length < 1" class="lobby">
      <h1>Waiting for your opponent to join<span class="dots">{{ loadingDots }}</span></h1>
      <div :class="['invite-link', {'copied': inviteCopied}]">
        <span id="invitelink">{{ invite }}</span>
        <div data-clipboard-action="copy" data-clipboard-target="#invitelink" id="invitebtn" :class="['copy', {'copied': inviteCopied}]">{{ inviteCopied ? 'Copied' : 'Copy'}}</div>
      </div>
    </div>

    <!-- Match -->
    <div class="chat">
      <div class="chatbox">

        <div class="messages">
          <div id="messages" class="messages-helper">
            <span v-for="(message, index) in chatMessages" :key="index" v-html="message"></span>
          </div>
        </div>
        <form @submit.prevent="sendChat(chatMessage)">
          <input type="text" v-model="chatMessage" placeholder="Type to chat">
        </form>  
      </div>

      <div class="actionbox handaction">
        <template v-if="handSelection">
          <span>{{ handSelection.name }}</span>
          <div @click="addToPlayzone()" :class="['btn', {'disabled': !handSelection.canBePlayed}]">Add to playzone</div>
          <div class="spacer"></div>
          <div @click="addToManazone()" :class="['btn', {'disabled': state.hasAddedManaThisRound}]">Add to manazone</div>
        </template>
        <template v-if="playzoneSelection">
          <span>{{ playzoneSelection.name }}</span>
          <div @click="attackPlayer()" class="btn">Attack player</div>
          <div class="spacer"></div>
          <div @click="attackCreature()" class="btn">Attack creature</div>
        </template>
      </div>

      <div class="actionbox">
        <div @click="endTurn()" :class="['btn', 'block', { 'disabled': !state.myTurn }]">End turn</div>
      </div>
    </div>
    
    <template v-if="!started">
      <div v-if="deck" class="deck-chooser waiting">
        <h1>Waiting for your opponent to choose a deck<span class="dots">{{ loadingDots }}</span></h1>
      </div>

      <div class="deck-chooser" v-if="decks.length > 0 && !deck">
        <h1>Choose your deck</h1>
        <div class="backdrop">
          <h3>My custom decks</h3>
          <span v-if="decks.filter(x => !x.standard).length < 1">No decks available in this category</span>
          <div @click="chooseDeck(deck.uid)" v-for="(deck, index) in decks.filter(x => !x.standard)" :key="index" class="btn">{{ deck.name }}</div>
        </div>
        

        <br><br>
        <div class="backdrop">
          <h3>Standard decks</h3>
          <span v-if="decks.filter(x => x.standard).length < 1">No decks available in this category</span>
          <div @click="chooseDeck(deck.uid)" v-for="(deck, index) in decks.filter(x => x.standard)" :key="index" class="btn">{{ deck.name }}</div>
        </div>
      </div>
    </template>

    <div v-if="started" class="stadium">
      <div class="stage opponent">

        <div class="manazone">
          <div class="card mana placeholder"><img src="/assets/cards/backside.png"></div>
          <div @contextmenu.prevent="showLarge(card)" v-for="(card, index) in state.opponent.manazone" :key="index" :class="['card', 'mana', { 'tapped': card.tapped }]"><img :src="`/assets/cards/all/${card.uid}.jpg`"></div>
        </div>

        <div class="shieldzone">
          <div class="card shield placeholder"><img src="/assets/cards/backside.png"></div>
          <div v-for="(card, index) in state.opponent.shieldzone" :key="index" class="card shield flipped"><img src="/assets/cards/backside.png"></div>
        </div>

        <div class="playzone">
          <div class="card placeholder"><img src="/assets/cards/backside.png"></div>
          <div @contextmenu.prevent="showLarge(card)" v-for="(card, index) in state.opponent.playzone" :key="index" :class="['card', { 'tapped': card.tapped }]"><img class="flipped" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
        </div>

      </div>

      <div class="right-stage">
        <div class="right-stage-content">
          <p>Graveyard [{{ state.opponent.graveyard.length }}]</p>
          <div class="card">
            <img @contextmenu.prevent="" v-if="state.opponent.graveyard.length < 1" style="height: 10vh; opacity: 0.3" src="/assets/cards/backside.png">
            <img @contextmenu.prevent="previewCards = state.opponent.graveyard; previewCardsText = 'Opponent\'s Graveyard'" v-if="state.opponent.graveyard.length > 0" style="height: 10vh" :src="`/assets/cards/all/${state.opponent.graveyard[0].uid}.jpg`">
          </div>

          <p>Deck [{{ state.opponent.deck }}]</p>
          <div class="card"><img @contextmenu.prevent="" style="height: 10vh" src="/assets/cards/backside.png"></div>
        </div>       
      </div>

      <div class="right-stage bt">
        <div class="right-stage-content">
          <p>Graveyard [{{ state.me.graveyard.length }}]</p>
          <div class="card">
            <img @contextmenu.prevent="" v-if="state.me.graveyard.length < 1" style="height: 10vh; opacity: 0.3" src="/assets/cards/backside.png">
            <img @contextmenu.prevent="previewCards = state.me.graveyard; previewCardsText = 'My Graveyard'" v-if="state.me.graveyard.length > 0" style="height: 10vh" :src="`/assets/cards/all/${state.me.graveyard[0].uid}.jpg`">
          </div>

          <p>Deck [{{ state.me.deck }}]</p>
          <div class="card"><img @contextmenu.prevent="" style="height: 10vh" src="/assets/cards/backside.png"></div>
        </div>
      </div>

      <div class="stage me bt">

        <div class="playzone">
          <div class="card placeholder"><img src="/assets/cards/backside.png"></div>
          <div @click="onPlayzoneClicked(card)" @contextmenu.prevent="showLarge(card)" v-for="(card, index) in state.me.playzone" :key="index" :class="['card', { 'tapped': card.tapped }]"><img :class="playzoneSelection === card ? 'glow-' + card.civilization : ''" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
        </div>

        <div class="shieldzone">
          <div class="card shield placeholder"><img src="/assets/cards/backside.png"></div>
          <div v-for="(card, index) in state.me.shieldzone" :key="index" class="card shield"><img src="/assets/cards/backside.png"></div>
        </div>

        <div class="manazone">
          <div class="card mana placeholder"><img src="/assets/cards/backside.png"></div>
          <div @contextmenu.prevent="showLarge(card)" v-for="(card, index) in state.me.manazone" :key="index" :class="['card', 'mana', { 'tapped': card.tapped }]"><img class="flipped" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
        </div>
      </div>    

      <div class="hand bt">
        <div class="card placeholder"><img src="/assets/cards/backside.png"></div>
        <div @contextmenu.prevent="showLarge(card)" @click="makeHandSelection(card)" class="card" v-for="(card, index) in state.me.hand" :key="index"><img :class="[handSelection === card ? 'glow-' + card.civilization : '']" :src="`/assets/cards/all/${card.uid}.jpg`"></div>
      </div>
    </div>

  </div>
</template>

<script>

import config from '../config'
import ClipboardJS from "clipboard"
import { ws_protocol } from "../remote"

const send = (client, message) => {
  client.send(JSON.stringify(message))
}

export default {
  name: "game",
  data() {
    return {
      ws: null,

      errorMessage: "",
      warning: "",
      wait: "",

      loadingDots: "",
      invite: window.location.host + "/duel/" + this.$route.params.id,
      inviteCopied: false,
      inviteCopyTask: null,

      chatMessages: [],
      chatMessage: "",

      started: false,

      opponent: "",
      decks:[],
      deck: null,

      state: {},
      handSelection: null,

      playzoneSelection: null,

      action: null,
      actionError: "",
      actionSelects: [],
      actionObject: null,
      actionDrowdownSelection: null,

      previewCard: null,
      previewCards: null,
      previewCardsText: null

    }
  },
  methods: {
    redirect(to) {
      this.$router.push('/' + to)
    },
    sendChat(message) {
      if(!message) {
        return
      }
      this.chatMessage = ""
      this.ws.send(JSON.stringify({ header: "chat", message }))
    },
    chat(message) {
      this.chatMessages.push(message)
      this.$nextTick(() => {
        let container = document.getElementById('messages')
        container.scrollTop = container.scrollHeight
      })
    },

    chooseDeck(uid) {
      this.deck = uid
      this.ws.send(JSON.stringify({ header: "choose_deck", uid }))
    },

    makeHandSelection(card) {
      if(!this.state.myTurn) {
        return
      }
      this.playzoneSelection = null
      if(this.handSelection === card) {
        this.handSelection = null
        return
      }
      this.handSelection = card
    },

    actionSelect(card) {
      if(this.actionSelects.includes(card)) {
        this.actionSelects = this.actionSelects.filter(x => x !== card)
        return
      }

      if(this.actionSelects.length >= this.action.maxSelections) {
        return
      }

      this.actionSelects.push(card)
    },

    cancelAction() {
      if(!this.action || !this.action.cancellable) {
        return
      }
      this.ws.send(JSON.stringify({ header: "action",  cancel: true }))
    },

    chooseAction() {
      if(!this.action) {
        return
      }
      let cards = []
      for(let card of this.actionSelects) {
        cards.push(card.virtualId)
      }
      this.ws.send(JSON.stringify({ header: "action", cards, cancel: false }))
    },

    addToManazone() {
      if(!this.handSelection) {
        return
      }
      this.ws.send(JSON.stringify({ header: "add_to_manazone", virtualId: this.handSelection.virtualId }))
    },

    addToPlayzone() {
      if(!this.handSelection) {
        return
      }
      this.ws.send(JSON.stringify({ header: "add_to_playzone", virtualId: this.handSelection.virtualId }))
    },

    endTurn() {
      if(!this.state.myTurn) {
        return
      }
      this.ws.send(JSON.stringify({ header: "end_turn" }))
    },

    showLarge(card) {
      this.previewCard = card
    },

    dismissLarge() {
      this.previewCard = null
    },

    onPlayzoneClicked(card) {
      if(!this.state.myTurn) {
        return
      }
      this.handSelection = null
      if(this.playzoneSelection && this.playzoneSelection === card) {
        this.playzoneSelection = null
        return
      }
      if(card.tapped) {
        return
      }
      this.playzoneSelection = card
    },

    attackPlayer() {
      if(!this.playzoneSelection) {
        return
      }
      this.ws.send(JSON.stringify({ header: "attack_player", virtualId: this.playzoneSelection.virtualId }))
    },
    attackCreature() {
      if(!this.playzoneSelection) {
        return
      }
      this.ws.send(JSON.stringify({ header: "attack_creature", virtualId: this.playzoneSelection.virtualId }))
    },

  },
  created() {

    // Connect to the server
    const ws = new WebSocket(remote.ws_protocol + window.location.host + "/ws/" + this.$route.params.id)
    this.ws = ws

    ws.onopen = () => {
      ws.send(localStorage.getItem("token"))
    }

    ws.onmessage = (event) => {
      
      const data = JSON.parse(event.data)

      switch(data.header) {

        case "mping": {
          send(ws, {
            header: "mpong"
          })
          break
        }
        
        case "hello": {
          send(ws, {
            header: "join_match"
          })
          break
        }

        case "error": {
          this.errorMessage = data.message
          break
        }

        case "warn": {
          this.warning = data.message
          break
        }

        case "player_joined": {
          this.opponent = data.username
          break
        }

        case "choose_deck": {
          this.decks = data.decks
          break
        }

        case "chat": {
          this.chat(`<span style="color: ${data.color}">[${data.sender}]</span> <span>${data.message}</span>`)
          break
        }

        case "state_update": {
          if(!this.started) {
            this.started = true
          }
          this.handSelection = null
          this.playzoneSelection = null
          this.state = data.state
          break
        }

        case "action": {
          this.actionError = "",
          this.actionSelects = []
          if(!(data.cards instanceof Array)) {
            this.actionObject = data.cards
            console.log(Object.keys(data.cards)[0])
            this.actionDrowdownSelection = Object.keys(data.cards)[0]
          }
          this.action = {
            cards: data.cards instanceof Array ? data.cards : Object.keys(data.cards)[0],
            text: data.text,
            minSelection: data.minSelection,
            maxSelections: data.maxSelections,
            cancellable: data.cancellable
          }
          break
        }

        case "action_error": {
          if(!this.action) {
            return
          }
          this.actionError = data.message
          break
        }

        case "close_action": {
          this.action = null
          this.actionError = ""
          this.actionSelects = []
          this.actionObject = null
          this.actionDrowdownSelection = null
          break
        }

        case "wait": {
          this.wait = data.message || "Waiting for your opponent to make an action"
          break
        }

        case "end_wait": {
          this.wait = ""
          break
        }

      }

    }

    ws.onclose = () => { 
      this.errorMessage = "Connection to the server has been closed."
      console.log("connection closed")
    }

    ws.onerror = (event) => {
      console.log(event)
      this.errorMessage = "An error occured when attempting to communicate with the server."
    }


    // Loading dots
    setInterval(() => {
      if(this.loadingDots.length >= 4)
        this.loadingDots = ""
      else this.loadingDots += "."
    }, 500)

    // clipboard
    let clipboard = new ClipboardJS('#invitebtn')
    clipboard.on('success', (e) => {
      if(this.inviteCopyTask)
        clearTimeout(this.inviteCopyTask)
      this.inviteCopied = true
      this.inviteCopyTask = setTimeout(() => {
        this.inviteCopied = false
      }, 2000)
      e.clearSelection()
    })

  },
  beforeDestroy() {
    this.ws.close()
  }
}
</script>

<style scoped lang="scss">

.card-preview {
  width: 300px;
  text-align: center;
  border-radius: 4px;
  height: 480px;
  z-index: 2005;
  position: absolute;
  left: calc(50% - 300px / 2);
  top: calc(50vh - 480px / 2);
  img {
    width: 300px;
    border-radius: 15px;
    margin-bottom: 10px;
  }
}

.action-select {
    border: none;
    background: #484C52;
    padding: 5px !important;
    width: auto !important;
    margin-left: 5px;
    border-radius: 4px;
    color: #ccc;
    resize: none;
}
.action-select:focus {
    outline: none
}
.action-select:active {
    outline: none
}

.action-select {
  margin-top: 10px;
}

.action {
  max-height: 425px;
  width: 790px;
  background: #2F3136;
  position: absolute;
  z-index: 3000;
  margin: 0 auto;
  left: calc(50% - 790px / 2);
  top: calc(50vh - 300px / 2);
  text-align: center;
  border-radius: 4px;
  border: 1px solid #666;
  overflow-x: auto;
  padding-bottom: 15px;
  span {
    color: #ccc;
    font-size: 13px;
    display: block;
    margin: 0 30px;
    margin-top: 15px;
  }
  .btn {
    margin: 0 7px;
  }
  .action-cards {
    background: #222428;
    margin: 15px;
    border-radius: 4px;
    padding: 10px;
    max-height: 300px;
    overflow: auto;
    img {
      height: 125px;
    }
    .card {
      margin: 0 7px;
    }
  }
}

.placeholder {
  width: 0 !important;
  margin-left: 0 !important;
  margin-right: 0 !important;
  padding-left: 0 !important;
  padding-right: 0 !important;
  opacity: 0;
  img {
    width: 0
  }
}

.glow {
  box-shadow: 0px 0px 4px 0px red;
}

.glow-fire {
  box-shadow: 0px 0px 4px 0px red;
}
.glow-water {
  box-shadow: 0px 0px 4px 0px cyan;
}
.glow-nature {
  box-shadow: 0px 0px 4px 0px green;
}
.glow-light {
  box-shadow: 0px 0px 4px 0px yellow;
}
.glow-darkness {
  box-shadow: 0px 0px 4px 0px black;
}

.waiting {
  h1 {
    display: inline-block;
  }
  span {
    display: inline-block !important;
    font-size: 26px !important;
    line-height: 0;
  }
  display: inline-block;
}

.deck-chooser {
  overflow: auto;
  padding: 15px 50px;
  h3 {
    margin: 0;
    color: #ccc;
  }
  .btn {
    margin-top: 10px;
    margin-right: 10px;
  }
  span {
    display: block;
    margin-top: 10px;
    color: #ccc;
    font-size: 14px;
  }
}

.warn {
  z-index: 50000;
}

.backdrop {
  background: #2F3136;
  padding: 10px;
  border-radius: 4px;
}

.block {
  display: block !important;
}

.chatbox {
  height: calc(100vh - 128px - 15px);
  background: #2F3136;
  margin: 5px;
  border-radius: 4px;
}

.handaction {
  height: 53px !important;
  span {
    display: block;
    font-size: 13px;
    margin-bottom: 7px;
  }
  .btn {
    width: 110px;
  }
  .spacer {
    display: inline-block;
    width: 10px;
  }
}

.btn:active {
  background: #5B6EAE;
}

.disabled {
    background: #7289DA !important;
    opacity: 0.5;
}

.disabled:hover {
    cursor: not-allowed !important;
    background: #7289DA !important;
}

.disabled:active {
    background: #7289DA !important;
}

.messages {
  height: calc(100% - 77px);
  padding: 10px;
  font-size: 14px;
  color: #ccc;
  position: relative;
}

.messages-helper {
  position: absolute;
  bottom: 0;
  overflow: auto;
  max-height: calc(100% - 10px);
  width: calc(100% - 20px);
  > span {
    display: block;
    margin-top: 7px;
  }
}

*::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
    border-radius: 10px;
    background-color: #484C52;
  }

*::-webkit-scrollbar {
  width: 6px;
  height: 6px;
  background-color: #484C52;
}

*::-webkit-scrollbar-thumb {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    box-shadow: inset 0 0 6px rgba(0,0,0,.3);
    background-color: #222;
}

.chatbox input {
  border: none;
  border-radius: 4px;
  margin: 10px;
  width: calc(100% - 40px);
  background: #484C52;
  padding: 10px;
  color: #ccc;
  &:focus {
    outline: none
  }
  &:active {
    outline: none
  }
}

.actionbox {
  background:#2F3136;
  height: 30px;
  margin: 5px;
  padding: 10px;
  border-radius: 4px;
}

.lobby {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: #36393F;
  z-index: 10;
  text-align: center;
  padding-top: 10vh;
}

.invite-link {
  padding: 5px;
  padding-left: 10px;
  background: #2B2E33;
  border: 1px solid #222428;
  border-radius: 4px;
  display: inline-block;
  color: #E3E3E5;
  transition: 0.1s;
}

.invite-link span {
  float: left;
  display: block;
  margin-top: 5px;
}

.invite-link .copy {
  display: inline-block;
  background: #7289DA;
  color: #E3E3E5;
  font-size: 14px;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  margin-left: 20px;
  transition: 0.1s;
  text-align: center !important;
  width: 45px;
}

.copy:hover {
  cursor: pointer;
  background: #677BC4;
}

.copied {
  border-color: #3CA374 !important;
  color: #fff !important;
}

.invite-link > .copied {
  background: #3CA374;
}

.dots {
  display: inline-block;
  width: 40px;
  text-align: left;
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

.error p {
  padding: 5px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: #2B2E33 !important;
  border: 1px solid #222428;
}

.btn:hover {
  cursor: pointer;
  background: #677BC4;
}

.error {
  border: 1px solid #666;
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #36393F;
  z-index: 3005;
  left: calc(50% - 300px / 2);
  top: 40vh;
  padding: 10px;
  font-size: 14px;
  color: #ccc;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: #000;
  opacity: 0.5;
  z-index: 100;
}

.chat {
  width: 300px;
  height: 100vh;
  border-right: 1px solid #555;
  float: left;
}

.stage {
  width: calc(93% - 301px);
  height: 41vh;
  float: left;
}

.right-stage {
  width: 7%;
  height: 41vh;
  float: right;
}

.right-stage-content {
  text-align: center;
  margin-top: 7vh;
}

.right-stage p {
  color: #ccc;
  font-size: 14px;
  margin-bottom: 3px;
}

.bt {
  border-top: 1px solid #555;
}

.hand {
  width: calc(100% - 301px);
  float: left;
}

.card {
  display: inline-block;
  margin: 0.8%;
  margin-bottom: 0;
  img {
    height: 18vh;
    border-radius: 8px;
  }
}

.flipped {
  transform: rotate(180deg) scaleX(-1);
}

.tapped {
  transform: rotate(90deg);
  margin-left: 25px;
  margin-right: 1%;
}

.playzone {
  overflow: auto;
  white-space: nowrap;
  overflow-y: hidden;
  height: 20vh;
}

.playzone .tapped {
  margin-left: 35px;
  margin-right: 35px;
}

.shield {
  img {
    height: 8.5vh;
  }
}

.shieldzone {
  overflow: auto;
  white-space: nowrap;
  overflow-y: hidden;
  height: 10vh;
}

.mana {
  img {
    height: 8.5vh;
  }
}

.manazone {
  overflow: auto;
  white-space: nowrap;
  overflow-y: hidden;
  height: 10vh;
}

.hand {
  img {
    height: 15vh;
  }
  overflow: auto;
  white-space: nowrap;
  overflow-y: hidden;
  height: 17vh;
}

.cards-preview {
  position: absolute;
  text-align: center;
  width: 80%;
  left: 10%;
  top: 25vh;
  z-index: 700;
}

.cards-preview img {
  height: 20vh;
  display: inline-block;
  border-radius: 7px;
  margin: 10px;
}

</style>
