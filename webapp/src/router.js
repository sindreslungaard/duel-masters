import { createRouter, createWebHistory } from "vue-router";
import { call } from "./remote";
import { store } from "./store";
import { getSettings } from "./helpers/settings";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
      path: "/recover-password",
      name: "recover_password",
      component: () => import("./views/RecoverPassword.vue"),
      meta: { noauth: true }
    },

    {
      path: "/recover-password/:code",
      name: "reset_password",
      component: () => import("./views/ResetPassword.vue"),
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
      component: () => {
        if (getSettings().legacyDeckBuilder) {
          return import("./views/Decks.vue")
        }  
        return import("./views/DecksNew.vue")
      },
      meta: { auth: true }
    },

    {
      path: "/deck/:uid",
      name: "deck",
      component: () => import("./views/Deck.vue")
    },

    {
      path: "/settings",
      name: "settings",
      component: () => import("./views/Settings.vue"),
      meta: { auth: true }
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
  call({
    path: `/preferences`,
    method: "GET"
  })
    .then(res => {
      store.preferences = res.data;
    })
    .catch(err => {
      if (err && err.response && err.response.data) {
        console.error(err.response.data.error);
      } else {
        console.error(e);
      }
    });

  if (to.matched.length < 1) {
    return next("/");
  }

  const hasToken = localStorage.getItem("token") ? true : false;

  if (to.matched.some(record => record.meta.auth) && !hasToken) {
    return next("/login?redirect_to=" + encodeURIComponent(to.fullPath));
  }

  if (to.matched.some(record => record.meta.noauth) && hasToken) {
    return next("/overview");
  }

  return next();
});

export default router;
