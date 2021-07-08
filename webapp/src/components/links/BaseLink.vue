<template>
  <a
    :target="target"
    :rel="rel"
    :href="href"
  ><slot /></a>
</template>

<script>
export default {
  name: "BaseLink",
  props: {
    href: {
      type: String,
      required: true,
    },
    target: {
      type: String,
      required: false,
      default: "",
    },
    rel: {
      type: String,
      required: false,
      default: "",
    },
  },
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
