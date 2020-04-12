import Vue from 'vue';
import VueRouter from 'vue-router';
import store from "@/store";
import jwtDecode from 'jwt-decode';

import { AUTH_TOKEN } from '@/constants';
import { IJWTDecoded } from '@/store/user/user';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {path: '/', name: 'home', component: () => import(/* webpackChunkName: "home" */ '@/layouts/MainLayout.vue')},
    {path: '/login', name: 'login', component: () => import(/* webpackChunkName: "login" */ '@/views/Login.vue')},
    {path: '/signup', name: 'signup', component: () => import(/* webpackChunkName: "signup" */ '@/views/Signup.vue')},
    {path: '/welcome', name: 'welcome', component: () => import(/* webpackChunkName: "welcome" */ '@/views/Welcome.vue')},
    {path: '/verify/:id/:token', name: 'verify', component: () => import(/* webpackChunkName: "verify" */ '@/views/Verify.vue')},
    {path: '/forgot-password', name: 'forgot-password', component: () => import(/* webpackChunkName: "forgot-password" */ '@/views/ForgotPassword.vue')},
    {path: '/reset-password/:id/:token', name: 'reset-password', component: () => import(/* webpackChunkName: "reset-password" */ '@/views/ResetPassword.vue')},
    {path: '/profile', name: 'profile', component: () => import(/* webpackChunkName: "profile" */ '@/views/Profile.vue'), meta: { requiresAuth: true }},

    {path: '/404', alias: '*', name: 'notfound', component: () => import(/* webpackChunkName: "notfound" */ '@/views/NotFound.vue')},
  ],
});

router.beforeEach((to, from, next) => {
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    const token = localStorage.getItem(AUTH_TOKEN);

    if (token === null) {
      next({
        path: '/login',
        query: { next: to.fullPath },
      });
    } else {
      const decoded: IJWTDecoded = jwtDecode(token);

      if (decoded.exp * 1000 < Date.now().valueOf()) {
        next({
          path: '/login',
          query: { next: to.fullPath },
        });
      } else {
        next();
      }
    }
  } else {
    next();
  }
});

export default router;
