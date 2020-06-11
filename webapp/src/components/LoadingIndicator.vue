<template>
  <div class="loading-indicator">
    Loading
  </div>
</template>

<script>
export default {
  name: "LoadingIndicator"
};
</script>

<style scoped>
div::after {
  content: ".";
  animation: 1.5s infinite loadingDots;
}

@keyframes loadingDots {
  from {
    content: ".";
  }

  25% {
    content: "..";
  }

  50% {
    content: "...";
  }

  75% {
    content: "....";
  }

  to {
    content: ".";
  }
}
</style>
