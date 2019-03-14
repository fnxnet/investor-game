<template>
  <div>
    <div class="filters" v-if="type!='own'">
      <label @click="inactive=!inactive">Show Inactive</label>
      <input type="checkbox" v-model="inactive"/>
    </div>
    <h3 class="head" @click="display=!display">
      <span class="icon"><i :class="['fa', chevronClass]"/></span> {{type=='own'?'Own':'Others'}} ({{offers.length}})
    </h3>
    <table v-if=display :class="['table','text-center',type]" style="width: 100%;">
      <thead class="thead-dark">
      <tr>
        <th scope="col">Coins</th>
        <th scope="col" @click="swapSort('price')">
          <i :class="['fa',sortDirectionClass]" v-if="sort=='price'"></i> Price
        </th>
        <th scope="col" @click="swapSort('total')">
          <i :class="['fa',sortDirectionClass]" v-if="sort=='total'"></i> Total
        </th>
        <th scope="col" width="50%">Action</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="offer in offers" :class="{inactive:!offer.active}">
        <td>{{offer.coins}}</td>
        <td>{{offer.price}}</td>
        <td>{{offer.price*offer.coins}}</td>
        <td>
          <div class="btn btn-dark" v-if="type=='own'" @click="removeOffer(offer)">remove</div>
          <div class="btn btn-info" v-else-if="offer.active" @click="buyOne(offer)">Buy</div>
        </td>
      </tr>
      </tbody>
    </table>
    <modal v-model="offerModal" @close="offerModal.display=false">
      <div class="row justify-content-md-center">
        <div class="col-12">
          <div class="input-group">
            <input type="number" v-model.number="coins" @change="validateCoins" class="form-control">
            <input disabled class="form-control" :value="totalOffer">
            <div class="input-group-append">
              <span @click="increase" class="controls input-group-text">+</span>
              <span @click="decrease" class="controls input-group-text">-</span>
              <div class="btn btn-info" @click="buyPart">Buy</div>
            </div>
          </div>
        </div>
      </div>
    </modal>
  </div>
</template>

<script>
    import Modal from "./Modal";

    export default {
        data()
        {
            return {
                display: true,
                inactive: true,
                sort: '',
                sortDirection: 'up',
                coins: '',
                offerModal: {
                    display: false,
                    title: "Select amount",
                    close: true,
                    offer: {}
                }
            }
        },
        props: {
            type: String
        },
        computed: {
            sortDirectionClass()
            {
                return 'fa-chevron-' + this.sortDirection
            },
            chevronClass()
            {
                return 'fa-chevron-' + (this.display? 'down' : 'up')
            },
            totalOffer()
            {
                return 'Total: ' + this.offerModal.offer.price * this.coins + ' $'
            },
            user()
            {
                return this.$store.state.user
            },
            offers()
            {
                let offers = this.$store.state.offers.filter(
                    offer => this.type == 'own' ? offer.my : !offer.my
                ).filter(
                    offer => !offer.active && this.type == 'own' ? false : this.inactive || offer.active
                );

                offers.sort((a, b) => {
                    let a_total = a.price * a.coins;
                    let b_total = b.price * b.coins;

                    if(this.sort == 'price') {
                        return a.price > b.price ? 1 : (a.price < b.price ? -1 : 0)
                    }

                    return a_total > b_total ? 1 : (a_total < b_total ? -1 : 0)
                });

                if(this.sortDirection == 'down') {
                    offers.reverse()
                }

                return offers;
            }
        },
        methods: {
            swapSort(type)
            {
                this.sort = type;
                this.sortDirection = this.sortDirection == 'down' ? 'up' : 'down'
            },
            increase()
            {
                this.coins++;
                this.validateCoins()
            },
            decrease()
            {
                this.coins--;
                this.validateCoins()
            },
            validateCoins()
            {
                this.coins = this.coins <= this.offerModal.offer.coins ? this.coins : this.offerModal.offer.coins;
                this.coins = this.coins >= 1 ? this.coins : 1;
            },
            removeOffer(offer)
            {
                let message = {
                    callback: "offerRemoved",
                    action: 'removeOffer',
                    payload: offer,
                    user: this.user,
                };

                this.$socket.sendObj(message);

            },
            buyOne(offer)
            {
                if(offer.coins == 1) {
                    return this.buyAll(offer);
                }

                console.log('dasad')
                this.offerModal.display = true;
                this.offerModal.offer = offer;
                this.coins = offer.coins;
            },
            buyPart()
            {
                let o = JSON.parse(JSON.stringify(this.offerModal.offer));
                console.log('buy part', o);
                o.coins = this.coins;
                this.buyAll(o)
                // console.log(o)
                this.offerModal.display = false
            },
            buyAll(offer)
            {
                let o = JSON.parse(JSON.stringify(offer));
                // o.coins = 1;
                delete o.active;
                let message = {
                    callback: "offerBought",
                    action: 'acceptOffer',
                    payload: o,
                    user: this.user,
                };

                this.$socket.sendObj(message);
            },
            test()
            {
                console.log('test')
            }
        },
        components: {
            Modal,
        },
        name: "offer-list"
    }
</script>
<style>
  tr.inactive td {
    background-color: lightgrey;
  }

  .own thead.thead-dark th {
    background-color: #17a2b8;
    border-color: #17a2b8;
  }

  .controls {
    cursor: pointer;
  }

  .input[disabled] {
    background-color: #fff;
  }

  .table td {
    padding: .25rem;
    vertical-align: middle;
  }

  .table td .btn {
    padding: 0.25rem 2rem;
    vertical-align: middle;
  }

  .head .icon .fa{
    margin-top: -5px;
    font-size: 16px;
  }

</style>