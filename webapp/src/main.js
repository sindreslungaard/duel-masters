import VueFinalModal from "vue-final-modal";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "./scss/style.scss";
import config from "@/config";
import Draggable from "@/directives/Draggable";

import { createApp } from "vue";

// self-invoking directives
//import "./directives/draggable";

const app = createApp({
  router,
  ...App,
});
app.config.productionTip = false;

app.use(store);
app.use(router);
app.use(config);
app.use(VueFinalModal());

app.directive("draggable", Draggable);

app.mount("#app");

/*app.use(VModal, {
  dynamic: true,
  dynamicDefaults: { height: "auto", scrollable: true, width: "300px" },
});*/
