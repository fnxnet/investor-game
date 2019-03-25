import VuexPersistence from 'vuex-persist';
import actions from "./actions";
import mutations from "./mutations";

Vue.use(Vuex);

const vuexLocal = new VuexPersistence({
    key: 'vuex-store'
});

export const store = new Vuex.Store({
    state: {
        url_prefix: "",
        user: {},
        offers: {},
        stats: {
            offers: [],
            sum: 0,
            max: 0,
            avg: 0,
            user_count: 0,
        },
        socket: {
            isConnected: false,
            message: '',
            reconnectError: false,
        },
        actionMessage: ""
    },
    mutations,
    actions,
    plugins: [vuexLocal.plugin]
});