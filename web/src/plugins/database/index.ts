import { VueConstructor } from 'vue';

import { RedShedDatabaseImpl, checkDatastore } from './database';

export * from './types';

export default {
  install(Vue: VueConstructor) {
    Vue.$database = new RedShedDatabaseImpl();
    Vue.$startup.onStart(checkDatastore);
  },
};
