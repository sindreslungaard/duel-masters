import { createApp } from 'vue'
import App from "./App.vue";
import router from "./router";
import { store } from "./store";

// self-invoking directives
import draggable from './directives/draggable';
import closable from './directives/closable';
import {vTooltip} from 'floating-vue'
import 'floating-vue/dist/style.css'

const myApp = createApp(App)
myApp
.use(store)
.use(router)

myApp.directive("click-outside", closable)
myApp.directive("draggable", draggable) 
myApp.directive('tooltip', vTooltip)

myApp.mount('#app')