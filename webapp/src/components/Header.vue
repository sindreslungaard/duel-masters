<template>
  <header class="header">
    <nav>
      <span class="welcome-message">Welcome, {{ username }}!</span>
      <ul>
        <li>
          <router-link to="overview" v-slot="{ href }">
            <InternalLink :href="href">Overview</InternalLink>
          </router-link>
        </li>
        <li>
          <router-link to="decks" v-slot="{ href }">
            <InternalLink :href="href">Decks</InternalLink>
          </router-link>
        </li>
        <li>
          <router-link to="logout" v-slot="{ href }">
            <InternalLink :href="href">Logout</InternalLink>
          </router-link>
        </li>
        <li>
          <ExternalLink href="https://patreon.com/shobu">
            <img class="patreon-btn" src="/assets/images/patreon.png" />
          </ExternalLink>
        </li>
      </ul>
    </nav>

    <div class="announcments">
      Join our
      <ExternalLink href="https://discord.gg/R3gJrX2">Discord</ExternalLink>
      to contribute with feedback and suggestions, or simply drop by to have a
      chat.
    </div>
  </header>
</template>

<script>
import InternalLink from "./links/InternalLink";
import ExternalLink from "./links/ExternalLink";

export default {
  name: "Header",
  components: {
    InternalLink,
    ExternalLink
  },
  computed: {
    username: () => localStorage.getItem("username")
  }
};
</script>

<style lang="scss" scoped>
nav {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
}

ul {
  list-style: none;
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-items: center;
  margin: 0;
  padding: 0;
}

li:not(:last-child)::after {
  content: "|";
  margin: 0 15px;
}

.patreon-btn {
  width: 125px;
  border-radius: 4px;
  opacity: 0.9;
}

.patreon-btn:hover {
  opacity: 1;
}

.announcments {
  background: #2b2c31;
  padding: var(--spacing);
  border-radius: 4px;
  font-size: 14px;
  color: var(--color-text-light);
  margin-top: 15px;
}
</style>
