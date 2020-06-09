import Vue from "vue";
import Router from "vue-router";
import axios from "axios";

Vue.use(Router);

const router = new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [

    {
      path: "/",
      name: "index",
      component: () => import("./views/Login.vue"),
      meta: { noauth: true }
    },

    {
      path: "/login",
      name: "login",
      component: () => import("./views/Login.vue"),
      meta: { noauth: true }
    },

    {
      path: "/register",
      name: "register",
      component: () => import("./views/Register.vue"),
      meta: { noauth: true }
    },

    {
      path: "/logout",
      name: "logout",
      component: () => import("./views/Logout.vue"),
      meta: { auth: true }
    },

    {
      path: "/overview",
      name: "overview",
      component: () => import("./views/Overview.vue"),
      meta: { auth: true }
    },

    {
      path: "/decks",
      name: "decks",
      component: () => import("./views/Decks.vue"),
      meta: { auth: true }
    },

    {
      path: "/deck/:uid",
      name: "deck",
      component: () => import("./views/Deck.vue")
    },

    {
      path: "/duel/:id",
      name: "duel",
      component: () => import("./views/Match.vue"),
      meta: { auth: true }
    },

    {
      path: "/duel/:id/:invite",
      name: "duel_invite",
      component: () => import("./views/Match.vue"),
      meta: { auth: true }
    }

  ]
});

router.beforeEach((to, from, next) => {

  if(to.matched.length < 1) {
    return next("/");
  }

  const hasToken = localStorage.getItem("token") ? true : false;

  if(to.matched.some(record => record.meta.auth) && !hasToken)  {
    return next("/login?redirect_to=" + encodeURIComponent(to.fullPath));
  }

  if(to.matched.some(record => record.meta.noauth) && hasToken)  {
    return next("/overview");
  }

  return next();

});

export default router;