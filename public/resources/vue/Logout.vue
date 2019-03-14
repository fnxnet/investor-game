<template>
  <div>
    Logged In as <b>{{login}}</b>
    <input type="button" value="Logout" @click="logout"/>
  </div>
</template>

<script>
    export default {
        name: "logout",
        computed: {
            login()
            {
                return this.$store.state.user.login
            }
        },
        methods: {
            logout: function () {
                this.$store.commit('logout')
            },
        }
    }
</script>

<style scoped>

</style>