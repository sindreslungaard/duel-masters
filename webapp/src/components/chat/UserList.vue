<template>
  <Panel
    title="Online"
    class="user-list"
  >
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <template v-for="(userList, category) in users">
        <div
          v-if="userList.length > 0"
          :key="category"
          class="user-list__category"
        >
          <h4>{{ category }}</h4>
          <ul>
            <li
              v-for="user in userList"
              :key="user.username"
            >
              <Username
                :hub="user.hub"
                :color="user.color"
              >
                {{
                  user.username
                }}
              </Username>
            </li>
          </ul>
        </div>
      </template>
    </template>
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
    LoadingIndicator,
  },
  props: {
    users: {
      type: Object,
      required: true,
    },
    isLoading: {
      type: Boolean,
      required: true,
    },
  },
  computed: {
    /**
     * Whether the component is ready to be displayed.
     */
    hasFinishedLoading() {
      return !this.isLoading && Object.keys(this.users).length > 0;
    },
  },
};
</script>

<style lang="scss" scoped>
.user-list__category:not(:last-child) {
  margin-bottom: calc(2 * var(--spacing));
}

ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

li {
  margin-bottom: calc(0.5 * var(--spacing));
}
</style>
