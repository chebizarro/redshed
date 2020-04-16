declare namespace NodeJS {
  interface ProcessEnv {
    NODE_ENV: string;
    VUE_ROUTER_MODE: 'hash' | 'history' | 'abstract' | undefined;
    VUE_ROUTER_BASE: string | undefined;
    AUTH_TOKEN: string | 'auth-token';
    API_ENDPOINT: string | undefined;
    GOOGLE_KEY: string | undefined;
  }
}
