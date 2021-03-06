import VueNativeSock from 'vue-native-websocket'
import App from '../vue/App';
import {store} from './store';

Vue.config.devtools = true;
let options = {
    format: 'json',
    store,
    reconnection: true,
    reconnectionAttempts: 5,
    reconnectionDelay: 3000,
};

let url = 'ws://' + window.location.host.replace('admin.', '') + '/ws'

Vue.use(VueNativeSock, url, options);

const app = new Vue({
    el: "#app",
    store,
    template: '<app/>',
    components: {App}
});
