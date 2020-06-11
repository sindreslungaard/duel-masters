<template>
  <Panel title="Chat" class="chat">
    <LoadingIndicator v-if="isLoading" />

    <div class="chatbox">
      <div id="messages" class="messages spaced">
        <div class="messages-helper">
          <div v-for="(msg, i) in messages" :key="i">
            <Username :color="msg.color">{{ msg.username }}</Username>
            <div class="user-messages">
              <div v-for="(message, j) in msg.messages" :key="j">
                <span>{{ message }}</span>
              </div>
            </div>
          </div>
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
    </div>
  </Panel>
</template>

<script>
import Username from "./Username";
import Panel from "../Panel";
import LoadingIndicator from "../LoadingIndicator";

export default {
  name: "UserList",
  components: {
    Username,
    Panel,
    LoadingIndicator
  },
  data() {
    return {
      chatMessage: this.value
    };
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
  }
};
</script>

<style lang="scss" scoped>
.chatbox {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
}

.chatbox input {
  border: none;
  border-radius: 4px;
  width: 100%;
  background: var(--color-background-input);
  padding: 10px;
  color: var(--color-text-light);
  &:focus {
    outline: none;
  }
  &:active {
    outline: none;
  }
}

.user-name-container {
  overflow: hidden;
}

.user-messages {
  margin-left: 20px;
  margin-top: 0px;
  margin-bottom: 15px;
}

.user-messages > div {
  margin: 3px 0;
  color: #fff;
}

.messages {
  overflow: auto;
  margin-bottom: 0;
  padding-bottom: 0;
}
</style>
