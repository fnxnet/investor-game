import VueNativeSock from 'vue-native-websocket'
import Admin from '../vue/Admin';
import {store} from './store';

Vue.config.devtools = true;
let options = {
    format: 'json',
    store,
    reconnection: true,
    reconnectionAttempts: 5,
    reconnectionDelay: 3000,
};

Vue.use(VueNativeSock, 'ws://localhost:3030/ws', options);

const app = new Vue({
    el: "#app",
    store,
    template: '<admin/>',
    components: {Admin}
});
