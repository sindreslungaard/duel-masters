<template>
  <BaseDialog class="warning-dialog">
    <template #default>
      <p>{{ message }}</p>
      <Button
        text="Close"
        @click="$emit('update:modelValue')"
      />
    </template>
  </BaseDialog>
</template>

<script>
import Button from "../buttons/Button";
import BaseDialog from "@/components/dialogs/BaseDialog";

export default {
  name: "WarningDialog",
  components: {
    Button,
    BaseDialog,
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
