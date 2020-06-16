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
      required: true
    },
    submit: {
      type: Boolean,
      required: false,
      default: false
    },
    type: {
      type: String,
      required: false,
      default: "primary",
      validator(value) {
        return ["primary", "secondary"].includes(value);
      }
    }
  },
  computed: {
    buttonClass() {
      return `button--${this.type}`;
    },
    buttonType() {
      return this.submit ? "submit" : "button";
    }
  }
};
</script>

<style lang="scss" scoped>
.button {
  display: inline-block;
  color: #e3e3e5;
  font-size: 14px;
  line-height: 20px;
  padding: 7px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
  cursor: pointer;
  border: none;

  &--primary {
    background-color: var(--primary);

    &:hover {
      background-color: #677bc4;
    }
  }

  &--secondary {
    background-color: var(--secondary);

    &:hover {
      background-color: #ed3e3e;
    }
  }
}
</style>
