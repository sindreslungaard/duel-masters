<template>
  <main class="overview">
    <Header />
    <UserList
      :users="users"
      :is-loading="isLoading"
    />
    <Chat
      v-model="chatMessage"
      :messages="chatMessages"
      :is-loading="isLoading"
      @submit="sendChatMessage"
    />
    <DuelList
      :matches="matches"
      :is-loading="isLoading"
    />
  </main>
</template>

<script>
import Header from "../components/Header.vue";
import UserList from "../components/chat/UserList";
import Chat from "../components/chat/Chat";
import DuelList from "../components/DuelList";
import BaseMixin from "@/mixins/BaseMixin";
import LobbyMessageHandler from "@/message_handler/LobbyMessageHandler";

export default {
  name: "Overview",
  components: {
    Header,
    UserList,
    Chat,
    DuelList,
  },
  mixins: [BaseMixin],
  data() {
    return {
      /**
       * The WebSocket used for handling lobby messages.
       *
       * @type {WebSocket}
       */
      socket: null,
      /**
       * The message that is currently being typed.
       *
       * @type {String}
       */
      chatMessage: "",
      /**
       * Array of chat messages.
       *
       * @type {Array.Object}
       */
      chatMessages: [],
      /**
       * An object containing user categories with the users belonging
       * to that category.
       *
       * @type {Object}
       */
      users: {},
      /**
       * Array of matches.
       *
       * @type {Array.Object}
       */
      matches: [],
    };
  },
  /**
   * When the component is created we connect to the WebSocket and
   * setup the message handler.
   */
  created() {
    document.title = document.title.replace("🔴", "");

    try {
      this.socket = this.connectToSocket("/ws/lobby");
      new LobbyMessageHandler(this.socket, this);
    } catch (error) {
      console.error(error);
      this.showError("Lost connection to the server");
    }
  },
  beforeUnmount() {
    this.socket.close();
  },
  methods: {
    /**
     * Sends a chat message to the WebSocket.
     *
     * @return {void}
     */
    sendChatMessage() {
      if (!this.chatMessage) {
        return;
      }

      this.sendMessage(this.socket, {
        header: "chat",
        message: this.chatMessage,
      });
      this.chatMessage = "";
    },
    /**
     * Adds a new message to the chat.
     *
     * @param {Object} data
     * @returns {void}
     */
    addChatMessage(data) {
      // If the last message was from the same user and it was
      // sent inside the append duration, we append the mssage to the
      // last one.
      if (this.chatMessages.length > 0) {
        const lastMessage = this.chatMessages[this.chatMessages.length - 1];

        if (
          lastMessage.username == data.username &&
          lastMessage.timestamp >
            data.timestamp - this.$config.CHAT_MESSAGE_APPEND_TIMEOUT
        ) {
          lastMessage.messages.push(data.message);
          return;
        }
      }

      this.chatMessages.push({
        username: data.username,
        color: data.color,
        timestamp: data.timestamp,
        messages: [data.message],
      });
    },
  },
};
</script>

<style scoped lang="scss">
main {
  display: grid;
  grid-gap: var(--spacing);
  padding: var(--spacing);
  grid-template-columns: minmax(auto, 15%) minmax(auto, 35%) minmax(auto, 50%);
  grid-template-rows: auto minmax(0, 1fr);
  width: 100%;
  height: 100vh;

  @include tablet {
    height: auto;
    grid-template-columns: minmax(auto, 25%) minmax(auto, 75%);
    grid-template-rows: repeat(3, auto);
  }

  @include mobile {
    grid-template-rows: repeat(4, auto);
  }
}

.header {
  grid-column: 1 / 4;
  grid-row: 1;

  @include tablet {
    grid-column: 1 / 3;
  }
}

.user-list {
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
  grid-column: 3;
  grid-row: 2 / 4;

  @include tablet {
    grid-column: 1/3;
    grid-row: 2;
  }
}
</style>
