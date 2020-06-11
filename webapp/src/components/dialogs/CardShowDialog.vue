<template>
  <div class="card-show-dialog">
    {{ message }}
    <div class="cards">
      <img
        v-for="card in cards"
        :key="card"
        :src="`/assets/cards/all/${card}.jpg`"
      />
    </div>
    <Button text="Close" @click="$emit('close')" />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "CardShowDialog",
  props: {
    message: {
      type: String,
      required: true
    },
    cards: {
      type: Array,
      required: true
    }
  },
  components: {
    Button
  }
};
</script>

<style lang="scss" scoped>
.card-show-dialog {
  text-align: center;

  .cards {
    display: flex;
    justify-content: center;

    img {
      display: block;
      margin: 20px 10px;
      flex: 0 1 0;
      border-radius: 4%;
      width: 40vh;

      &:first-child {
        margin-left: 0;
      }

      &:last-child {
        margin-right: 0;
      }
    }
  }
}
</style>
