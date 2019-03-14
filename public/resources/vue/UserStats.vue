<template>
  <div class="userData">
    Cash: {{user.cash}} | Coins: {{user.coins - user.coins_locked}} / {{user.coins}}
  </div>
</template>

<script>
    export default {
        data()
        {
            return {}
        },
        computed: {
            user()
            {
                return this.$store.state.user
            }
        },
        methods: {},
        name: "user-stats"
    }
</script>

<style scoped>
</style>