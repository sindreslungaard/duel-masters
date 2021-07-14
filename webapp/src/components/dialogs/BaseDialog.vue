<template>
  <vue-final-modal
    v-slot="{ close }"
    v-bind="$attrs"
    classes="modal-container"
    content-class="modal-content"
    @click-outside="$emit('close', close)"
  >
    <h2 class="modal__title">
      <slot name="title" />
    </h2>
    <div class="modal__content">
      <slot />
    </div>
  </vue-final-modal>
</template>
<script>
export default {
  name: "BaseDialog",
  inheritAttrs: false,
  emits: ["close"],
};
</script>
<style scoped>
::v-deep(.modal-container) {
  display: flex;
  justify-content: center;
  align-items: center;
}
::v-deep(.modal-content) {
  position: relative;
  display: flex;
  flex-direction: column;
  max-height: 90%;
  margin: 0 1rem;
  padding: var(--spacing);
  border-radius: var(--border-radius);
  background-color: var(--color-background);
  min-width: 300px;
}
.modal__title {
  margin-bottom: var(--spacing);
}
.modal__content {
  flex-grow: 1;
  overflow-y: auto;
}
</style>
