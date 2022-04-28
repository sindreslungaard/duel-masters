<template>
  <div>
    <div
      v-show="
        wait ||
          previewCard ||
          previewCards ||
          errorMessage ||
          warning ||
          action ||
          opponentDisconnected ||
          reconnecting
      "
      @click="handleOverlayClick()"
      class="overlay"
    ></div>

    <div v-if="opponentDisconnected && !errorMessage" class="error">
      <p>
        Your opponent disconnected or left the match. Waiting for them to
        reconnect{{ loadingDots }}
      </p>
      <div @click="redirect('overview')" class="btn">Leave duel</div>
    </div>

    <div v-if="reconnecting && !errorMessage" class="error">
      <p>Disconnected from server. Attempting to reconnect{{ loadingDots }}</p>
      <div @click="redirect('overview')" class="btn">Leave duel</div>
    </div>

    <div v-show="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
      <div @click="redirect('overview')" class="btn">Back to overview</div>
    </div>

    <div v-show="warning" class="error warn">
      <p>{{ warning }}</p>
      <div @click="warning = ''" class="btn">Close</div>
    </div>

    <div v-show="wait" class="error">
      <p>
        {{ wait }}<span class="dots">{{ loadingDots }}</span>
      </p>
    </div>

    <div v-if="previewCard" class="card-preview">
      <img :src="`/assets/cards/all/${previewCard.uid}.jpg`" />
      <div @click="dismissLarge()" class="btn">Close</div>
    </div>

    <div v-if="previewCards" class="cards-preview" @click="dismissLarge()">
      <h1>{{ previewCardsText }}</h1>
      <img
        @contextmenu.prevent="
          dismissLarge();
          showLarge(card);
        "
        v-for="(card, index) in previewCards"
        :key="index"
        :src="`/assets/cards/all/${card.uid}.jpg`"
      />
      <br /><br />
      <div @click="dismissLarge()" class="btn">
        Close
      </div>
    </div>

    <!-- action (card selection) -->
    <div v-if="action" id="action" class="action noselect">
      <span v-draggable data-ref="action">{{ action.text }}</span>
      <span><i> Tip: click and drag to (de)select faster</i></span>
      <template v-if="actionObject">
        <select class="action-select" v-model="actionDrowdownSelection">
          <option
            v-for="(option, index) in actionObject"
            :key="index"
            :value="index"
            >{{ index }}</option
          >
        </select>
      </template>
      <div v-if="!actionObject" class="action-cards">
        <div v-for="(card, index) in action.cards" :key="index" class="card">
          <img
            @dragstart.prevent=""
            @mouseenter="actionSelectMouseEnter($event, card)"
            @mousedown="actionSelect(card)"
            :class="[
              'no-drag',
              actionSelects.includes(card) ? 'glow-' + card.civilization : ''
            ]"
            :src="`/assets/cards/all/${card.uid}.jpg`"
          />
        </div>
      </div>
      <div v-if="actionObject" class="action-cards">
        <div
          v-for="(card, index) in actionObject[actionDrowdownSelection]"
          :key="index"
          class="card"
        >
          <img
            @dragstart.prevent=""
            @mouseenter="actionSelectMouseEnter($event, card)"
            @mousedown="actionSelect(card)"
            :class="[
              'no-drag',
              actionSelects.includes(card) ? 'glow-' + card.civilization : ''
            ]"
            :src="`/assets/cards/all/${card.uid}.jpg`"
          />
        </div>
        <p v-if="actionObject[actionDrowdownSelection].length < 1">
          There's no cards in this category. Use the dropdown above to switch
          category.
        </p>
      </div>
      <div @click="chooseAction()" class="btn">Choose</div>
      <div @click="cancelAction()" v-if="action.cancellable" class="btn">
        Close
      </div>
      <span style="color: red">{{ actionError }}</span>
    </div>

    <!-- Lobby -->
    <div v-if="!started && decks.length < 1" class="lobby">
      <h1>
        Waiting for your opponent to join<span class="dots">{{
          loadingDots
        }}</span>
      </h1>
      <div :class="['invite-link', { copied: inviteCopied }]">
        <span id="invitelink">{{ invite }}</span>
        <div
          data-clipboard-action="copy"
          data-clipboard-target="#invitelink"
          id="invitebtn"
          :class="['copy', { copied: inviteCopied }]"
        >
          {{ inviteCopied ? "Copied" : "Copy" }}
        </div>
      </div>
    </div>

    <!-- Match -->
    <div class="chat">
      <div :class="state.spectator ? 'fullsize-chatbox' : 'chatbox'">
        <div class="messages">
          <div id="messages" class="messages-helper">
            <div
              class="message"
              :style="{
                background:
                  message.sender.toLowerCase() === 'server' ? 'none' : '#202124'
              }"
              v-for="(message, index) in chatMessages.filter(
                m => !settings.muted.includes(m.sender)
              )"
              :key="index"
            >
              <div
                class="message-sender"
                :style="{ color: message.color || 'orange' }"
              >
                {{
                  message.sender.toLowerCase() == "server"
                    ? "-"
                    : message.sender + ":"
                }}
              </div>
              <div class="message-text">{{ message.message }}</div>
              <div class="mute-icon-container">
                <MuteIcon
                  v-if="
                    !['server', username].includes(message.sender.toLowerCase())
                  "
                  :player="message.sender"
                  @toggled="onSettingsChanged()"
                />
              </div>
            </div>
          </div>
        </div>
        <form @submit.prevent="sendChat(chatMessage)">
          <input type="text" v-model="chatMessage" placeholder="Type to chat" />
        </form>
      </div>

      <div v-if="!state.spectator" class="actionbox handaction">
        <template v-if="handSelection">
          <span>{{ handSelection.name }}</span>
          <div
            @click="addToPlayzone()"
            :class="['btn', { disabled: !handSelection.canBePlayed }]"
          >
            Add to playzone
          </div>
          <div class="spacer"></div>
          <div
            @click="addToManazone()"
            :class="['btn', { disabled: state.hasAddedManaThisRound }]"
          >
            Add to manazone
          </div>
        </template>
        <template v-if="playzoneSelection">
          <span>{{ playzoneSelection.name }}</span>
          <div @click="attackPlayer()" class="btn">Attack player</div>
          <div class="spacer"></div>
          <div @click="attackCreature()" class="btn">Attack creature</div>
        </template>
      </div>

      <div v-if="!state.spectator" class="actionbox">
        <div
          @click="endTurn()"
          :class="['btn', 'block', { disabled: !state.myTurn }]"
        >
          End turn
        </div>
      </div>
    </div>

    <template v-if="!started">
      <div v-if="deck" class="deck-chooser waiting">
        <h1>
          Waiting for your opponent to choose a deck<span class="dots">{{
            loadingDots
          }}</span>
        </h1>
      </div>

      <div class="deck-chooser" v-if="decks.length > 0 && !deck">
        <h1>Choose your deck</h1>
        <div class="backdrop">
          <h3>My custom decks</h3>
          <span v-if="decks.filter(x => !x.standard).length < 1"
            >No decks available in this category</span
          >
          <div
            @click="chooseDeck(deck.uid)"
            v-for="(deck, index) in decks.filter(x => !x.standard)"
            :key="index"
            class="btn"
          >
            {{ deck.name }}
          </div>
        </div>

        <br /><br />
        <div class="backdrop">
          <h3>Standard decks</h3>
          <span v-if="decks.filter(x => x.standard).length < 1"
            >No decks available in this category</span
          >
          <div
            @click="chooseDeck(deck.uid)"
            v-for="(deck, index) in decks.filter(x => x.standard)"
            :key="index"
            class="btn"
          >
            {{ deck.name }}
          </div>
        </div>
      </div>
    </template>

    <div v-if="started" class="stadium">
      <div class="stage opponent">
        <div class="manazone">
          <div class="card mana placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            @contextmenu.prevent="showLarge(card)"
            v-for="(card, index) in state.opponent.manazone"
            :key="index"
            :class="['card', 'mana', { tapped: card.tapped }]"
          >
            <img :src="`/assets/cards/all/${card.uid}.jpg`" />
          </div>
        </div>

        <div
          class="shieldzone"
          @drop.prevent="drop($event, 'opponentshieldzone')"
          @dragover.prevent
          @dragenter.prevent
          ref="opponentshieldzone"
        >
          <div class="card shield placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            v-for="(card, index) in state.opponent.shieldzone"
            :key="index"
            :class="
              'card shield' + (!settings.noUpsideDownCards ? ' flipped' : '')
            "
          >
            <img src="/assets/cards/backside.png" />
          </div>
        </div>

        <div
          class="playzone"
          @drop.prevent="drop($event, 'opponentsplayzone')"
          @dragover.prevent
          @dragenter.prevent
          ref="opponentsplayzone"
        >
          <div class="card placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            @contextmenu.prevent="showLarge(card)"
            v-for="(card, index) in state.opponent.playzone"
            :key="index"
            :class="['card', { tapped: card.tapped }]"
          >
            <img
              :class="!settings.noUpsideDownCards ? 'flipped' : ''"
              :src="`/assets/cards/all/${card.uid}.jpg`"
            />
          </div>
        </div>
      </div>

      <div class="right-stage">
        <div class="right-stage-content">
          <p>Hand [{{ state.opponent.handCount }}]</p>
          <p>Graveyard [{{ state.opponent.graveyard.length }}]</p>
          <div class="card">
            <img
              @contextmenu.prevent=""
              v-if="state.opponent.graveyard.length < 1"
              style="height: 10vh; opacity: 0.3"
              src="/assets/cards/backside.png"
            />
            <img
              @contextmenu.prevent="
                previewCards = state.opponent.graveyard;
                previewCardsText =
                  (state.spectator ? state.opponent.username : 'Opponent') +
                  '\'s Graveyard';
              "
              v-if="state.opponent.graveyard.length > 0"
              style="height: 10vh"
              :src="`/assets/cards/all/${state.opponent.graveyard[0].uid}.jpg`"
            />
          </div>

          <p>Deck [{{ state.opponent.deck }}]</p>
          <div class="card">
            <img
              @contextmenu.prevent=""
              style="height: 10vh"
              src="/assets/cards/backside.png"
            />
          </div>
        </div>
      </div>

      <div class="right-stage bt">
        <div class="right-stage-content">
          <p v-if="state.spectator">Hand [{{ state.me.handCount }}]</p>
          <p>Graveyard [{{ state.me.graveyard.length }}]</p>
          <div class="card">
            <img
              @contextmenu.prevent=""
              v-if="state.me.graveyard.length < 1"
              style="height: 10vh; opacity: 0.3"
              src="/assets/cards/backside.png"
            />
            <img
              @contextmenu.prevent="
                previewCards = state.me.graveyard;
                previewCardsText =
                  (state.spectator ? state.me.username + '\'s' : 'My') +
                  ' Graveyard';
              "
              v-if="state.me.graveyard.length > 0"
              style="height: 10vh"
              :src="`/assets/cards/all/${state.me.graveyard[0].uid}.jpg`"
            />
          </div>

          <p>Deck [{{ state.me.deck }}]</p>
          <div class="card">
            <img
              @contextmenu.prevent=""
              style="height: 10vh"
              src="/assets/cards/backside.png"
            />
          </div>
        </div>
      </div>

      <div class="stage me bt">
        <div
          @drop.prevent="drop($event, 'playzone')"
          @dragover.prevent
          @dragenter.prevent
          ref="myplayzone"
          class="playzone"
        >
          <div class="card placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            @click="onPlayzoneClicked(card)"
            @contextmenu.prevent="showLarge(card)"
            v-for="(card, index) in state.me.playzone"
            :key="index"
            :class="['card', { tapped: card.tapped }]"
          >
            <img
              draggable
              @dragstart="startDrag($event, card, 'playzone')"
              @dragend="stopDrag('playzone')"
              :class="
                playzoneSelection === card ? 'glow-' + card.civilization : ''
              "
              :src="`/assets/cards/all/${card.uid}.jpg`"
            />
          </div>
        </div>

        <div class="shieldzone">
          <div class="card shield placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            v-for="(card, index) in state.me.shieldzone"
            :key="index"
            class="card shield"
          >
            <img src="/assets/cards/backside.png" />
          </div>
        </div>

        <div
          @drop.prevent="drop($event, 'manazone')"
          @dragover.prevent
          @dragenter.prevent
          ref="mymanazone"
          class="manazone"
        >
          <div class="card mana placeholder">
            <img src="/assets/cards/backside.png" />
          </div>
          <div
            @contextmenu.prevent="showLarge(card)"
            v-for="(card, index) in state.me.manazone"
            :key="index"
            :class="['card', 'mana', { tapped: card.tapped }]"
          >
            <img
              :class="!settings.noUpsideDownCards ? 'flipped' : ''"
              :src="`/assets/cards/all/${card.uid}.jpg`"
            />
          </div>
        </div>
      </div>

      <div class="hand bt">
        <div class="spectator-info" v-if="state.spectator">
          <div>You are spectating</div>
          <Username :color="state.me.color">{{ state.me.username }}</Username>
          <div>vs</div>
          <Username :color="state.opponent.color">{{
            state.opponent.username
          }}</Username>
        </div>
        <div class="spectator-leave" v-if="state.spectator">
          <span @click="$router.push('/overview')">Stop spectating</span>
        </div>

        <div class="card placeholder">
          <img src="/assets/cards/backside.png" />
        </div>
        <div
          @contextmenu.prevent="showLarge(card)"
          @click="makeHandSelection(card)"
          class="card"
          v-for="(card, index) in state.me.hand"
          :key="index"
        >
          <img
            draggable
            @dragstart="startDrag($event, card, 'hand')"
            @dragend="stopDrag('hand')"
            :class="[handSelection === card ? 'glow-' + card.civilization : '']"
            :src="`/assets/cards/all/${card.uid}.jpg`"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ClipboardJS from "clipboard";
import { call, ws_protocol, host } from "../remote";
import CardShowDialog from "../components/dialogs/CardShowDialog";
import Username from "../components/Username.vue";
import MuteIcon from "../components/MuteIcon.vue";
import { getSettings, didSeeMuteWarning } from "../helpers/settings";

const send = (client, message) => {
  client.send(JSON.stringify(message));
};

function sound(src) {
  this.sound = document.createElement("audio");
  this.sound.src = src;
  this.sound.setAttribute("preload", "auto");
  this.sound.setAttribute("controls", "none");
  this.sound.style.display = "none";
  this.sound.volume = 0.3;
  document.body.appendChild(this.sound);
  this.play = function() {
    this.sound.play();
  };
  this.stop = function() {
    this.sound.pause();
  };
}

let turnSound = new sound("/assets/turn.mp3");
let playerJoinedSound = new sound("/assets/player_joined.mp3");

export default {
  name: "game",
  components: {
    Username,
    MuteIcon
  },
  data() {
    return {
      ws: null,
      reconnecting: false,
      preventReconnect: false,

      errorMessage: "",
      warning: "",
      wait: "",

      loadingDots: "",
      invite:
        location.protocol + "//" + host + "/invite/" + this.$route.params.id,
      inviteCopied: false,
      inviteCopyTask: null,

      chatMessages: [], // { sender, message, color }
      chatMessage: "",

      started: false,

      opponent: "",
      opponentDisconnected: false,
      decks: [],
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
      previewCardsText: null,
      settings: getSettings()
    };
  },
  computed: {
    username: () => localStorage.getItem("username")
  },
  methods: {
    redirect(to) {
      this.$router.push("/" + to);
    },
    sendChat(message) {
      if (!message) {
        return;
      }
      this.chatMessage = "";
      this.ws.send(JSON.stringify({ header: "chat", message }));
    },

    chat(sender, color, message) {
      this.chatMessages.push({ sender, color, message });
      this.$nextTick(() => {
        let container = document.getElementById("messages");
        container.scrollTop = container.scrollHeight;
      });
    },

    onSettingsChanged(e) {
      if (!e && !didSeeMuteWarning()) {
        this.warning =
          "You can unmute players at any time from the settings page";
      }

      this.settings = getSettings();
    },

    chooseDeck(uid) {
      this.deck = uid;
      this.ws.send(JSON.stringify({ header: "choose_deck", uid }));
    },

    handleOverlayClick() {
      if (this.previewCard || this.previewCards) {
        this.dismissLarge();
      }
    },

    makeHandSelection(card) {
      if (!this.state.myTurn) {
        return;
      }
      this.playzoneSelection = null;
      if (this.handSelection === card) {
        this.handSelection = null;
        return;
      }
      this.handSelection = card;
    },

    actionSelect(card) {
      if (this.actionSelects.includes(card)) {
        this.actionSelects = this.actionSelects.filter(x => x !== card);
        return;
      }

      if (this.actionSelects.length >= this.action.maxSelections) {
        return;
      }

      this.actionSelects.push(card);
    },

    actionSelectMouseEnter(event, card) {
      const isLeftClick = event.buttons === 1;

      if (isLeftClick) {
        this.actionSelect(card);
      }
    },

    cancelAction() {
      if (!this.action || !this.action.cancellable) {
        return;
      }
      this.ws.send(JSON.stringify({ header: "action", cancel: true }));
    },

    chooseAction() {
      if (!this.action) {
        return;
      }
      let cards = [];
      for (let card of this.actionSelects) {
        cards.push(card.virtualId);
      }
      this.ws.send(JSON.stringify({ header: "action", cards, cancel: false }));
    },

    addToManazone() {
      if (!this.handSelection) {
        return;
      }
      this.ws.send(
        JSON.stringify({
          header: "add_to_manazone",
          virtualId: this.handSelection.virtualId
        })
      );
    },

    addToPlayzone() {
      if (!this.handSelection) {
        return;
      }
      this.ws.send(
        JSON.stringify({
          header: "add_to_playzone",
          virtualId: this.handSelection.virtualId
        })
      );
    },

    endTurn() {
      if (!this.state.myTurn) {
        return;
      }
      this.ws.send(JSON.stringify({ header: "end_turn" }));
    },

    showLarge(card) {
      this.previewCard = card;
    },

    dismissLarge() {
      this.previewCard = null;
      this.previewCards = null;
      this.previewCardsText = null;
    },

    onPlayzoneClicked(card) {
      if (!this.state.myTurn) {
        return;
      }
      this.handSelection = null;
      if (this.playzoneSelection && this.playzoneSelection === card) {
        this.playzoneSelection = null;
        return;
      }
      if (card.tapped) {
        return;
      }
      this.playzoneSelection = card;
    },

    attackPlayer() {
      if (!this.playzoneSelection) {
        return;
      }
      this.ws.send(
        JSON.stringify({
          header: "attack_player",
          virtualId: this.playzoneSelection.virtualId
        })
      );
    },
    attackCreature() {
      if (!this.playzoneSelection) {
        return;
      }
      this.ws.send(
        JSON.stringify({
          header: "attack_creature",
          virtualId: this.playzoneSelection.virtualId
        })
      );
    },
    startDrag(evt, card, source) {
      evt.dataTransfer.dropEffect = "move";
      evt.dataTransfer.effectAllowed = "move";
      evt.dataTransfer.setData("vid", card.virtualId);

      const greenHighlight = "#507053";
      const redHighlight = "#7d5252";

      if (source === "hand") {
        if (card.canBePlayed) {
          this.$refs.myplayzone.style.backgroundColor = greenHighlight;
        } else {
          this.$refs.myplayzone.style.backgroundColor = redHighlight;
        }

        if (!this.state.hasAddedManaThisRound) {
          this.$refs.mymanazone.style.backgroundColor = greenHighlight;
        } else {
          this.$refs.mymanazone.style.backgroundColor = redHighlight;
        }
      } else if (source === "playzone") {
        if (this.state.opponent.playzone.length) {
          this.$refs.opponentsplayzone.style.backgroundColor = greenHighlight;
        } else {
          this.$refs.opponentsplayzone.style.backgroundColor = redHighlight;
        }

        this.$refs.opponentshieldzone.style.backgroundColor = greenHighlight;
      }
    },
    stopDrag(source) {
      if (source === "hand") {
        this.$refs.myplayzone.style.backgroundColor = "transparent";
        this.$refs.mymanazone.style.backgroundColor = "transparent";
      } else if (source === "playzone") {
        this.$refs.opponentsplayzone.style.backgroundColor = "transparent";
        this.$refs.opponentshieldzone.style.backgroundColor = "transparent";
      }
    },
    drop(event, zone) {
      const vid = event.dataTransfer.getData("vid");

      if (["manazone", "playzone"].includes(zone)) {
        this.handSelection = this.state.me.hand.find(x => x.virtualId === vid);

        if (zone === "manazone") {
          this.addToManazone();
        } else if (zone === "playzone") {
          this.addToPlayzone();
        }
      }

      if (["opponentshieldzone", "opponentsplayzone"].includes(zone)) {
        this.playzoneSelection = this.state.me.playzone.find(
          x => x.virtualId === vid
        );

        if (zone === "opponentshieldzone") {
          this.attackPlayer();
        } else if (zone === "opponentsplayzone") {
          this.attackCreature();
        }
      }
    }
  },
  created() {
    addEventListener("storage", this.onSettingsChanged);

    let lastReconnect = 0;

    const connect = async () => {
      if (this.errorMessage.includes("won")) {
        return;
      }

      if (this.preventReconnect) {
        return;
      }

      if (lastReconnect > 0) {
        console.log("Attempting to reconnect..");
      }

      if (lastReconnect > Date.now() - 5000) {
        setTimeout(connect, 1000);
        return;
      }

      // make sure the match actually exists
      try {
        await call({
          path: `/match/${this.$route.params.id}`,
          method: "GET"
        });
      } catch (err) {
        if (err.response && err.response.status == 404) {
          this.errorMessage = "This duel has been closed";
        }
      }

      lastReconnect = Date.now();

      // Connect to the server
      const ws = new WebSocket(
        ws_protocol + host + "/ws/" + this.$route.params.id
      );
      this.ws = ws;

      ws.onopen = () => {
        this.reconnecting = false;
        ws.send(localStorage.getItem("token"));
      };

      ws.onclose = () => {
        console.log("connection closed");
        this.reconnecting = true;
        connect();
      };

      ws.onerror = event => {
        console.error(event);
      };

      ws.onmessage = event => {
        const data = JSON.parse(event.data);

        switch (data.header) {
          case "mping": {
            send(ws, {
              header: "mpong"
            });
            break;
          }

          case "hello": {
            send(ws, {
              header: "join_match"
            });
            break;
          }

          case "error": {
            this.errorMessage = data.message;
            break;
          }

          case "warn": {
            this.warning = data.message;
            break;
          }

          case "opponent_disconnected": {
            this.opponentDisconnected = true;
            break;
          }

          case "opponent_reconnected": {
            this.opponentDisconnected = false;
            break;
          }

          // don't think this is ever fired, todo: remove or implement
          case "player_joined": {
            this.opponent = data.username;
            break;
          }

          case "choose_deck": {
            playerJoinedSound.play();
            document.title = "🔴 " + document.title;
            this.decks = data.decks;
            break;
          }

          case "chat": {
            this.chat(data.sender, data.color, data.message);
            break;
          }

          case "state_update": {
            if (!this.started) {
              this.started = true;
            }
            this.handSelection = null;
            this.playzoneSelection = null;

            if (this.state.myTurn !== data.state.myTurn) {
              turnSound.play();
              console.log("turn change");
            }

            this.state = data.state;
            break;
          }

          case "action": {
            (this.actionError = ""), (this.actionSelects = []);
            if (!(data.cards instanceof Array)) {
              this.actionObject = data.cards;
              console.log(Object.keys(data.cards)[0]);
              this.actionDrowdownSelection = Object.keys(data.cards)[0];
            }
            this.action = {
              cards:
                data.cards instanceof Array
                  ? data.cards
                  : Object.keys(data.cards)[0],
              text: data.text,
              minSelection: data.minSelection,
              maxSelections: data.maxSelections,
              cancellable: data.cancellable
            };
            break;
          }

          case "action_error": {
            if (!this.action) {
              return;
            }
            this.actionError = data.message;
            break;
          }

          case "close_action": {
            this.action = null;
            this.actionError = "";
            this.actionSelects = [];
            this.actionObject = null;
            this.actionDrowdownSelection = null;
            break;
          }

          case "wait": {
            this.wait =
              data.message || "Waiting for your opponent to make an action";
            break;
          }

          case "end_wait": {
            this.wait = "";
            break;
          }

          case "show_cards": {
            this.$modal.show(
              CardShowDialog,
              {
                message: data.message,
                cards: data.cards
              },
              {
                width: data.cards.length * 25 + "%"
              },
              {}
            );
          }
        }
      };
    };

    connect();

    // Loading dots
    setInterval(() => {
      if (this.loadingDots.length >= 4) this.loadingDots = "";
      else this.loadingDots += ".";
    }, 500);

    // clipboard
    let clipboard = new ClipboardJS("#invitebtn");
    clipboard.on("success", e => {
      if (this.inviteCopyTask) clearTimeout(this.inviteCopyTask);
      this.inviteCopied = true;
      this.inviteCopyTask = setTimeout(() => {
        this.inviteCopied = false;
      }, 2000);
      e.clearSelection();
    });
  },
  beforeDestroy() {
    removeEventListener("storage", this.onSettingsChanged);

    this.preventReconnect = true;
    this.ws.close();
  }
};
</script>

<style scoped lang="scss">
.spectator-info {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  margin-top: 7vh;
  div {
    margin: 0 3px;
  }
}

.spectator-leave {
  color: #999;
  text-align: center;
  margin-top: 10px;
  margin-bottom: 7vh;
  text-decoration: underline;
  color: #777;
  span {
    &:hover {
      cursor: pointer;
    }
  }
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
  img {
    width: 300px;
    border-radius: 15px;
    margin-bottom: 10px;
  }
}

.action-select {
  border: none;
  background: #484c52;
  padding: 5px !important;
  width: auto !important;
  margin-left: 5px;
  border-radius: 4px;
  color: #ccc;
  resize: none;
}
.action-select:focus {
  outline: none;
}
.action-select:active {
  outline: none;
}

.action-select {
  margin-top: 10px;
}

.action {
  max-height: 425px;
  width: 790px;
  background: #2f3136;
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
    padding: 15px 30px;
    &:hover {
      cursor: move;
    }
  }
  .btn {
    margin: 0 7px;
  }
  .action-cards {
    background: #222428;
    margin: 15px;
    margin-top: 0;
    border-radius: 4px;
    padding: 10px;
    max-height: 300px;
    overflow: auto;
    img {
      height: 125px;

      &.no-drag {
        user-drag: none;
        -webkit-user-drag: none;
        user-select: none;
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
      }
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
    width: 0;
  }
}

.glow {
  box-shadow: 0px 0px 4px 0px red;
}

.glow-fire {
  box-shadow: 0px 0px 4px 0px skyblue;
}
.glow-water {
  box-shadow: 0px 0px 4px 0px skyblue;
}
.glow-nature {
  box-shadow: 0px 0px 4px 0px skyblue;
}
.glow-light {
  box-shadow: 0px 0px 4px 0px skyblue;
}
.glow-darkness {
  box-shadow: 0px 0px 4px 0px skyblue;
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
  background: #2f3136;
  padding: 10px;
  border-radius: 4px;
}

.block {
  display: block !important;
}

.chatbox {
  height: calc(100vh - 128px - 15px);
  background: #2f3136;
  margin: 5px;
  border-radius: 4px;
}

.fullsize-chatbox {
  height: calc(100vh - 15px);
  background: #2f3136;
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
  background: #5b6eae;
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

.message {
  display: flex;
  margin-top: 5px;
  border-radius: 3px;
  padding: 3px 6px;
  line-height: 26px;

  .mute-icon-container {
    display: flex;
    align-items: center;
    opacity: 0;
  }

  &:hover .mute-icon-container {
    opacity: 1;
  }
}

.message-sender {
  margin-right: 5px;
}

.message-text {
  flex-grow: 1;
}

*::-webkit-scrollbar-track {
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  border-radius: 10px;
  background-color: #484c52;
}

*::-webkit-scrollbar {
  width: 6px;
  height: 6px;
  background-color: #484c52;
}

*::-webkit-scrollbar-thumb {
  border-radius: 10px;
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  background-color: #222;
}

.chatbox input {
  border: none;
  border-radius: 4px;
  margin: 10px;
  width: calc(100% - 40px);
  background: #484c52;
  padding: 10px;
  color: #ccc;
  &:focus {
    outline: none;
  }
  &:active {
    outline: none;
  }
}

.fullsize-chatbox input {
  border: none;
  border-radius: 4px;
  margin: 10px;
  width: calc(100% - 40px);
  background: #484c52;
  padding: 10px;
  color: #ccc;
  &:focus {
    outline: none;
  }
  &:active {
    outline: none;
  }
}

.actionbox {
  background: #2f3136;
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
  background: #36393f;
  z-index: 10;
  text-align: center;
  padding-top: 10vh;
}

.invite-link {
  padding: 5px;
  padding-left: 10px;
  background: #2b2e33;
  border: 1px solid #222428;
  border-radius: 4px;
  display: inline-block;
  color: #e3e3e5;
  transition: 0.1s;
}

.invite-link span {
  float: left;
  display: block;
  margin-top: 5px;
}

.invite-link .copy {
  display: inline-block;
  background: #7289da;
  color: #e3e3e5;
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
  background: #677bc4;
}

.copied {
  border-color: #3ca374 !important;
  color: #fff !important;
}

.invite-link > .copied {
  background: #3ca374;
}

.dots {
  display: inline-block;
  width: 40px;
  text-align: left;
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

.error p {
  padding: 5px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: #2b2e33 !important;
  border: 1px solid #222428;
}

.btn:hover {
  cursor: pointer;
  background: #677bc4;
}

.error {
  border: 1px solid #666;
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #36393f;
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

.noselect {
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Old versions of Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome, Edge, Opera and Firefox */
}
</style>
