import enUS from '../i18n/en-us/auth'
import GAuth from 'vue-google-oauth2'
import Vue from 'vue'

// "async" is optional
export default async ({ app }) => {
  app.i18n.mergeLocaleMessage('en-us', enUS);

  Vue.use(GAuth, {
    clientId: process.env.GOOGLE_KEY,
    scope: 'email',
    prompt: 'consent',
    fetch_basic_profile: true
  });

}
