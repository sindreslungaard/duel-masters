<template>
  <div class="warning-dialog">
    <p>{{ message }}</p>
    <Button
      text="Close"
      @click="$emit('close')"
    />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "WarningDialog",
  components: {
    Button,
  },
  props: {
    message: {
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.button {
  margin-top: var(--spacing);
}
</style>
