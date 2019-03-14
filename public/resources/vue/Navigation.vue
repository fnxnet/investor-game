<template>
  <div>
    <div class="nav">
      <div class="col-12 col-md-6">
        <logout></logout>
      </div>
      <div class="col-12 col-md-6 text-right">
        <user-stats></user-stats>
      </div>
    </div>
  </div>
</template>

<script>
    import Logout from './Logout'
    import UserStats from './UserStats'

    export default {
        name: "navigation",
        data()
        {
            return {}
        },
        components: {
            Logout,
            UserStats
        }
    }
</script>

<style scoped>
  .nav {
    padding-top: 10px;
  }

  .no-gutters {
    margin-right: 0;
    margin-left: 0;
  }

  .no-gutters.col,
  .no-gutters[class*="col-"] {
    padding-right: 0;
    padding-left: 0;
  }
</style>