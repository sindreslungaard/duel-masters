<template>
  <Panel title="Chat" class="chat">
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <div class="message" v-for="(msg, i) in messages" :key="i">
        <div class="message__meta">
          <Username :color="msg.color">{{ msg.username }}</Username>
          <span>{{ formatTimestamp(msg.timestamp) }}</span>
        </div>

        <div class="message__content">
          <p v-for="(message, j) in msg.messages" :key="j">
            {{ message }}
          </p>
        </div>
      </div>

      <form @submit.prevent="$emit('submit')">
        <input
          @input="$emit('input', chatMessage)"
          type="text"
          v-model="chatMessage"
          placeholder="Type to chat"
        />
      </form>
    </template>
  </Panel>
</template>

<script>
import Username from "./Username";
import Panel from "../Panel";
import LoadingIndicator from "../LoadingIndicator";

export default {
  name: "Chat",
  components: {
    Username,
    Panel,
    LoadingIndicator
  },
  props: {
    messages: {
      type: Array,
      required: true
    },
    isLoading: {
      type: Boolean,
      required: true
    },
    value: {}
  },
  data() {
    return {
      /**
       * The chat message being types.
       *
       * @type {String}
       */
      chatMessage: this.value
    };
  },
  methods: {
    /**
     * Formats a timestamp to a string according to the locale.
     *
     * @param {Number} dateTime
     * @return {String}
     */
    formatTimestamp(dateTime) {
      return new Date(dateTime * 1000).toLocaleString();
    }
  },
  computed: {
    /**
     * Whether the component is ready to be displayed.
     */
    hasFinishedLoading() {
      return !this.isLoading && Object.keys(this.messages).length > 0;
    }
  },
  watch: {
    /**
     * When the v-model is updated from outside, we need to move it
     * to our "internal" v-model.
     *
     * @param {String} value
     * @returns {void}
     */
    value(value) {
      this.chatMessage = value;
    },
    /**
     * When new message are added we scroll to the bottom of the chat.
     *
     * @returns {void}
     */
    messages() {
      this.$nextTick(() => {
        const chatContainer = document.querySelector(".chat .panel__content");
        chatContainer.scrollTop = chatContainer.scrollHeight;
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.message {
  &__meta {
    display: flex;
    justify-content: space-between;
  }

  &__content {
    margin: 0 var(--spacing) var(--spacing);
  }
}

p {
  margin: 0;
}
</style>
