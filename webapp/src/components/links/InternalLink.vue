<template>
  <a :target="target" :rel="rel" :href="href"><slot></slot></a>
</template>

<script>
export default {
  name: "InternalLink",
  props: {
    href: {
      type: String,
      required: true
    },
    target: {
      type: String,
      required: true
    },
    rel: {
      type: String,
      required: true
    }
  }
};
</script>

<style lang="scss" scoped>
a {
  color: #fff;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

img {
  display: block;
}
</style>
