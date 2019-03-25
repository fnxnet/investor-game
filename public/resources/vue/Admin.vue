<template>
  <div>
    <div class="row" style="padding-top:20px;">
      <div class="col-6">
      </div>
      <div class="col-6 text-right">
        <div class="input-group mb-4">
          <input type="password" class="form-control" v-model="password" placeholder="password">
          <div class="input-group-append">
            <button class="btn btn-info" @click="loginAdmin">Login</button>
            <button class="btn btn-danger" @click="clearDb">Clear DB</button>
            <button type="button" :class="['btn', lockClass]" @click="swapLock">{{lockLabel}}</button>
          </div>
        </div>
      </div>
    </div>

    <!--<div class="row lower">-->
    <!--<div class="col-6">-->
    <!--<h2 @click="firstIteration=!firstIteration">First Iteration</h2>-->
    <!--</div>-->
    <!--</div>-->
    <!--<div class="row" v-if="firstIteration">-->
    <!--<div class="col">-->
    <!--<new-offer @newOffer="newOffer"></new-offer>-->
    <!--</div>-->
    <!--</div>-->

    <div class="row lower">
      <div class="col">
        <h2 @click="showStats=!showStats">Stats</h2>
      </div>
    </div>

    <div class="row" :class="['divider', {'d-none':!showStats}]">
      <div class="stats col-5">
        <div>
          <h3>Transactions: {{stats.offers.length}}</h3>
          <h3>Coins: {{sumCoins}}</h3>
          <h3>Users: {{userCount}}</h3>
          <hr/>
          <h3>Cash spend: {{sumOffers | toCurrency}}</h3>
          <h3><b>BestOffer: {{bestOffer | toCurrency}}</b></h3>
          <h3><u>Average: {{avgAmount | toCurrency}}</u></h3>
        </div>
      </div>
      <div class="col-7">
        <div class="chart" style="width:100%;height:300px">
          <canvas id="myChart"></canvas>
        </div>
      </div>
    </div>

    <div class="row lower">
      <div class="col">
        <h2 @click="income=!income">Share income</h2>
      </div>
    </div>

    <div class="row " v-if="income">
      <div class="col">
        <div class="input-group mb-4">
          <input type="number" class="form-control" placeholder="Previous" v-model.number="shares.prev"/>
          <div class="input-group-prepend" @click="swapShares">
            <button class="btn btn-outline-info" type="button">Swap</button>
          </div>
          <input type="number" class="form-control" placeholder="Current" v-model.number="shares.current"/>
          <div class="input-group-append" @click="shareIncome">
            <button class="btn btn-outline-success" type="button">Share Income</button>
          </div>

        </div>
      </div>
    </div>

    <div class="row lower divider" v-if="income">
      <div class="col-6">
        <h2><span :style="{ color: diffColor}">Diff: {{shareDiff}}</span></h2>
      </div>
      <div class="col-6">
        <h2><span style="color: green;"> Share: {{shareAmount}}</span></h2>
      </div>
    </div>

    <div class="row">
      <div class="col">
        <h2 @click="getTheWinner">Find Winners!!!</h2>
      </div>
    </div>

    <div v-if="winners.length" class="alert alert-success">
      <li v-for="winner in winners">{{winner.login}} ({{winner.cash|amount}})</li>
    </div>

  </div>
</template>

<script>
    import axios from 'axios'
    import NewOffer from './NewOffer'

    export default {
        data()
        {
            return {
                firstIteration: false,
                showStats: false,
                controls: true,
                income: false,
                locked: false,
                chart: '',
                password: '',
                winners: [],
                statsInterval: false,
                shares: {
                    current: '',
                    prev: '',
                },
                offer: {
                    price: '',
                    amount: '',
                }
            }
        },
        filters: {
            amount(value, decimals)
            {
                if(decimals == undefined) {
                    decimals = 0
                }

                let pow = Math.pow(10, decimals);
                return Math.ceil(value * pow) / pow
            }
        },
        computed: {
            userCount()
            {
                return this.$store.state.stats.user_count;
            },
            lockLabel()
            {
                return this.locked ? "Unlock" : "Lock"
            },
            stats()
            {
                return this.$store.state.stats;
            },
            sumOffers()
            {
                console.log("sumOffers");
                let sum = 0;

                this.stats.offers.map(offer => {
                    sum += offer.coins * offer.price
                });

                return sum
            },
            avgAmount()
            {
                let total = this.sumOffers;
                let len = this.sumCoins;

                return total && len ? total / len : 0;
            },
            sumCoins()
            {
                let total = 0;

                this.stats.offers.map(offer => {
                    total += offer.coins
                });

                return total;
            },
            globalOffers()
            {
                let offers = [];

                this.stats.offers.map(offer => {
                    offers.push(offer.price)
                });

                return offers;
            },
            bestOffer()
            {
                let best = 0;

                this.stats.offers.map(offer => {
                    if(offer.price >= best) {
                        best = offer.price;
                    }
                });

                return best
            },
            diffColor()
            {
                return this.shareDiff < 0 ? 'red' : ''
            },
            shareDiff()
            {
                return this.shares.current - this.shares.prev
            },
            shareAmount()
            {
                if(this.shareDiff <= 0) {
                    return 0;
                }

                return Math.ceil(this.shareDiff / 1000 * 1000) / 1000
            },
            lockClass()
            {
                return this.locked ? 'btn-danger' : 'btn-success'
            },
            avg()
            {
                let chunk = this.stats.stats.chunk;
                let avg = chunk.TotalCash && chunk.CoinsSold ? chunk.TotalCash / chunk.CoinsSold : 0

                return Math.ceil(avg * 100) / 100
            }
        },
        mounted()
        {
            var ctx = document.getElementById("myChart").getContext('2d');
            this.chart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: Array.from(Array(this.stats.offers.length).keys()),
                    datasets: [
                        {
                            label: 'Price',
                            data: this.globalOffers,
                            backgroundColor: [
                                'rgba(255, 255, 255, 0.2)',
                            ],
                            borderColor: [
                                'rgba(255,99,132,1)',
                            ],
                            borderWidth: 1
                        }
                    ]
                },
                options: {
                    scales: {
                        yAxes: [{
                            ticks: {
                                beginAtZero: true
                            }
                        }]
                    }
                }
            });

            setInterval(this.updateStats, 1000);
        },
        methods: {
            toggleStats()
            {
                if(!this.statsInterval) {
                    this.updateStats();
                    this.statsInterval = setInterval(this.updateStats, 3000)
                }
                else {
                    clearInterval(this.statsInterval)
                    this.statsInterval = false
                }
            },
            updateStats()
            {
                console.log('update chart');
                let sets = this.chart.data.datasets;
                this.chart.data.labels = Array.from(Array(this.stats.offers.length).keys());

                sets[0].data = this.globalOffers;
                sets[0].data = this.globalOffers;

                this.chart.update();
            },
            swapLock()
            {
                let url = this.$store.state.url_prefix + '/admin/' + (this.locked ? 'unlock' : 'lock');
                axios.post(url, {password: this.password}).then(data => {
                    console.log(data);
                    this.locked = !this.locked;
                }).catch(error => {
                    alert('Ups');
                });

            },
            swapShares()
            {
                let tmp = this.shares.current;
                this.shares.current = this.shares.prev;
                this.shares.prev = tmp
            },
            newOffer(offer)
            {
                let message = {
                    callback: "offerReceived",
                    action: 'newOffer',
                    payload: offer,
                    user: this.$store.state.user,
                };

                this.$socket.sendObj(message);
            },
            shareIncome()
            {
                let url = this.$store.state.url_prefix + '/admin/share';
                let data = {
                    password: this.password,
                    current: this.shares.current,
                    prev: this.shares.prev,
                };

                axios.post(url, data).then(data => {
                    this.shares = {
                        current: '',
                        prev: this.shares.current,
                    }
                }).catch(error => {
                    alert('Ups!');
                    console.log(error);
                });
            },
            getTheWinner()
            {
                let url = this.$store.state.url_prefix + '/admin/best';
                let data = {
                    password: this.password,
                };

                axios.post(url, data).then(data => {
                    this.winners = data.data || [];
                }).catch(error => {
                    console.log(error);
                });
            },
            loginAdmin()
            {
                let url = this.$store.state.url_prefix + '/register-admin';
                axios.post(url, {password: this.password}).then(data => {
                    this.$store.commit('register', data.data).then(() => {
                        console.log('success')
                    })
                }).catch(error => {
                    if(error.response === undefined) {
                        return
                    }

                    let status = error.response.status || 0;
                    switch (status) {
                        case 409:
                            alert('Registered');
                            break;
                        case 500:
                            alert('Wrong password');
                    }
                });
            },
            clearDb()
            {
                let url = this.$store.state.url_prefix + '/admin/clear';
                axios.post(url, {password: this.password}).then(data => {
                    this.controls = false;

                    this.$store.commit("clearStats");

                }).catch(error => {
                    alert('invalid data');
                    console.log(error)
                });
            },
        },
        components: {
            NewOffer,
        }
    }
</script>
<style scoped>
  .row.lower {
    margin-bottom: 30px;
  }

  .divider {
    border-bottom: 1px solid #ccc;
  }
</style>