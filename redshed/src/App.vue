<template>
  <div id="q-app">
    <router-view />
  </div>
</template>

<script lang="ts">


import { Vue } from 'vue-property-decorator';
import { Getter } from 'vuex-class';

export default class App extends Vue {

  @Getter('isAuthenticated', { namespace: 'user' }) private isAuthenticated!: boolean;

  private drawer: boolean = false;
}
</script>
