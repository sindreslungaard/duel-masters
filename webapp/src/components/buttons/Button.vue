<template>
  <button
    class="button"
    :class="buttonClass"
    :type="buttonType"
    @click="$emit('click')"
  >
    {{ text }}
  </button>
</template>

<script>
export default {
  name: "Button",
  props: {
    text: {
      type: String,
      required: true,
    },
    submit: {
      type: Boolean,
      required: false,
      default: false,
    },
    type: {
      type: String,
      required: false,
      default: "info",
      validator(value) {
        return ["info", "error", "success"].includes(value);
      },
    },
  },
  computed: {
    buttonClass() {
      return `button--${this.type}`;
    },
    buttonType() {
      return this.submit ? "submit" : "button";
    },
  },
};
</script>

<style lang="scss" scoped>
.button {
  display: inline-block;
  color: #e3e3e5;
  font-size: 14px;
  line-height: 20px;
  padding: 7px 10px;
  border-radius: var(--border-radius);
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
  cursor: pointer;
  border: none;
  font-weight: bold;

  &--info {
    background-color: var(--info);

    &:hover {
      background-color: var(--info-hover);
    }
  }

  &--error {
    background-color: var(--error);

    &:hover {
      background-color: var(--error-hover);
    }
  }

  &--success {
    background-color: var(--success);

    &:hover {
      background-color: var(--success-hover);
    }
  }
}
</style>
