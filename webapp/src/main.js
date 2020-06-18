import Vue from "vue";
import VModal from "vue-js-modal";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "./scss/style.scss";
import config from "@/config/Config";

Vue.config.productionTip = false;

Vue.use(VModal, { dynamic: true, dynamicDefaults: { height: "auto", scrollable: true, width: "300px" } });

Vue.prototype.$config = Object.freeze(config);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
