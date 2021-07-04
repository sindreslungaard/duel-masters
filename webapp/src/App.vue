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

html, body {
  height: 100%;
}

body {
  background: #36393f;
  margin: 0;
  font-family: "Roboto", sans-serif;
  color: #eee;
}

.vm--modal {
  background-color: #36393f;
  border-color: #222428;
  border-radius: var(--border-radius);
  padding: 20px;
  border: 1px solid #666;
}

.vm--overlay {
  background: rgba(0, 0, 0, 0.4);
}
</style>
