<template>
  <div>
    <div v-if="$store.state.user.login" class="container">
      <navigation></navigation>
      <div class="row">
        <div class="col-12 col-md-6">
          <new-offer @newOffer="newOffer"></new-offer>
          <offer-list type="own"></offer-list>
        </div>
        <offer-list class="col-12 col-md-6" type="others"></offer-list>
      </div>
    </div>
    <div v-else>
      <registration></registration>
    </div>
    <error></error>
  </div>
</template>

<script>
    import Navigation from './Navigation'
    import Registration from './Registration'
    import Error from './Error'
    import NewOffer from "./NewOffer"
    import OfferList from "./OfferList"

    export default {
        created()
        {
            this.$store.dispatch('refreshOffers')
        },
        methods: {
            newOffer(offer){

                let message = {
                    callback: "offerReceived",
                    action: 'newOffer',
                    payload: offer,
                    user: this.$store.state.user,
                };

                this.$socket.sendObj(message);
            }
        },
        components: {
            Navigation,
            Registration,
            Error,
            NewOffer,
            OfferList
        }
    }
</script>

<style scoped>
</style>