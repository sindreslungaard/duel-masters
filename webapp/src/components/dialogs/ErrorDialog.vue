<template>
  <div class="error-dialog">
    <p>{{ message }}</p>
    <Button
      text="Reload"
      @click="reloadPage"
    />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "ErrorDialog",
  components: {
    Button,
  },
  props: {
    message: {
      type: String,
      required: true,
    },
  },
  methods: {
    reloadPage() {
      window.location.reload();
    },
  },
};
</script>

<style lang="scss" scoped>
.button {
  margin-top: var(--spacing);
}
</style>
