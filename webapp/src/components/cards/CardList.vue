<template>
  <Panel
    title="Cards"
    class="card-list"
  >
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <div class="view-mode-switcher">
        <div class="field">
          Grid <input
            v-model="viewMode"
            value="grid"
            type="radio"
          > List
          <input
            v-model="viewMode"
            value="list"
            type="radio"
          >
        </div>
      </div>
      <div class="filter">
        <div class="field">
          <label for="filter-name">Name</label>
          <input
            id="filter-name"
            v-model="filterName"
            type="search"
            placeholder="Type to search"
          >
        </div>
        <div class="field">
          <label for="filter-set">Set</label>
          <select
            id="filter-set"
            v-model="filterSet"
          >
            <option value="all">
              All
            </option>
            <option
              v-for="(set, index) in sets"
              :key="index"
              :value="set"
            >
              {{
                set.toUpperCase()
              }}
            </option>
          </select>
        </div>
        <div class="field">
          <label for="filter-civilization">Civilization</label>
          <select
            id="filter-civilization"
            v-model="filterCivilization"
          >
            <option value="all">
              All
            </option>
            <option value="fire">
              Fire
            </option>
            <option value="water">
              Water
            </option>
            <option value="nature">
              Nature
            </option>
            <option value="light">
              Light
            </option>
            <option value="darkness">
              Darkness
            </option>
          </select>
        </div>
      </div>

      <div class="filter__count">
        {{ filteredCards.length }} cards
      </div>

      <table
        v-if="viewMode === 'list'"
        class="list-view"
      >
        <thead>
          <tr>
            <th>Name</th>
            <th>Set</th>
            <th>Civilization</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="card in filteredCards"
            :key="card.uid"
            @click.left="addCardToDeck(card)"
          >
            <td class="fill">
              {{ card.name }}
            </td>
            <td class="set no-break">
              {{ card.set }}
            </td>
            <td class="civilization no-break">
              {{ card.civilization }}
            </td>
          </tr>
        </tbody>
      </table>

      <div
        v-if="viewMode === 'grid'"
        class="grid-view"
      >
        <div
          v-for="card in filteredCards"
          :key="card.uid"
          @click.left="addCardToDeck(card)"
        >
          <img
            loading="lazy"
            :alt="card.name"
            :title="card.name"
            class="image--card"
            :src="`https://shobu.io/assets/cards/all/${card.uid}.jpg`"
          >
        </div>
      </div>
    </template>
  </Panel>
</template>

<script>
import Panel from "@/components/Panel";
import LoadingIndicator from "@/components/LoadingIndicator";
import { call } from "@/remote";
import _ from "lodash";

export default {
  name: "CardList",
  components: {
    Panel,
    LoadingIndicator,
  },
  props: {
    isLoading: {
      type: Boolean,
      required: false,
      default: false,
    },
    cards: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      viewMode: "grid",
      filterName: "",
      filterSet: "all",
      filterCivilization: "all",
    };
  },
  computed: {
    filteredCards() {
      let filteredCards = this.cards.filter(card =>
        card.name.toLowerCase().includes(this.filterName.toLowerCase()),
      );

      if (this.filterSet !== "all") {
        filteredCards = filteredCards.filter(
          card => card.set === this.filterSet,
        );
      }

      if (this.filterCivilization !== "all") {
        filteredCards = filteredCards.filter(
          card => card.civilization === this.filterCivilization,
        );
      }

      filteredCards = _.sortBy(filteredCards, "name");

      return filteredCards;
    },
    sets() {
      const sets = Array.from(new Set(this.cards.map(card => card.set)));
      sets.sort();
      return sets;
    },
    /**
     * Whether the component is ready to be displayed.
     */
    hasFinishedLoading() {
      return !this.isLoading && Object.keys(this.cards).length > 0;
    },
  },
  async created() {},
  methods: {
    addCardToDeck(card) {
      this.$store.commit("addCardToDeck", card);
    },
  },
};
</script>

<style lang="scss" scoped>
table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid var(--color-foreground-border);
  padding: calc(0.5 * var(--spacing));
  text-align: left;
}

.filter {
  display: flex;
  margin: 0 calc(-0.5 * var(--spacing));

  .field {
    width: 50%;
    padding: 0 calc(0.5 * var(--spacing));
  }

  &__count {
    margin-bottom: var(--spacing);
  }
}

.view-mode-switcher {
  margin-bottom: var(--spacing);
}

.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(125px, 1fr));
  grid-gap: calc(0.5 * var(--spacing));
}
</style>
