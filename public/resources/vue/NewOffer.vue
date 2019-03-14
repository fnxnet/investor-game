<template>
  <div class="offer">
    <div class="input-group mb-4">
      <input type="number" class="form-control" placeholder="Coins" v-model.number="offer.coins"/>
      <input type="number" class="form-control" placeholder="Price" v-model.number="offer.price"/>
      <div class="input-group-append">
        <button class="btn btn-outline-success" type="button" @click="sendOffer">Add</button>
      </div>
    </div>
  </div>
</template>

<script>
    export default {
        data(){
            return{
                offer: {
                    coins: '',
                    price: '',
                },
            }
        },
        methods:{
            sendOffer: function () {
                let offer = this.offer;
                if(offer.coins <= 0 || offer.price <= 0) {
                    return
                }

                this.$emit('newOffer', offer)
            },
        },
        name: "new-offer"
    }
</script>

<style scoped>
.offer{
  margin-top:20px;
}
</style>