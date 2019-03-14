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

    <div class="row" :class="{'d-none':!showStats}">
      <div class="stats col-4">
        <div class="stats" v-if="stats">
          <h3>Cash spend: {{stats.stats.chunk.TotalCash}}</h3>
          <h3>Coins: {{stats.stats.chunk.CoinsSold}}</h3>
          <h3><u>Average: {{avg?avg:0}}</u></h3>
          <h3>BestOffer: {{stats.best_offer}}</h3>
          <h3>Users: {{stats.user_count}}</h3>
        </div>
        <button type="button" class="btn btn-danger" @click="resetStats">ResetStats</button>
        <button v-if="!statsInterval" type="button" class="btn btn-success" @click="toggleStats">Start stats</button>
        <button v-else type="button" class="btn btn-warning" @click="toggleStats">Stop stats</button>
      </div>
      <div class="col-8">
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

    <div class="row" v-if="income">
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

    <div class="row lower" v-if="income">
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
                stats: {
                    stats: {
                        chunk: {
                            TotalCash: '',
                            CoinsSold: '',
                        },
                        global: {
                            TotalCash: '',
                            CoinsSold: '',
                        },
                    },
                    user_count: "",
                    best_offer: ""
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
            lockLabel()
            {
                return this.locked ? "Unlock" : "Lock"
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
                    labels: [],
                    datasets: [
                        {
                            label: 'Avg Price',
                            data: [],
                            backgroundColor: [
                                'rgba(255, 255, 255, 0.2)',
                            ],
                            borderColor: [
                                'rgba(255,99,132,1)',
                            ],
                            borderWidth: 1
                        },
                        {
                            label: 'Users',
                            data: [],
                            backgroundColor: [
                                'rgba(255, 255, 255, 0.2)',
                            ],
                            borderColor: [
                                'rgba(8,100,149,1)',
                            ],
                            borderWidth: 1
                        },
                        {
                            label: 'Best',
                            data: [],
                            backgroundColor: [
                                'rgba(255, 255, 255, 0.2)',
                            ],
                            borderColor: [
                                'rgba(8, 181, 135,1)',
                            ],
                            borderWidth: 1
                        },
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
                axios.post(this.$store.state.url_prefix + '/admin/stats', {password: this.password}).then(data => {
                    this.stats = data.data;

                    console.log(data.data);

                    let sets = this.chart.data.datasets;
                    this.chart.data.labels.push(null);

                    console.log(this.stats.user_count);
                    sets[0].data.push([this.avg]);
                    sets[1].data.push([this.stats.user_count]);
                    sets[2].data.push([this.stats.best_offer]);

                    if(sets[0].data.length > 15) {
                        this.chart.data.labels.shift();
                        sets.map(set => {
                            set.data.shift()
                        })
                    }
                    this.chart.update();
                });
            },
            resetStats()
            {
                let url = this.$store.state.url_prefix + '/admin/reset';

                axios.post(url, {password: this.password}).then(data => {
                    this.updateStats()
                }).catch(error => {
                    console.log(error);
                });
            },
            swapLock()
            {
                let url = this.$store.state.url_prefix + '/admin/' + (this.locked ? 'unlock' : 'lock');
                axios.post(url, {password: this.password}).then(data => {
                    console.log(data);
                    this.locked = !this.locked;
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
                        prev: '',
                    }
                }).catch(error => {
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
                }).catch(error => {
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
</style>