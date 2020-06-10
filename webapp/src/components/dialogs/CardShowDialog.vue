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

    img {
      display: block;
      margin: 20px auto;
      padding: 0 10px;
      flex: 1 1 0;
      width: 0;

      &:first-child {
        padding-left: 0;
      }

      &:last-child {
        padding-right: 0;
      }
    }
  }
}
</style>
