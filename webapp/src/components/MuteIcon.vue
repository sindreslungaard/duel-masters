<template>
  <div class="icon-container">
    <img
      class="mute-icon"
      @click="toggleMute(player)"
      width="20px"
      height="20px"
      :src="`/assets/images/volume_${!isMuted ? 'off' : 'on'}.png`"
    />
    <div style="position: relative;">
      <div class="tooltip">
        {{ isMuted ? "Unmute" : "Mute" }} player
      </div>
    </div>
  </div>
</template>

<script>
import { toggleMute, isMuted } from "../helpers/settings";

export default {
  props: ["player"],
  name: "muteicon",
  methods: {
    toggleMute() {
      this.isMuted = toggleMute(this.player);
      this.$emit("toggled");
    },
    updateIsMuted() {
      const newIsMuted = isMuted(this.player);
      this.isMuted = newIsMuted;

      if (newIsMuted !== this.isMuted) {
        this.$emit("toggled");
      }
    }
  },

  data() {
    return {
      isMuted: isMuted(this.player)
    };
  },

  created() {
    addEventListener("storage", this.updateIsMuted);
    addEventListener("storageUpdated", this.updateIsMuted);
  },

  beforeDestroy() {
    removeEventListener("storage", this.updateIsMuted);
    removeEventListener("storageUpdated", this.updateIsMuted);
  }
};
</script>

<style scoped lang="scss">
.icon-container {
  height: 20px;
  display: flex;

  .tooltip {
    line-height: 16px;
    position: absolute;
    right: 25px;
    bottom: 0px;
    background: #333;
    border: 1px solid #eee;
    border-radius: 4px;
    width: 80px;
    padding: 1px;
    text-align: center;
    transition: opacity 0.3s;
    font-size: 12px;
    visibility: hidden;
    opacity: 0;
  }

  &:hover .tooltip {
    visibility: visible;
    opacity: 1;
  }

  .mute-icon {
    color: #eee;
    -webkit-filter: invert(70%);
    filter: invert(70%);

    &:hover {
      cursor: pointer;
    }
  }
}
</style>
