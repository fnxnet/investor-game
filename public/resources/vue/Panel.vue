<template>
  <div>

    <div class="row lower">
      <div class="col">
        <h2>Stats</h2>
      </div>
    </div>

    <div class="row">
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
                showStats: true,
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
                let sets = this.chart.data.datasets;
                this.chart.data.labels = Array.from(Array(this.stats.offers.length).keys());

                sets[0].data = this.globalOffers;
                sets[0].data = this.globalOffers;

                this.chart.update();
            },
            swapLock()
            {
                let url = this.$store.state.url_prefix + '/panel/' + (this.locked ? 'unlock' : 'lock');
                axios.post(url, {password: this.password}).then(data => {
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
            shareIncome(amount)
            {
                let url = this.$store.state.url_prefix + '/panel/share';
                let data = {
                    password: this.password,
                    // current: this.shares.current,
                    // prev: this.shares.prev,
                    diff: amount,
                };

                axios.post(url, data).then(data => {
                    console.log('shared', data);
                }).catch(error => {
                    alert('Ups!');
                });
            },
            getTheWinner()
            {
                let url = this.$store.state.url_prefix + '/panel/best';
                let data = {
                    password: this.password,
                };

                axios.post(url, data).then(data => {
                    this.winners = data.data || [];
                }).catch(error => {
                    alert('Ups!');
                });
            },
            loginAdmin()
            {
                let url = this.$store.state.url_prefix + '/register-admin';
                axios.post(url, {password: this.password}).then(data => {
                    this.$store.commit('register', data.data);
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
                let url = this.$store.state.url_prefix + '/panel/clear';
                axios.post(url, {password: this.password}).then(() => {
                    this.controls = false;

                    this.$store.commit("clearStats");

                }).catch(error => {
                    alert('Ups!');
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