<template>
  <transition name="right-in-left-out">
    <div class="box"
        v-show="show"
    >
        <img :src="background.url" >
        <span>{{formatDate}}</span>

        <span
            v-for="(c,idx) in formatCopyRihgt"
            :key="idx">{{c}}</span>
            
            <div class="bottom">   
              <span 
                  class="button download"
                  @click.stop="download(background.url)"

              >Download</span>                
              <span 
                  class="button"
                  @click.stop="dismiss"            
              >Dismiss</span>   
            </div>
    </div> 

  </transition>
</template>

<script>
import download from '../../api/download.js'

export default {
    props:{
        show:{
            type: Boolean,
            required: true
        },
        background:{
            type : Object,
            required: true
        }        
    },
    computed:{
        formatDate: function (){    
            let date = this.background.date
            if(date){
                 return `${date.substring(0,4)}/${date.substring(4,6)}/${date.substring(6,8)}`
            }            
        },
        formatCopyRihgt: function (){
            let cpr = this.background.copyright
            if(cpr){
                cpr = cpr.split('，')
            }
            return cpr
        }
    },
    methods:{
        dismiss(){
            this.$emit('dismiss')
        },
        download(url){
            download(url)
            this.dismiss()
        }
    }
  
}
</script>

<style lang="scss" scoped>
.box{
    width: 400px;
    min-height: 53vh;
    position: fixed;
    right: 5px;
    bottom: 60px;   
    background: #fff;
    border: 1px solid #eee;
    border-radius: 4px;
    z-index: 99;
    padding: 0 0 20px; 
    & >img{    
        width: 100%;
        height: 34vh;
        padding: 6px;
    }
    & >span{
        display: block;
        margin: 2px 0 4px;
        color: #252b32;
        padding: 0 2px;
    }
    & >.bottom{
        margin-top: 15px;       
        & .download{
            margin-right: 8px;
        }
    }
}

</style>
