import VueNativeSock from 'vue-native-websocket'
import Panel from '../vue/Panel';
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

Vue.filter('toCurrency', function (value) {
    if(typeof value !== "number") {
        return value;
    }
    var formatter = new Intl.NumberFormat('pl-PL', {
        style: 'currency',
        currency: 'PLN',
        minimumFractionDigits: 2
    });
    return formatter.format(value);
});

const app = new Vue({
    el: "#app",
    store,
    template: '<panel/>',
    components: {Panel}
});
