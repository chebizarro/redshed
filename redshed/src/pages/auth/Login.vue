<template>
  <q-page class="flex flex-center">
    <q-card
      square
      style="width: 400px; padding:50px"
    >

      <q-card-section>
        <img alt="RedShed" src="~assets/redshed-md.png">

        <div class="text-h5">
          RedShed
        </div>

        <div class="text-h6">
          {{ $t('auth.login.login') }}
        </div>
      </q-card-section>

      <q-card-section>
        <EmailTextField v-model="email" />
        <PasswordTextField v-model="password" />
      </q-card-section>

      <q-card-actions>
        <q-btn
          color="primary"
          :loading="loading"
          @click="doLogin"
        >
          {{ $t('auth.login.login') }}
        </q-btn>
        <q-btn>
          {{ $t('auth.login.social.google') }}
        </q-btn>
      </q-card-actions>

      <router-link to="/password/forgot">
        <a>{{ this.$t('auth.login.password_forgot') }}</a>
      </router-link>

    </q-card>
  </q-page>

</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Action } from 'vuex-class';
import EmailTextField from '@/components/auth/EmailTextField.vue';
import PasswordTextField from '@/components/auth/PasswordTextField.vue';

@Component({
  components: {
    //Alert,
    EmailTextField,
    PasswordTextField,
  },
})
export default class Login extends Vue {
  @Action('login', { namespace: 'user' }) private login: any;

  private email: string = '';
  private password: string = '';
  private loading: boolean = false;
  private error: string = '';

  private doLogin() {
    if ((this.$refs.form as HTMLFormElement).validate()) {
      this.loading = true;

      this.login({email: this.email, password: this.password}).then(() => {
        this.loading = false;
        if (this.$route.query.next) {
          this.$router.replace({ path: this.$route.query.next as string });
        } else {
          this.$router.replace({ path: '/' });
        }
      }).catch((err: any) => {
        this.error = err.message;
        this.loading = false;
      });
    }
  }
}
</script>
