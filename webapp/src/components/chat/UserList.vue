<template>
  <Panel title="Online" class="user-list">
    <!-- TODO: Make a Loader component -->
    <div v-if="isLoading">Loading</div>

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

export default {
  name: "UserList",
  components: {
    Username,
    Panel
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
.user-list__category {
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
