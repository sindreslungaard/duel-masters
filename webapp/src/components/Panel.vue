<template>
  <div class="panel">
    <header class="panel__header">
      <h3>{{ title }}</h3>
    </header>

    <div class="panel__content">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "Panel",
  props: {
    title: {
      type: String,
      require: true
    }
  }
};
</script>

<style lang="scss" scoped>
.panel__content {
  overflow: auto;
  background: #2b2c31;
  border-radius: 4px;
  font-size: 14px;
  color: var(--color-text-light);
  padding: var(--spacing);
  flex-grow: 1;
}

.panel__header {
  display: flex;
  align-items: center;
}
</style>
