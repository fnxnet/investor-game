<template>
  <div v-if="value.display">
    <div class="modal-background">
      <div class="modal show" tabindex="-1" role="dialog">
        <div class="modal-dialog" role="document">
          <div class="modal-content">
            <div v-if="value.close && value.title" class="modal-header">
              <h5 class="modal-title">{{value.title}}</h5>
              <button type="button" class="close" data-dismiss="modal" @click="closeModal" aria-label="Close">
                <span aria-hidden="true">&times;</span>
              </button>
            </div>
            <div class="modal-body">
              <button v-if="value.close && !value.title" type="button" class="close top" data-dismiss="modal"
                      @click="closeModal" aria-label="Close">
                <span aria-hidden="true">&times;</span>
              </button>
              <slot></slot>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
    export default {
        name: "modal",
        data()
        {
            return {}
        },
        computed: {},
        props: {
            title: String,
            display: Boolean,
            close: Boolean,
            value: Object,
        },
        methods: {
            closeModal()
            {
                this.$emit('close')
            },
        },
    }
</script>

<style scoped>
  .modal {
    display: table-cell;
    vertical-align: middle;
  }

  .close.top {
    display: block;
    position: absolute;
    top: 0;
    right: 20px;
  }

  .modal-background {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, .5);
    display: table;
    transition: opacity .3s ease;
  }
</style>