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
    offerRemoved(context, offer)
    {
        this.commit('removeOffer', offer.payload, [1,2,3]);
        console.log("remove offer", offer)
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