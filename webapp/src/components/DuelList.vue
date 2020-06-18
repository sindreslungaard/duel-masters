<template>
  <Panel title="Duels" class="duels">
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <Button type="success" @click="openNewDuelDialog" text="New Duel" />
      <div v-if="!hasMatches" class="empty-message">
        <p>No matches to show, click the button above to create one.</p>
      </div>
      <div class="matches" v-if="hasMatches">
        <div class="match" v-for="match in matches" :key="match.id">
          <div class="match__owner">
            <Username :color="match.color">{{ match.owner }}</Username>
          </div>
          <div class="match__name">{{ match.name }}</div>
          <div class="match__actions">
            <router-link
              :to="{ name: 'duel', params: { id: match.id } }"
              v-slot="{ href }"
            >
              <ButtonLink
                :href="href"
                :type="match.spectate ? 'success' : 'info'"
                :text="match.spectate ? 'Spectate' : 'Join match'"
              ></ButtonLink>
            </router-link>
          </div>
        </div>
      </div>
    </template>
  </Panel>
</template>

<script>
import Panel from "./Panel";
import Button from "./buttons/Button";
import ButtonLink from "@/components/links/ButtonLink";
import LoadingIndicator from "./LoadingIndicator";
import NewDuelDialog from "@/components/dialogs/NewDuelDialog";
import Username from "@/components/chat/Username";

export default {
  name: "DuelList",
  components: {
    Panel,
    Button,
    LoadingIndicator,
    Username,
    ButtonLink
  },
  props: {
    matches: {
      type: Array,
      required: true
    },
    isLoading: {
      type: Boolean,
      required: true
    }
  },
  methods: {
    /**
     * Opens the dialog for creating a new duel.
     *
     * @returns {void}
     */
    openNewDuelDialog() {
      this.$modal.show(NewDuelDialog);
    }
  },
  computed: {
    /**
     * Whether the component is ready to be displayed.
     *
     * @returns {void}
     */
    hasFinishedLoading() {
      return !this.isLoading && this.matches;
    },
    /**
     * Whether there are any matches.
     *
     * @returns {void}
     */
    hasMatches() {
      return this.matches.length > 0;
    }
  }
};
</script>

<style lang="scss" scoped>
.matches {
  margin-top: var(--spacing);
  display: flex;
  flex-direction: column;
}

.match {
  display: flex;
  align-items: center;
  padding: var(--spacing);

  &:nth-child(odd) {
    background-color: var(--color-foreground-dark);
  }

  &__owner {
    width: 12ch;
    margin-right: var(--spacing);
  }

  &__actions {
    margin-left: auto;
  }
}
</style>
