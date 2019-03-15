import axios from 'axios'

export default {
    offerReceived(context, offer)
    {
        console.log(offer)
        this.commit('addNewOffer', offer.payload);
    },
    offerBought(context, offer)
    {
        this.commit('offerBought', offer.payload);
    },
    offerRemoved(context, data)
    {
        this.commit('removeOffer', data.payload);
    },
    error(context, data)
    {
        this.commit('errorReceived', data);
    },
    refreshOffers({state})
    {
        axios.get(state.url_prefix + '/offers').then(response => {
            this.commit('refreshOffers', response.data)
        });
    }
}