import Vue from 'vue';

import { RedShedStartup } from "@/plugins/startup";
import { RedShedDatabase } from "@/plugins/database";

declare module 'vue/types/vue' {

  interface Vue {
    $dense: boolean;
  }
  interface VueConstructor {
    $database: RedShedDatabase;
    $startup: RedShedStartup;
  }
}
