<template>
  <div id="app">
    <router-view />
    <modals-container />
  </div>
</template>

<style>
/* prettier-ignore */
@import url(https://fonts.googleapis.com/css?family=Roboto:300,400,500,600,700);

html {
  min-height: 100%;
  background: rgb(29,33,42);
  background: linear-gradient(45deg, rgba(29,33,42,1) 0%, rgba(20,16,21,1) 100%);
}

body {
  margin: 0;
  font-family: "Roboto", sans-serif;
  color: #eee;
}

.vm--modal {
  background-color: #36393f;
  border-color: #222428;
  border-radius: 4px;
  padding: 20px;
}

*::-webkit-scrollbar-track {
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  background-color: #484c52;
}

*::-webkit-scrollbar {
  width: 6px;
  height: 6px;
  background-color: #484c52;
}

*::-webkit-scrollbar-thumb {
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  background-color: #222;
}
</style>
