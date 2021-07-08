<template>
  <div class="panel">
    <header class="panel__header">
      <h3>{{ title }}</h3>
    </header>

    <div class="panel__content">
      <slot />
    </div>
  </div>
</template>

<script>
export default {
  name: "Panel",
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.panel {
  position: relative;
  display: flex;
  flex-direction: column;

  &__header {
    display: flex;
    align-items: center;
  }

  &__content {
    overflow: auto;
    background-color: var(--color-foreground);
    border-radius: var(--border-radius);
    font-size: 14px;
    color: var(--color-text-light);
    padding: var(--spacing);
    flex-grow: 1;
  }
}
</style>
