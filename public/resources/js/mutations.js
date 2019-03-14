export default {
    register(state, user)
    {
        state.user = user;
        state.actionMessage = "";
    },
    logout(state)
    {
        state.user = {};
        state.offers = [];
    },
    addNewOffer(state, offer)
    {
        offer.active = true;

        if(offer.user.token == state.user.token) {
            state.user = offer.user;
            offer.my = true
        }

        state.offers.push(offer);
    },
    offerBought(state, data)
    {
        console.log(data);
        let buyer = data.payload.buyer;
        let owner = data.payload.owner;

        if(state.user.token === buyer.token) {
            state.user = buyer
        }
        else if(state.user.token === owner.token) {
            state.user = owner
        }

        let offer = state.offers.find(offer => offer.id === data.id);

        if(offer) {
            if(offer.coins - data.coins <= 0) {
                offer.active = false;
            }
            else {

                offer.coins -= data.coins;
                console.log("OFFER", offer);
            }
        }
    },
    removeOffer(state, data, e)
    {
        let offer = state.offers.find(offer => offer.id === data.id);

        console.log('e', e);

        if(offer) {
            offer.active = false;
        }
    },
    errorReceived(state, data)
    {
        console.log(data);
        if(data.user.token === state.user.token) {
            state.actionMessage = data.error;
        }
    },
    hideErrorMessage(state, data)
    {
        state.actionMessage = "";
    },
    refreshOffers(state, offers)
    {
        offers.map(offer => {
            offer.active = true
            offer.my = offer.user.token === state.user.token;
        });

        state.offers = offers;
    },
    SOCKET_ONOPEN(state, event)
    {
        Vue.prototype.$socket = event.currentTarget;
        state.socket.isConnected = true
    },
    SOCKET_ONCLOSE(state, event)
    {
        state.socket.isConnected = false;
        console.log(event)
    },
    SOCKET_ONERROR(state, event)
    {
        console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE(state, message)
    {
        state.socket.message = message
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count)
    {
        console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state)
    {
        state.socket.reconnectError = true;
    }
}