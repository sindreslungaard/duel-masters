<template>
  <div />
</template>


<script>
export default {
  name: "Logout",
  created() {
    localStorage.clear();
    this.$router.push("login");
  },
};
</script>
