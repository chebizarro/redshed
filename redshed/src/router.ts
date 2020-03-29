import Vue from 'vue';
import VueRouter from 'vue-router';
import store from "@/store";

Vue.use(VueRouter);

const router = new VueRouter({
  routes: [
    {
      path: '/',
      component: () => import('layouts/MainLayout.vue'),
      children: [
        {
          path: '', component: () => import('pages/Index.vue')
        },
        {
          path: '/login',
          component: () => import('pages/Login.vue')
        },
        {
          path: '/register',
          component: () => import('pages/Register.vue')
        },
      ],
      beforeEnter: (to, from, next) => {
        let auth = localStorage.getItem('token')
        if (!auth) {
          store.commit('auth/logout')
          next('/login')
        } else {
          Vue.http.get('auth/user')
            .then(response => {
              next()
            }, response => {
              next('/login')
            })
        }
      },

    },
    { // Always leave this as last one
      path: '*',
      component: () => import('pages/Error404.vue')
    },
  ],
  mode: 'history',
  base: '/',
});

export default router;
