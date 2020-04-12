// import something here
import enUS from '../i18n/en-us/auth'

// "async" is optional
export default async ({ app }) => {
  app.i18n.mergeLocaleMessage('en-us', enUS)
}
