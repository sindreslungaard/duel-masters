<template>
  <div class="new-duel-dialog">
    <form @submit.prevent="createDuel()">
      <h2>Create a new duel</h2>
      <div
        v-if="errorMessage"
        class="field error"
      >
        <p>{{ errorMessage }}</p>
      </div>
      <div class="field">
        <label for="duel-name">Name</label>
        <input
          id="duel-name"
          v-model="duelData.name"
          autofocus
          type="text"
          placeholder="Name"
        >
      </div>
      <div class="field">
        <label for="duel-visibility">Visibility</label>
        <select
          id="duel-visibility"
          v-model="duelData.visibility"
        >
          <option value="public">
            Show in list of duels
          </option>
          <option value="private">
            Hide from list of duels
          </option>
        </select>
      </div>
      <div class="field buttons">
        <Button
          submit
          text="Create"
        />
        <Button
          text="Cancel"
          type="error"
          @click="$emit('close')"
        />
      </div>
    </form>
  </div>
</template>

<script>
import Button from "../buttons/Button";
import { call } from "../../remote";
import Storage from "@/utils/Storage";

export default {
  name: "NewDuelDialog",
  components: {
    Button,
  },
  data() {
    return {
      duelData: {
        name: Storage.getItem(this.$config.STORAGE_KEY_DUEL_NAME) ?? "",
        description: "",
        visibility: "public",
      },
      errorMessage: "",
    };
  },
  methods: {
    async createDuel() {
      if (this.duelData.name.length < 5 || this.duelData.name.length > 30) {
        this.errorMessage = "Duel name must be between 5-30 characters";
        return;
      }

      Storage.setItem(this.$config.STORAGE_KEY_DUEL_NAME, this.duelData.name);

      try {
        let res = await call({
          path: "/match",
          method: "POST",
          body: this.duelData,
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
