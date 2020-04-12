import PortalVue from 'portal-vue';
import { PluginObject, VueConstructor } from 'vue';

import { externalComponent } from '@/helpers/component-ref';

import { dashboardStore } from '@/store/dashboards';
import { pluginStore, UIPlugin } from '@/store/plugins';
import { serviceStore } from '@/store/services';
import { systemStore } from '@/store/system';

import database from '@/plugins/database';
import startup from '@/plugins/startup';
import store from '@/store';

const loadRemotePlugin = async (plugin: UIPlugin): Promise<PluginObject<any>> => {
  try {
    const obj = await externalComponent(plugin.url);
    pluginStore.commitResult({ id: plugin.id, loaded: false, error: null });
    return {
      install: (Vue, { store }) => {
        try {
          Vue.use(obj, { store });
          pluginStore.commitResult({ id: plugin.id, loaded: true, error: null });
        } catch (e) {
          pluginStore.commitResult({ id: plugin.id, loaded: false, error: e.toString() });
        }
      },
    };
  } catch (e) {
    pluginStore.commitResult({ id: plugin.id, loaded: false, error: e.toString() });
    // return a dummy so we don't have to filter the list before calling Vue.use()
    return { install: () => { } };
  }
};

export default async ({ Vue }: { Vue: VueConstructor }): Promise<void> => {
  // Make Vue accessible as global variable in umd plugins
  Object.defineProperty(window, 'Vue', { value: Vue });

  // Enable the Vue devtools performance tab
  Vue.config.performance = Boolean(process.env.DEV && process.env.BLOX_PERFORMANCE);

  Vue.use(startup);
  Vue.use(database);

  const plugins: PluginObject<any>[] = [
    PortalVue,
    history,
  ];

  try {
    Vue.$startup.onStart(() => systemStore.start());
    Vue.$startup.onStart(() => serviceStore.start());

    const remotePlugins = await pluginStore.fetchPlugins();
    const loaded = await Promise
      .all(remotePlugins.map(loadRemotePlugin));

    // Append '?safe' to the URL to disable installing remote plugins
    const url = new URL(window.location.href);
    const safeArg = url.searchParams.get('safe');

    if (safeArg === null) {
      plugins.push(...loaded);
    }
  } catch (e) { }

  plugins.forEach(plugin => Vue.use(plugin, { store }));
};
