<template>
  <div>
    <div    class="background fade"
            v-if="bing && bing.status"            
            :style="backgroundImage"           
            />
    <div class="background" 
            v-else :style="setBackgroundColor"     
    > </div>
   </div>
   
</template>

<script>
import { mapGetters } from 'vuex'

export default {
    props:{
        background:{
            type: Object,
            required: true
        }
    },        
    computed:{       
        backgroundImage :function(){
            let style = {}

            this.bing.inital && this.background.base64 ? style.backgroundImage = `url(${this.background.base64})`
                             : style.backgroundImage = `url(${this.background.url})`

            return style
        },
        setBackgroundColor:function(){
           
           return {
               background : this.color,
               opacity : 0.84
           }
        },
        ...mapGetters(['bing','color'])
    }
  
}
</script>

<style lang="scss" scoped>
.background{
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;  
    z-index: -1;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;     
    
}

@keyframes fade{
    from{       
         opacity: 0.1
    }
    to{       
         opacity: 1    
    }
}

.fade{
    animation-fill-mode: forwards;
    animation-name: fade;
    animation-duration: 0.2s 
}

</style>
