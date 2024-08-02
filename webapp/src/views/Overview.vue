<template>
  <div>
    <div
      v-show="errorMessage || wizardVisible || warning"
      @click="closeOverlay()"
      class="overlay"
    ></div>

    <div v-show="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
      <div @click="refreshPage()" class="btn">Reconnect</div>
    </div>

    <div v-show="warning" class="error">
      <p>{{ warning }}</p>
      <div @click="warning = ''" class="btn">Ok</div>
    </div>

    <div v-show="wizardVisible" class="new-duel">
      <div class="wizard">
        <div class="spacer">
          <span class="headline">Create a new duel</span>
          <br /><br />
          <form @submit="handleSubmit()" v-on:submit.prevent="onSubmit">
            <input v-model="wizard.name" type="text" placeholder="Name" />
            <br /><br />
            <span class="helper">Visibility</span>
            <select v-model="wizard.visibility">
              <option value="public">Show in list of duels</option>
              <option value="private">Hide from list of duels</option>
            </select>

            <span v-if="wizardError" class="errorMsg">{{ wizardError }}</span>

            <div @click="createDuel()" class="btn">
              Create
            </div>
            <div @click="toggleWizard()" class="btn cancel">
              Cancel
            </div>
          </form>
        </div>
      </div>
    </div>

    <main>
      <Header style="width: 100%"></Header>

      <div class="spaced">
        <div class="categories">
          <h3 class="user-list">Online</h3>
          <h3 class="chat">Chat</h3>
          <h3 class="duels" style="position: relative;">
            Duels<span @click="toggleWizard(true)" class="new-duel-btn"
              >New Duel</span
            >
          </h3>
        </div>

        <!-- Users online -->
        <div class="box user-list">
          <div class="spaced">
            <div v-if="wsLoading">Loading{{ loadingDots }}</div>

            <div
              class="user-name-container"
              v-for="(category, index) in users"
              :key="index"
            >
              <div class="user-category">
                <span>{{ category.category }}</span>
              </div>

              <Username
                v-for="(user, index) in category.users"
                :key="index"
                :hub="user.hub"
                :color="user.color"
                >{{ user.username }}</Username
              >

              <br />
            </div>
          </div>
        </div>

        <!-- Chat -->
        <div class="box chat">
          <div v-if="wsLoading" class="spaced" style="position: absolute">
            Loading{{ loadingDots }}
          </div>

          <div class="chatbox">
            <div v-if="pinnedMessages.length" class="pinned-messages">
              <div v-for="(msg, i) in pinnedMessages" :key="i">
                {{ msg.message }}
                <span v-if="msg.timeString">{{ msg.timeString }}</span>
              </div>
            </div>
            <div id="messages" class="messages spaced">
              <div class="messages-helper">
                <div
                  v-for="(msg, i) in chatMessages.filter(
                    m => !settings.muted.includes(m.username)
                  )"
                  :key="i"
                >
                  <Username
                    :color="msg.color"
                    :muteUser="
                      ['[Server]', username].includes(msg.username)
                        ? null
                        : msg.username
                    "
                    @muteToggled="onSettingsChanged()"
                    >{{ msg.username }}
                    <span class="message-ts"
                      >{{ tsformat(msg.timestamp) }} ago
                    </span>
                  </Username>
                  <div class="user-messages">
                    <div v-for="(message, j) in msg.messages" :key="j">
                      <span>{{ message }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <form @submit.prevent="sendChat(chatMessage)">
              <input
                type="text"
                v-model="chatMessage"
                placeholder="Type to chat"
              />
            </form>
          </div>
        </div>

        <!-- Duels -->
        <div class="box duels">
          <div v-if="wsLoading" class="spaced">Loading{{ loadingDots }}</div>

          <table>
            <!-- Loading -->
            <tr
              v-if="
                !wsLoading && matches.length < 1 && matchRequests.length < 1
              "
            >
              <td>No matches to show, click the button above to create one.</td>
            </tr>

            <!-- Match requests -->
            <tr v-for="(request, index) in matchRequests" :key="index">
              <td style="width: 30%">
                <div class="match-players">
                  <Username :color="request.host_color">{{
                    request.host_name
                  }}</Username>
                  <div v-show="request.guest_id">vs</div>
                  <Username
                    v-show="request.guest_id"
                    :color="request.guest_color"
                    >{{ request.guest_name }}</Username
                  >
                  <span
                    v-show="request.host_id == uid && request.guest_id"
                    @click="kickPlayer(request)"
                    style="color: #f36a6a; text-decoration: underline dotted; cursor: pointer; font-weight: bold"
                    >remove</span
                  >
                </div>
              </td>
              <td style="width: 45%">
                {{
                  request.guest_id == uid
                    ? "Waiting for the host to start the match" + loadingDots
                    : request.name
                }}
              </td>
              <td style="width: 25%">
                <div
                  @click="leaveMatch(request)"
                  v-show="request.host_id == uid && !request.guest_id"
                  class="btn-colorless bg-red-500 hover:bg-red-600 cursor-pointer"
                  style="margin-left: 10px"
                >
                  Close
                </div>

                <div
                  @click="startMatch(request)"
                  v-show="request.host_id == uid && request.guest_id"
                  class="btn save"
                >
                  Start
                </div>

                <div
                  @click="leaveMatch(request)"
                  v-show="request.guest_id == uid"
                  class="btn-colorless bg-red-500 hover:bg-red-600 cursor-pointer"
                >
                  Leave
                </div>

                <div
                  @click="joinMatch(request)"
                  v-show="request.host_id != uid && !request.guest_id"
                  class="btn save"
                >
                  Join match
                </div>

                <div
                  v-show="
                    request.guest_id != '' &&
                      request.host_id != uid &&
                      request.guest_id != uid
                  "
                  class="float-right opacity-50"
                >
                  <div style="width: 120px; padding-top: 8px; padding-bottom: 7px;">
                    Waiting to start{{ loadingDots }}
                  </div>
                </div>

                <div
                  v-show="request.host_id == uid && !request.guest_id"
                  @click="copyToClipboard(protocol + '//' + host + '/invite/' + request.link_code)"
                  :class="['copy', { copied: inviteCopied }]"
                >
                  {{ inviteCopied ? "Copied" : "Copy invite link" }}
                </div>
              </td>
            </tr>

            <!-- Matches -->
            <tr v-for="(match, index) in matches" :key="index">
              <td>
                <div class="match-players">
                  <Username :color="match.p1color">{{ match.p1 }}</Username>
                  <div v-show="match.p2">vs</div>
                  <Username v-show="match.p2" :color="match.p2color">{{
                    match.p2
                  }}</Username>
                </div>
              </td>
              <td>{{ match.name }}</td>
              <td>
                <div
                  @click="$router.push('/duel/' + match.id)"
                  :class="'btn' + (match.spectate ? '' : ' save')"
                >
                  {{ match.spectate ? "Spectate" : "Join match" }}
                </div>
              </td>
            </tr>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { call, ws_protocol, host } from "../remote";
import Header from "../components/Header.vue";
import Username from "../components/Username.vue";
import { getSettings, didSeeMuteWarning } from "../helpers/settings";
import 'setimmediate';

import {
  format,
  fromUnixTime,
  formatDistanceToNowStrict,
  isBefore,
  formatDistance
} from "date-fns";

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

let requestAcceptedSound = new sound("/assets/request_accepted.wav");
let playerLeftSound = new sound("/assets/player_left.wav");

export default {
  name: "overview",
  components: {
    Header,
    Username
  },
  computed: {
    uid: () => localStorage.getItem("uid"),
    username: () => localStorage.getItem("username"),
    protocol: () => window.location.protocol
  },
  data() {
    return {
      ws: null,
      wizardVisible: false,
      wizardError: "",
      wizard: {
        name: "",
        description: "",
        visibility: "public"
      },
      chatMessage: "",
      chatMessages: [],
      pinnedMessages: [], // { message, time }
      users: [],
      matches: [],
      matchRequests: [],
      errorMessage: "",
      warning: "",
      wsLoading: true,
      loadingDots: ".",
      settings: getSettings(),
      inviteCopied: false,
      inviteCopyTask: null,
      host,
      linkCode: "",
    };
  },
  methods: {
    onSettingsChanged(e) {
      if (!e && !didSeeMuteWarning()) {
        this.warning =
          "You can unmute players at any time from the settings page";
      }

      this.settings = getSettings();
    },
    tsformat(ts) {
      return formatDistance(Date.now(), fromUnixTime(ts));
    },
    refreshPage() {
      location.reload();
    },
    toggleWizard(state) {
      this.wizardError = "";
      this.wizard = {
        name: "",
        description: "",
        visibility: "public"
      };
      this.wizardVisible = state;
    },
    closeOverlay() {
      this.toggleWizard();
      this.errorMessage = "";
      this.warning = "";
    },
    handleSubmit() {
      this.createDuel();
    },
    async createDuel() {
      if (this.wizard.name != "" && this.wizard.name.length > 30) {
        this.wizardError = "Duel name cannot exceed 30 characters";
        return;
      }

      if (this.wizard.visibility == "private") {
        try {
          let res = await call({
            path: "/match",
            method: "POST",
            body: this.wizard
          });

          this.$router.push({ path: "/duel/" + res.data.id });
        } catch (e) {
          try {
            console.log(e);
            this.wizardError = e.response.data.message;
          } catch (err) {
            console.log(err);
            this.wizardError =
              "Unable to communicate with the server. Please try again later.";
          }
        }

        return;
      }

      this.ws.send(
        JSON.stringify({
          header: "create_match_request",
          name: this.wizard.name
        })
      );
      this.wizardVisible = false;
      this.wizard = {
        name: "",
        description: "",
        visibility: "public"
      };
    },
    sendChat(message) {
      if (!message) {
        return;
      }
      this.chatMessage = "";
      this.ws.send(JSON.stringify({ header: "chat", message }));
    },
    chat(data) {
      let createNew = true;

      if (this.chatMessages.length > 0) {
        let lastMsg = this.chatMessages[this.chatMessages.length - 1];

        if (
          lastMsg.username == data.username &&
          lastMsg.timestamp > data.timestamp - 15
        ) {
          lastMsg.messages.push(data.message);
          createNew = false;
        }
      }

      if (createNew) {
        this.chatMessages.push({
          username: data.username,
          color: data.color,
          timestamp: data.timestamp,
          messages: [data.message]
        });
      }

      this.$nextTick(() => {
        let container = document.getElementById("messages");
        container.scrollTop = container.scrollHeight;
      });
    },

    leaveMatch(request) {
      this.ws.send(JSON.stringify({ header: "leave_match_request" }));
    },

    joinMatch(request) {
      this.ws.send(
        JSON.stringify({ header: "join_match_request", id: request.id })
      );
    },

    startMatch(request) {
      this.ws.send(JSON.stringify({ header: "start_match", id: request.id }));
    },

    kickPlayer(request) {
      this.ws.send(
        JSON.stringify({
          header: "kick_from_request",
          id: request.id,
          player_id: request.guest_id
        })
      );
    },

    copyToClipboard(text) {
      navigator.clipboard.writeText(text);

      if (this.inviteCopyTask) clearTimeout(this.inviteCopyTask);
      this.inviteCopied = true;
      this.inviteCopyTask = setTimeout(() => {
        this.inviteCopied = false;
      }, 2000);
    }
  },
  created() {
    if(this.$route.query.invite) {
      this.linkCode = this.$route.query.invite;
    }

    addEventListener("storage", this.onSettingsChanged);

    document.title = document.title.replace("🔴", "");

    // Loading dots
    setInterval(() => {
      if (this.loadingDots.length >= 4) this.loadingDots = "";
      else this.loadingDots += ".";
    }, 500);

    let timeUpdates = setInterval(() => {
      for (let msg of this.pinnedMessages) {
        if (!msg.time) continue;

        if (isBefore(fromUnixTime(msg.time), Date.now())) {
          clearInterval(timeUpdates);
          return;
        }

        msg.timeString = formatDistanceToNowStrict(fromUnixTime(msg.time));
      }
    }, 500);

    // Connect to the server
    try {
      const ws = new WebSocket(ws_protocol + host + "/ws/lobby");
      this.ws = ws;

      ws.onopen = () => {
        ws.send(localStorage.getItem("token"));
        this.wsLoading = false;
      };

      ws.onclose = () => {
        this.errorMessage = "Lost connection to the server";
      };

      ws.onerror = () => {
        this.errorMessage = "Lost connection to the server";
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
              header: "subscribe"
            });
            break;
          }

          case "chat": {
            for (let message of data.messages) {
              this.chat(message);
            }
            break;
          }

          case "pinned_messages": {
            this.pinnedMessages = [];

            for (let message of data.messages) {
              if (message.includes("time:")) {
                let time = message.split("time:")[1];
                this.pinnedMessages.push({
                  message: message.split("time:")[0],
                  time,
                  timeString: formatDistanceToNowStrict(fromUnixTime(time))
                });
              } else {
                this.pinnedMessages.push({ message });
              }
            }

            setImmediate(() => {
              let container = document.getElementById("messages");
              container.scrollTop = container.scrollHeight;
            });

            break;
          }

          case "users": {
            this.users = [
              {
                category: "player",
                users: []
              }
            ];

            for (let user of data.users) {
              let chatroles = user.permissions.filter(x =>
                x.includes("chat.role.")
              );

              if (chatroles.length > 0) {
                let role = chatroles[0].split("chat.role.")[1];

                let category = this.users.find(x => x.category == role);

                let toPushCategory = false;
                if (!category) {
                  category = {
                    category: role,
                    users: []
                  };
                  toPushCategory = true;
                }

                category.users.push(user);
                if (toPushCategory) {
                  this.users.push(category);
                }
              } else {
                let category = this.users.find(x => x.category == "player");
                category.users.push(user);
              }

              this.users.sort((a, b) => a.category.localeCompare(b.category));
            }
            break;
          }

          case "matches": {
            this.matches = data.matches;
            break;
          }

          case "match_requests": {
            this.matchRequests = data.requests;
            if(this.linkCode != "") {
              let found = false;
              for(let req of data.requests) {
                console.log("comparing", req.link_code, this.linkCode)
                if(req.link_code == this.linkCode) {
                  found = true;
                  this.joinMatch(req)
                }
              }
              if(!found) {
                this.chat({
                  username: "[Server -> you]",
                  color: "#777",
                  timestamp: Math.round(Date.now() / 1000),
                  message: "Could not find the duel you were invited to. It has probably been closed or started already."
                })
              }
              this.linkCode = "";
            }
            
            break;
          }

          case "warn": {
            this.warning = data.message;
            break;
          }

          case "match_forward": {
            this.$router.push({ path: "/duel/" + data.id });
            break;
          }

          case "play_sound": {
            switch (data.sound) {
              case "request_accepted":
                requestAcceptedSound.play();
                break;
              case "player_left":
                playerLeftSound.play();
                break;
            }
            break;
          }
        }
      };
    } catch (err) {
      this.errorMessage = "Connection lost";
    }
  },
  beforeDestroy() {
    removeEventListener("storage", this.onSettingsChanged);
    this.ws.close();
  }
};
</script>

<style scoped lang="scss">
.match-players {
  display: flex;
  div {
    margin: 0 3px;
  }
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
  background: #0c0c0f;
  border: 1px solid #333;
  width: 250px;
  border-radius: 4px;
  color: #fff;
  z-index: 100;
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

select option {
  background: #111;
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
  background: url(/assets/images/overlay_30.png);
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
  background: #5865f2;
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
  background: #515de2;
}

.btn:active {
  background: #4c58d3 !important;
}

.btn-colorless {
  display: inline-block;
  color: #e3e3e5;
  font-size: 14px;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn-colorless:hover {
  cursor: pointer;
}

main {
  width: 100%;
  height: 100vh;
  margin: auto;
  overflow: hidden;
}

.box {
  overflow: auto;
  background: url(/assets/images/overlay_30.png);
  min-height: 20px;
  border-radius: 4px;
  font-size: 14px;
  color: #ccc;
  display: inline-block;
  height: calc(100vh - 140px);
}

.user-list {
  width: 10%;
}

.chat {
  width: calc(35% - 15px);
  margin-left: 15px;
}

.duels {
  width: calc(55% - 15px);
  margin-left: 15px;
}

.spaced {
  margin: 15px;
}

.categories > h3 {
  margin-top: 0;
  margin-bottom: 7px;
  display: inline-block;
  color: #eee;
  font-weight: 400;
  font-size: 16px;
}

.duels > table {
  width: 100%;
  border-collapse: collapse;
}

.duels td {
  border: none;
  text-align: left;
  padding: 15px;
}

.duels tr:nth-child(odd) {
  background: url(/assets/images/overlay_30.png);
}

.duels .btn {
  float: right;
}

.duels .btn-colorless {
  float: right;
}

.save {
  background: #3ca374 !important;
}

.save:hover {
  background: #35966a !important;
}

.user-category {
  margin-bottom: 16px;
  border-bottom: 1px solid #333;
  color: #555;
  padding-bottom: 5px;
  font-weight: 400;
  text-transform: capitalize;
}

.chatbox {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.chatbox input {
  border: none;
  border-radius: 4px;
  margin: 10px;
  width: calc(100% - 40px);
  background: url(/assets/images/overlay_30.png);
  padding: 10px;
  color: #ccc;
  &:focus {
    outline: none;
  }
  &:active {
    outline: none;
  }
}

.chatbox form {
  justify-self: end;
}

.duels .btn {
  width: 70px;
}

.duels .btn-colorless {
  width: 70px;
}

.user-list .user-name {
  margin-bottom: 10px;
}

.user-name-container {
  overflow: hidden;
}

.user-messages {
  margin-left: 20px;
  margin-top: 0px;
  margin-bottom: 15px;
  word-break: break-word;
}

.user-messages > div {
  margin: 3px 0;
  color: #e1e1e1;
}

.messages {
  overflow: auto;
  margin-bottom: 0;
  padding-bottom: 0;
  flex-grow: 1;

  &-helper {
    margin-right: 10px;
  }
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

.error p {
  padding: 10px 12px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: url(/assets/images/overlay_15.png) !important;
}

.error {
  border: 1px solid #333;
  position: absolute;
  top: 0;
  left: 0;
  width: 300px;
  border-radius: 4px;
  background: #111214;
  z-index: 3005;
  left: calc(50% - 300px / 2);
  top: 40vh;
  padding: 15px;
  font-size: 14px;
  color: #ccc;
}

.new-duel-btn {
  margin-left: 7px;
  color: #fff;
  background: #3ca374;
  padding: 3px 5px;
  border-radius: 4px;
  font-size: 12px;
  text-shadow: 1px 1px #0f2c1f;
  font-weight: 600;
  text-transform: uppercase;
  position: absolute;
  top: -1px;
}

.new-duel-btn:hover {
  cursor: pointer;
  background: #35966a;
}

.pinned-messages {
  background: url(/assets/images/overlay_30.png);
  padding: 10px;
  font-size: 13px;
  color: yellow;
}

.message-ts {
  font-size: 11px;
  color: #999;
  text-shadow: none;
  opacity: 0.5;
}

.kick-btn {
  text-decoration-style: dotted;
  color: white;
  border-radius: 4px;
  padding: 2px 6px;
  background: rgb(243, 106, 106);
  text-shadow: 1px 1px #0f2c1f;
}

.kick-btn:hover {
  cursor: pointer;
  background: #e64343;
}

.copy {
  text-decoration: underline dotted;
  color: #666;
  font-size: 14px;
  transition: 0.1s;
  padding-top: 7px;
  width: 100px;
  text-align: right;
  float: right;
}

.copy:hover {
  cursor: pointer;
}

.copied {
  color: #3ca374;
}
</style>
