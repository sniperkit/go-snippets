<template>
  <div class="overlay"
    @click="click"
  ></div>
</template>

<script>
export default {
    
    methods:{
         click(){
               this.$emit('overlay')
         }       
    }
}
</script>

<style lang="scss" scoped>
.overlay{
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.84);
    z-index: 8;
}

</style>