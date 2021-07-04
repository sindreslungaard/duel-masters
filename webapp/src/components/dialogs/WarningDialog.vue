<template>
  <div class="warning-dialog">
    <p>{{ message }}</p>
    <Button text="Close" @click="$emit('close')" />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "WarningDialog",
  props: {
    message: {
      type: String,
      required: true
    }
  },
  components: {
    Button
  }
};
</script>

<style lang="scss" scoped>
.button {
  margin-top: var(--spacing);
}
</style>
