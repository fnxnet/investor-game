import axios from 'axios'

export default {
    offerReceived(context, offer)
    {
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
    userRegistered(context, data)
    {
        this.commit('increaseUserCount');
    },
    incomeReceived(context, data)
    {
        this.commit('updateUser', data.user);
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