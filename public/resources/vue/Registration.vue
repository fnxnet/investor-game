<template>
  <div>
    <modal v-model="modal">
      <div class="row justify-content-md-center">
        <div class="col-sm-8">
          <div class="input-group">
            <input type="text" class="form-control" placeholder="login" v-model="registrationLogin"
                   aria-label="Input group example" aria-describedby="btnGroupAddon">
            <div class="input-group-append">
              <div class="btn btn-success" id="btnGroupAddon" @click="register">Register</div>
            </div>
          </div>
        </div>
        <div class="alert alert-danger errorMessage" @click="registrationMessage=''" v-if="registrationMessage">
          {{registrationMessage}}
        </div>
      </div>
    </modal>
  </div>
</template>

<script>
    import axios from 'axios';
    import Modal from './Modal';

    export default {
        name: "registration",
        data()
        {
            return {
                registrationLogin: '',
                registrationMessage: '',
                modal: {
                    display: true
                }
            }
        },
        methods: {
            register()
            {
                if(this.registrationLogin.length < 3) {
                    this.registrationMessage = "You need to oprovide login min 3 characters long";
                    return
                }
                this.registrationMessage = '';
                axios.post(this.$store.state.url_prefix + '/register', {login: this.registrationLogin})
                    .then((r) => {
                        this.$store.dispatch('refreshOffers').then(() => {
                            this.$store.commit('register', r.data)
                        });
                    })
                    .catch(error => {

                        if(error.response === undefined) {
                            return
                        }

                        let status = error.response.status || 0;

                        switch (status) {
                            case 409:
                                this.registrationMessage = "Login you used is already taken please use different one"
                        }
                    })
            }
        },
        components: {
            Modal,
        }
    }
</script>

<style scoped>
  .errorMessage {
    margin-top: 20px;
    cursor: pointer;
  }
</style>