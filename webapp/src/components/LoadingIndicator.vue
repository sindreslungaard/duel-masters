<template>
  <div class="loading-indicator">
    Loading
  </div>
</template>

<script>
export default {
  name: "LoadingIndicator",
};
</script>

<style lang="scss" scoped>
.loading-indicator {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%);

  &::after {
    content: ".";
    animation: 1.5s infinite loadingDots;
  }
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
