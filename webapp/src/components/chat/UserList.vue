<template>
  <Panel title="Online" class="user-list">
    <LoadingIndicator v-if="isLoading" />

    <div
      class="user-list__category"
      v-for="(category, index) in users"
      :key="index"
    >
      <h4>{{ category.category }}</h4>
      <ul>
        <li v-for="(user, index) in category.users" :key="index">
          <Username :hub="user.hub" :color="user.color">{{
            user.username
          }}</Username>
        </li>
      </ul>
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
  props: {
    users: {
      type: Array,
      required: true
    },
    isLoading: {
      type: Boolean,
      required: true
    }
  }
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
