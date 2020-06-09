<template>
  <div class="user-name" :style="hub ? 'width: 200px' : ''">
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
    <span :style="'color: ' + (color ? color : 'orange')"><slot></slot></span>
  </div>
</template>

<script>
export default {
  name: "username",
  props: ["hub", "color"]
};
</script>

<style scoped>
.user-name {
  text-shadow: 1px 1px #000;
  overflow: hidden;
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
