<template>
  <div class="message-dialog">
    <p>{{ message }}</p>
    <Button text="Close" @click="$emit('close')" />
  </div>
</template>

<script>
import Button from "../buttons/Button";

export default {
  name: "MessageDialog",
  props: {
    message: {
      type: String,
      required: true
    },
  },
  components: {
    Button
  }
};
</script>

<style lang="scss" scoped>
.message-dialog {
  text-align: center;
  background: #222428;
}

.message-dialog > *:first-child{
  padding-top: 20px;
}

.message-dialog > *:last-child{
  margin-bottom: 10px;
}

.message-dialog > button {
  padding-top: 5px;
  padding-bottom: 5px;
  padding-left: 15px;
  padding-right: 15px;
}

</style>
