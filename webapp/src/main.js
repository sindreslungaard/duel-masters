import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// self-invoking directives
import "./directives/draggable";
import "./directives/closable";

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
