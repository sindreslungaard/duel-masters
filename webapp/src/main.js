import Vue from "vue";
import VModal from "vue-js-modal";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// self-invoking directives
import "./directives/draggable";

Vue.config.productionTip = false;

Vue.use(VModal, { dynamic: true, dynamicDefaults: { height: "auto", scrollable: true } });

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
