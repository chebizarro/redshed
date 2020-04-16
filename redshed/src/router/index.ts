import { route } from 'quasar/wrappers'
import VueRouter from 'vue-router'
import { StoreInterface } from '@/store'
import routes from './routes'
import {IJWTDecoded} from "@/store/auth/user";
import jwtDecode from 'jwt-decode';
/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation
 */

export default route<StoreInterface>(function ({ Vue }) {
  Vue.use(VueRouter);

  const Router = new VueRouter({
    scrollBehavior: () => ({ x: 0, y: 0 }),
    routes,

    // Leave these as is and change from quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    mode: process.env.VUE_ROUTER_MODE,
    base: process.env.VUE_ROUTER_BASE,

  });

  Router.beforeEach((to, from, next) => {

    if (to.matched.some((record) => record.meta.requiresAuth)) {
      const token = localStorage.getItem(process.env.AUTH_TOKEN);

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

  return Router
})
