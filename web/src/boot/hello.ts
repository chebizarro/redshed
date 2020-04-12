import * as hello from 'hellojs'

export default ({ Vue }) => {
  hello.init({
    google: ''
  });
  Vue.prototype.$hello = hello;
}
