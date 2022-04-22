<template>
  <div
    :class="['user-name', muteUser ? 'allow-overflow' : '']"
    :style="hub ? 'width: 200px' : ''"
  >
    <div
      v-if="hub"
      class="user-status"
      :style="
        hub ? (hub == 'match' ? 'background: cyan' : 'background: #00FF00') : ''
      "
    >
      <div class="user-status-tooltip">
        {{ hub == "match" ? "In a duel" : "Online" }}
      </div>
    </div>
    <span
      :style="
        'color: ' + (color ? color : 'orange') + '; height: 20px; flex-grow: 1;'
      "
      ><slot></slot
    ></span>
    <div class="mute-icon-container">
      <MuteIcon
        v-if="muteUser"
        :player="muteUser"
        @toggled="$emit('muteToggled')"
      />
    </div>
  </div>
</template>

<script>
import MuteIcon from "./MuteIcon.vue";
export default {
  components: { MuteIcon },
  name: "username",
  props: ["hub", "color", "muteUser"]
};
</script>

<style scoped lang="scss">
.user-name {
  display: flex;
  align-items: center;
  text-shadow: 1px 1px #000;
  overflow: hidden;

  &.allow-overflow {
    overflow: visible;
  }

  .mute-icon-container {
    display: flex;
    align-items: center;
    margin-left: 5px;
    opacity: 0;
  }

  &:hover .mute-icon-container {
    opacity: 1;
  }
}

.user-status {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  margin-right: 5px;
}

.user-status-tooltip {
  position: absolute;
  background: #333;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 2px;
  width: 60px;
  text-align: center;
  transition: opacity 0.3s;
  font-size: 12px;
  margin: 15px 15px;
  visibility: hidden;
  opacity: 0;
}

.user-status:hover .user-status-tooltip {
  visibility: visible;
  opacity: 1;
}
</style>
