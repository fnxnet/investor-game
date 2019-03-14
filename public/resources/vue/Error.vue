<template>
  <modal v-model="modal" v-if="message" @close="close">
    <h4 class="errorMessage text-center">{{message}}</h4>
  </modal>
</template>

<script>
    import Modal from './Modal'

    export default {
        data()
        {
            return {
                modal: {
                    display: true,
                    close: true,
                }
            }
        },
        computed: {
            message()
            {
                return this.$store.state.actionMessage
            },
        },
        methods: {
            close(){
                this.$store.commit('hideErrorMessage')
            }
        },
        components: {
            Modal
        },
        name: "error"
    }
</script>

<style scoped>

</style>