<template>
  <div class="error-dialog">
    {{ message }}
    <Button text="Reload" @click="reloadPage" />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "ErrorDialog",
  props: {
    message: {
      type: String,
      required: true
    }
  },
  components: {
    Button
  },
  methods: {
    reloadPage() {
      window.location.reload();
    }
  }
};
</script>

<style lang="scss" scoped>
.button {
  margin-top: var(--spacing);
}
</style>
