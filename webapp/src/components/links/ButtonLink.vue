<template>
  <BaseLink class="button-link" :class="buttonClass" :href="href">{{
    text
  }}</BaseLink>
</template>

<script>
import BaseLink from "./BaseLink";

export default {
  name: "ButtonLink",
  components: {
    BaseLink
  },
  props: {
    text: {
      type: String,
      required: true
    },
    href: {
      type: String,
      required: true
    },
    type: {
      type: String,
      required: false,
      default: "info",
      validator(value) {
        return ["info", "error", "success"].includes(value);
      }
    }
  },
  computed: {
    buttonClass() {
      return `button-link--${this.type}`;
    }
  }
};
</script>

<style lang="scss" scoped>
.button-link {
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

  &:hover {
    text-decoration: none;
  }

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
