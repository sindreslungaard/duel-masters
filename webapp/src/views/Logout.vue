<template>
  <div></div>
</template>


<script>
export default {
  name: "logout",
  created() {
    localStorage.clear();
    this.$router.push("login");
  }
};
</script>
