<template>
  <InternalLink target="_blank" rel="noopener noreferrer" :href="href">
    <slot></slot>
  </InternalLink>
</template>

<script>
import InternalLink from "./InternalLink";

export default {
  name: "ExternalLink",
  components: {
    InternalLink
  },
  props: {
    href: {
      type: String,
      required: true
    }
  }
};
</script>
