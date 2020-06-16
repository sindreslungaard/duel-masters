<template>
  <div class="overview">
    <div v-show="errorMessage" class="overlay"></div>

    <div v-show="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
      <div @click="refreshPage()" class="btn">Reconnect</div>
    </div>

    <main>
      <Header />
      <UserList :users="users" :isLoading="wsLoading" />
      <Chat
        @submit="sendChat(chatMessage)"
        v-model="chatMessage"
        :messages="chatMessages"
        :isLoading="wsLoading"
      />
      <DuelList
        @newDuel="openNewDuelDialog()"
        :matches="matches"
        :isLoading="wsLoading"
      />
    </main>
  </div>
</template>

<script>
import { call, ws_protocol } from "../remote";
import Header from "../components/Header.vue";
import UserList from "../components/chat/UserList";
import Chat from "../components/chat/Chat";
import DuelList from "../components/DuelList";
import NewDuelDialog from "../components/dialogs/NewDuelDialog";

const send = (client, message) => {
  client.send(JSON.stringify(message));
};

export default {
  name: "Overview",
  components: {
    Header,
    UserList,
    Chat,
    DuelList
  },
  computed: {
    username: () => localStorage.getItem("username")
  },
  data() {
    return {
      ws: null,
      chatMessage: "",
      chatMessages: [],
      users: [],
      matches: [],
      errorMessage: "",
      wsLoading: true
    };
  },
  methods: {
    refreshPage() {
      location.reload();
    },
    openNewDuelDialog() {
      this.$modal.show(NewDuelDialog);
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
    }
  },
  created() {
    // Loading dots
    setInterval(() => {
      if (this.loadingDots.length >= 4) this.loadingDots = "";
      else this.loadingDots += ".";
    }, 500);

    // Connect to the server
    try {
      const ws = new WebSocket(
        ws_protocol + window.location.host + "/ws/lobby"
      );
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
        }
      };
    } catch (err) {
      this.errorMessage = "Connection lost";
    }
  },
  beforeDestroy() {
    this.ws.close();
  }
};
</script>

<style scoped lang="scss">
input,
textarea,
select {
  border: none;
  background: var(--color-background-input);
  padding: 10px;
  border-radius: 4px;
  width: 200px;
  color: var(--color-text-light);
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
  color: var(--color-text-light);
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

.box {
  overflow: auto;
  background: #2b2c31;
  border-radius: 4px;
  font-size: 14px;
  color: var(--color-text-light);
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
  padding: 5px;
  border-radius: 4px;
  margin: 0;
  margin-bottom: 10px;
  background: #2b2e33 !important;
  border: 1px solid #222428;
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
  color: var(--color-text-light);
}

.new-duel-btn {
  margin-left: 7px;
  color: #fff;
  background: #3ca374;
  padding: 3px 5px;
  border-radius: 4px;
  font-size: 12px;
  text-transform: uppercase;
  font-weight: 400;
  position: absolute;
  top: -1px;
}

.new-duel-btn:hover {
  cursor: pointer;
  background: #35966a;
}

main {
  display: grid;
  grid-gap: 15px;
  padding: var(--spacing);
  grid-template-columns: minmax(auto, 15%) minmax(auto, 35%) minmax(auto, 50%);
  grid-template-rows: auto minmax(0, 1fr);
  width: 100%;
  height: 100vh;

  @include tablet {
    height: auto;
    grid-template-columns: minmax(auto, 25%) minmax(auto, 75%);
    grid-template-rows: auto auto auto;
  }

  @include mobile {
    grid-template-rows: repeat(4, auto);
  }

  .header {
    grid-column: 1 / 4;
    grid-row: 1;

    @include tablet {
      grid-column: 1 / 3;
    }
  }

  .user-list {
    display: flex;
    flex-direction: column;
    grid-column: 1;
    grid-row: 2 / 4;

    @include tablet {
      grid-column: 1;
      grid-row: 3;
    }

    @include mobile {
      grid-column: 1 / 3;
    }
  }

  .chat {
    display: flex;
    flex-direction: column;
    grid-column: 2;
    grid-row: 2 / 4;

    @include tablet {
      max-height: 400px;
      grid-column: 2;
      grid-row: 3;
    }

    @include mobile {
      grid-column: 1 / 3;
      grid-row: 4;
    }
  }

  .duels {
    display: flex;
    flex-direction: column;
    grid-column: 3;
    grid-row: 2 / 4;

    @include tablet {
      grid-column: 1/3;
      grid-row: 2;
    }
  }
}
</style>
