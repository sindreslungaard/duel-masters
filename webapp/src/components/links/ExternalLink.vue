<template>
  <BaseLink
    target="_blank"
    rel="noopener noreferrer"
    :href="href"
  >
    <slot />
  </BaseLink>
</template>

<script>
import BaseLink from "./BaseLink";

export default {
  name: "ExternalLink",
  components: {
    BaseLink,
  },
  props: {
    href: {
      type: String,
      required: true,
    },
  },
};
</script>
