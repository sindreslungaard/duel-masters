<template>
  <span v-if="props.uid">
    <div 
      class="card-preview" 
      :style="{
        'top': previewCardTop + 'px', 
        'left': previewCardLeft + 'px',
      }" 
    >
      <img
        :src="`https://scans.shobu.io/${uid}.jpg`"
        alt="Full card"
      />
    </div>
  </span>
</template>

<script setup>
import { watchEffect } from 'vue'

const props = defineProps(['uid', 'event', 'xPos', 'side']);
let previewCardTop = 0;
let previewCardLeft = 0;

watchEffect(() => {
  if (!props.uid) {
    return
  }  
  let positionY = props.event.clientY - 420/2;
  let positionX = props.xPos ? props.xPos : (props.event.clientX -  props.event.offsetX + props.event.target.offsetWidth);
  if (props.side == "right") {
    positionX += 5
  }  
  if (props.side == "left") {
    positionX -= 300
  }
  if (positionY + 420 > window.innerHeight) {
    positionY = window.innerHeight - 420;
  }
  if (positionY < 10) {
    positionY = 10;
  }
  previewCardTop = positionY;
  previewCardLeft = positionX;
})
</script>

<style scoped lang="scss">
.card-preview {
  width: 300px;
  text-align: center;
  border-radius: 4px;
  height: 420px;
  z-index: 2005;
  position: absolute;
}

.card-preview > img {
  width: 300px;
  border-radius: 15px;
}
</style>
