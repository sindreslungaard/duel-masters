<template>
  <div class="deck-edit-dialog">
    <form @submit.prevent="updateDeck">
      <h2>Edit your deck</h2>
      <div class="field">
        <label for="deck-name">Name</label>
        <input
          id="deck-name"
          v-model="deckData.name"
          autofocus
          type="text"
          placeholder="Name"
        >
      </div>
      <div class="field">
        <label for="deck-visibility">Visibility</label>
        <select
          id="deck-visibility"
          v-model="deckData.visibility"
        >
          <option value="true">
            Public
          </option>
          <option value="false">
            Private
          </option>
        </select>
      </div>
      <div class="field buttons">
        <Button
          submit
          text="Done"
        />
      </div>
    </form>
  </div>
</template>
<script>
import Button from "../buttons/Button";

export default {
  name: "DeckEditDialog",
  components: {
    Button,
  },
  data() {
    return {
      deckData: {
        name: null,
        visibility: null,
      },
    };
  },
  created() {
    this.deckData.name = this.$store.state.deck.name;
    this.deckData.visibility = this.$store.state.deck.public;
  },
  methods: {
    updateDeck() {
      this.$store.commit("updateDeckName", this.deckData.name);
      this.$store.commit("updateDeckVisibility", !!this.deckData.visibility);
      this.$emit("close");
    },
  },
};
</script>

<style lang="scss" scoped>
form {
  display: flex;
  flex-direction: column;
}
</style>
