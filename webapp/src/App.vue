<template>
  <div id="app">
    <router-view />
    <modals-container />
  </div>
</template>

<style>
@import url(
  https://fonts.googleapis.com/css?family=Roboto:300,
  400,
  500,
  600,
  700
);

body {
  background: #36393f;
  margin: 0;
  font-family: "Roboto", sans-serif;
  color: #eee;
}

.vm--modal {
  background-color: #222428;
  border-radius: 4px;
  padding: 20px;
}
</style>
