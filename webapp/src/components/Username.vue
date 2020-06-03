<template>
    
    <div class="user-name">
        <span :style="'color: ' + (color ? color : 'orange')"><slot></slot></span>
    </div>
    
</template>

<script>
export default {
  name: 'username',
  props: ["color"]
}
</script>

<style scoped>
.user-name {
	text-shadow: 1px 1px #000;
}
</style>