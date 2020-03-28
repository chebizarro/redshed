<template>
  <div
    class="login-form white"
    color="primary"
  >
    <v-alert
      class="login-form__fail-alert px-5"
      type="error"
      :value="loginFailed"
      dismissible
    >
      {{ $t('global.login.failed') }}
    </v-alert>
    <div
      class="
        login-form__wrapper
        d-flex
        flex-column
        align-center
        justify-center
      "
    >
      <v-progress-circular
        v-if="loginWait"
        :size="100"
        :width="5"
        class="login-form__loader"
        color="primary"
        indeterminate
      />

      <template v-else>

        <!-- logo -->
        <img
          v-if="showLogo"
          class="mb-2"
          :src="require(`@/assets/images/${logo}`)"
        >
        <!-- app title -->
        <h1
          class="
            login-form__title
            text-center
            primary--text
            display-1
            font-weight-bold
            mb-10
          "
        >
          {{ $t('global.login.title') }}
        </h1>

        <!-- locale select -->
        <v-menu v-if="localeSelectable">
          <template v-slot:activator="{ on }">
            <v-btn
              v-on="on"
              dark fab small
              color="secondary"
              class="mb-2"
            >
              <v-icon>translate</v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              v-for="(locale, i) in locales"
              :key="i"
              @click="changeLocale(locale.name)"
            >
              <v-list-item-title>{{ locale.text }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

        <!-- login form -->
        <v-form
          v-model="valid"
          class="login-form__form"
          ref="form"
          lazy-validation
          @submit.prevent
        >
          <v-text-field
            :label="$t('global.login.login')"
            v-model="user"
            :rules="loginRules"
            required
          ></v-text-field>
          <v-text-field
            :label="$t('global.login.password')"
            v-model="password"
            :rules="passwordRules"
            :counter="30"
            required
            :append-icon="passAppendIcon"
            @click:append="() => (passwordHidden = !passwordHidden)"
            :type="passTextFieldType"
            class="mb-5"
          ></v-text-field>
          <v-btn
            block
            @click="loginAttempt()"
            :disabled="!valid"
            class="primary white--text"
          >
            {{ $t('global.login.submit') }}
          </v-btn>

          <div class="login-form__fancy">
          OR
          </div>

          <v-btn
            block
            :to="{ path: 'register' }"
            class="primary white--text"
          >
            {{ $t('global.login.signup') }}
          </v-btn>
        </v-form>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import { Auth } from '../../../config/auth'
import { Prop, Vue } from 'vue-property-decorator'
import { Mutation } from 'vuex-module-decorators'
import { namespace } from 'vuex-class'
const auth = namespace('auth')
import Component from 'vue-class-component'

@Component
export default class LoginForm extends Vue {
  private valid: boolean = false
  private password: string
  private user: string
  private passwordHidden: boolean = true
  private alphanumericRegex: RegExp = /^[a-zA-Z0-9]+$/
  private emailRegex: RegExp = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/

  @Prop(String) readonly redirect: string = '/'
  @Prop(Boolean) readonly showLogo: boolean = true
  @Prop(String) readonly logo: string = 'vue-crud-sm.png'
  @Prop(Boolean) readonly localeSelectable: boolean = true


  get loginRegex(): RegExp {
    return Auth.loginRegex ? Auth.loginRegex : (Auth.loginWithEmail ? this.emailRegex : this.alphanumericRegex)
  }

  get loginRules() {
    return [
      v => !!v || this.$t('global.login.loginRequired'),
      v => this.emailRegex.test(v) || this.$t('global.login.incorrectLogin'),
    ]
  }

  get passwordRegex(): RegExp {
    return Auth.passwordRegex ? Auth.passwordRegex : this.alphanumericRegex
  }

  get passwordRules() {
    return [
      v => !!v || this.$t('global.login.passwordRequired'),
      v => this.passwordRegex.test(v) || this.$t('global.login.incorrectPassword')
    ]
  }

  get credential(): {} {
    let credentials = {}
    credentials[Auth.loginFieldName || 'login'] = this.user
    credentials[Auth.passwordFieldName || 'password'] = this.password
    return credentials
  }

  get passTextFieldType(): string {
    return this.passwordHidden ? 'password' : 'text'
  }

  get passAppendIcon(): string {
    return this.passwordHidden ? 'visibility' : 'visibility_off'
  }

  public changeLocale(locale: string) {
    this.$i18n.locale = locale
    this.$vuetify.lang.current = locale
    this.setLocale(locale)
  }

  public loginAttempt(): void {
    this.login(this.credential).then(() => {
      this.$router.push({ path: this.redirect })
    })
  }

  @auth.State
  public loginWait!: boolean

  @auth.State
  public loginFailed!: boolean

  @auth.Action
  public login!

  @Mutation
  public setLocale!: (newLocale: string) => void

}

</script>

<style lang="scss" scoped>
.login-form {
  &__fail-alert {
    width:100%;
    position:absolute;
    top: 0;
    left:0;
  }
  &__wrapper {
    height: 100vh;
    width: 100%;
  }
  &__form {
    width: 300px;
  }
  &__logo {
    width:100%;
    height:auto;
  }
  &__fancy {
     display: flex;
     width: 100%;
     justify-content: center;
     align-items: center;
     text-align: center;
     margin:25px 0;
  }
  &__fancy:before, &__fancy:after {
    content: '';
    border-top: 1px solid #ccc;
    margin: 0 10px 0 50px;
    flex: 1 0 20px;
  }
  &__fancy:after {
    margin: 0 50px 0 10px;
  }
}
</style>
