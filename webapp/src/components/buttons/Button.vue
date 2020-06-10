<template>
  <button @click="$emit('click')" class="button">{{ text }}</button>
</template>

<script>
export default {
  name: "Button",
  props: {
    text: {
      type: String,
      required: true
    }
  }
};
</script>

<style lang="scss" scoped>
.button {
  display: inline-block;
  background-color: #7289da;
  color: #e3e3e5;
  font-size: 14px;
  line-height: 20px;
  padding: 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
  cursor: pointer;
  border: none;

  &:hover {
    background-color: #677bc4;
  }
}
</style>
