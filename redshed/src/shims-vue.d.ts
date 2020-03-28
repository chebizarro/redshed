import Vue from 'vue';
import { RedShedStartup } from "./plugins/startup";

declare module 'vue/types/vue' {
  interface Vue {
    $dense: boolean;
  }
  interface VueConstructor {
    $startup: RedShedStartup;
  }
}
