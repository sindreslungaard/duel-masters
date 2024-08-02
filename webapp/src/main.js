import { createApp } from 'vue'
import App from "./App.vue";
import router from "./router";
import { store } from "./store";

// self-invoking directives
import draggable from './directives/draggable';
import closable from './directives/closable';

const myApp = createApp(App)
myApp
.use(store)
.use(router)

myApp.directive("click-outside", closable)
myApp.directive("draggable", draggable) 

myApp.mount('#app')