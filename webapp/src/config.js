import config from "./config/Config";

export default {
  install(app, options) {
    app.config.globalProperties.$config = Object.freeze(config);
  },
};
