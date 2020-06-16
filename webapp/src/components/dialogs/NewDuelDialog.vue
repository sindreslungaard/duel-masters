<template>
  <div class="new-duel-dialog">
    <form @submit.prevent="createDuel()">
      <h2>Create a new duel</h2>
      <div v-if="errorMessage" class="field error">
        <p>{{ errorMessage }}</p>
      </div>
      <div class="field">
        <label for="duel-name">Name</label>
        <input
          id="duel-name"
          autofocus
          v-model="duelData.name"
          type="text"
          placeholder="Name"
        />
      </div>
      <div class="field">
        <label for="duel-visibility">Visibility</label>
        <select id="duel-visibility" v-model="duelData.visibility">
          <option value="public">Show in list of duels</option>
          <option value="private">Hide from list of duels</option>
        </select>
      </div>
      <div class="field buttons">
        <Button submit text="Create" />
        <Button @click="$emit('close')" text="Cancel" type="secondary" />
      </div>
    </form>
  </div>
</template>

<script>
import Button from "../buttons/Button";
import { call } from "../../remote";

export default {
  name: "NewDuelDialog",
  data() {
    return {
      duelData: {
        name: "",
        description: "",
        visibility: "public"
      },
      errorMessage: ""
    };
  },
  props: {
    value: {}
  },
  components: {
    Button
  },
  methods: {
    async createDuel() {
      if (this.duelData.name.length < 5 || this.duelData.name.length > 30) {
        this.errorMessage = "Duel name must be between 5-30 characters";
        return;
      }

      try {
        let res = await call({
          path: "/match",
          method: "POST",
          body: this.duelData
        });

        this.$emit("close");
        this.$router.push({ path: "/duel/" + res.data.id });
      } catch (e) {
        try {
          console.log(e);
          this.wizardError = e.response.data.message;
        } catch (err) {
          console.log(err);
          this.errorMessage =
            "Unable to communicate with the server. Please try again later.";
        }
      }
    }
  }
};
</script>

<style lang="scss" scoped>
form {
  display: flex;
  flex-direction: column;
}
</style>
