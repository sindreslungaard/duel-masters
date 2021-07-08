<template>
  <header class="header">
    <nav>
      <span>Welcome, {{ username }}!</span>
      <ul>
        <li>
          <router-link
            v-slot="{ href }"
            to="lobby"
          >
            <BaseLink :href="href">
              Lobby
            </BaseLink>
          </router-link>
        </li>
        <li>
          <router-link
            v-slot="{ href }"
            to="decks"
          >
            <BaseLink :href="href">
              Decks
            </BaseLink>
          </router-link>
        </li>
        <li>
          <router-link
            v-slot="{ href }"
            to="logout"
          >
            <BaseLink :href="href">
              Logout
            </BaseLink>
          </router-link>
        </li>
        <li>
          <ExternalLink href="https://patreon.com/shobu">
            <img
              class="button--patreon"
              src="/assets/images/patreon.png"
            >
          </ExternalLink>
        </li>
      </ul>
    </nav>

    <div class="announcments">
      Join our
      <ExternalLink href="https://discord.gg/R3gJrX2">
        Discord
      </ExternalLink>
      to contribute with feedback and suggestions. Consider supporting us on <ExternalLink href="https://patreon.com/shobu">
        Patreon
      </ExternalLink> if you like what we're doing. ❤️
    </div>
  </header>
</template>

<script>
import BaseLink from "@/components/links/BaseLink";
import ExternalLink from "@/components/links/ExternalLink";

export default {
  name: "Header",
  components: {
    BaseLink,
    ExternalLink,
  },
  computed: {
    username: () => localStorage.getItem("username"),
  },
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
  margin: 0 var(--spacing);
}

.button--patreon {
  width: 125px;
  border-radius: var(--border-radius);
  opacity: 0.9;
}

.button--patreon:hover {
  opacity: 1;
}

.announcments {
  background: #2b2c31;
  padding: var(--spacing);
  border-radius: var(--border-radius);
  font-size: 14px;
  color: var(--color-text-light);
  margin-top: var(--spacing);
}
</style>
