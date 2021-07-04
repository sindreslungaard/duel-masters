<template>
  <Panel title="Deck" class="deck-view">
    <LoadingIndicator v-if="!hasFinishedLoading" />

    <template v-if="hasFinishedLoading">
      <div class="field deck-selection">
        <label for="decks">Decks</label>
        <select
          @input="$emit('input', selectedDeck)"
          id="decks"
          v-model="selectedDeck"
        >
          <option v-for="deck in decks" :key="deck.uid" :value="deck.uid">{{
            deck.name
          }}</option>
        </select>
      </div>

      <div class="actions">
        <Button type="success" text="New Deck"></Button>
        <Button text="Edit Deck"></Button>
        <Button type="error" text="Delete Deck"></Button>
      </div>
    </template>
  </Panel>
</template>

<script>
import Panel from "@/components/Panel";
import LoadingIndicator from "@/components/LoadingIndicator";
import Button from "@/components/buttons/Button";
import { call } from "@/remote";
import _ from "lodash";

export default {
  name: "DeckView",
  components: {
    Panel,
    LoadingIndicator,
    Button
  },
  props: {
    isLoading: {
      type: Boolean,
      required: false,
      default: false
    },
    value: {}
  },
  data() {
    return {
      decks: [],
      selectedDeck: this.value
    };
  },
  async created() {
    let decksResponse = await call({ path: "/decks", method: "GET" });
    this.decks = decksResponse.data;

    // TODO: Check if 0 decks.

    this.selectedDeck = this.decks[0].uid;
  },
  methods: {
    onLeftClick(card) {
      this.$emit("leftClick", card);
    },
    onMiddleClick(card) {
      this.$emit("middleClick", card);
    },
    onRightClick(card) {
      this.$emit("rightClick", card);
    }
  },
  watch: {
    value(value) {
      this.selectedDeck = value;
    },
    selectedDeck(value) {
      const deck = this.decks.filter(deck => deck.uid === value)[0];
      console.log(deck);
    }
  },
  computed: {
    /**
     * Whether the component is ready to be displayed.
     */
    hasFinishedLoading() {
      return !this.isLoading && this.decks.length > 0;
    }
  }
};
</script>
