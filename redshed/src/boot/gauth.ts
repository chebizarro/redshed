import { boot } from 'quasar/wrappers'
import Vue from 'vue';
import GAuth from 'vue-google-oauth2';

declare module 'vue/types/vue' {
  interface Vue {
    $gAuth: any;
  }
}

Vue.use(GAuth,
  {
    clientId: process.env.GOOGLE_KEY,
    scope: 'email',
    prompt: 'consent',
    fetch_basic_profile: true
  });


export default boot(async ({ app, router, store, Vue }) => {
  //app.$gAuth = GAuth();
})

