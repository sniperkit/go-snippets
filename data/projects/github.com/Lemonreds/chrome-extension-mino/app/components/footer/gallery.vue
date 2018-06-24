<template>
  <transition name="scale">
    <div class="gallery">
        <div class="bing">        
             <span>Images from</span>
             <a target="_blank" href="https://bing.com/">Bing daily wallpaper.</a>         
             <span 
                  class="button dismiss"
                  @click.stop="dismiss"            
              >Dismiss</span>     
        </div>
        <ul class="images">
            <li class="item"
                v-for="image in images"
                :key="image.name"
                :style="backgroundImage(image.url)">      


               <span class="options">
                    <i @click="download(image.url)"  class="material-icons">file_download</i>
                    <i @click="setImage(image)"  class="material-icons">photo</i>    
               </span>
                
                <span class="date">{{image.date}}</span>
                <p class="copyright"> {{image.copyright}}</p>
                 
            </li>
        </ul>  

    </div>
  </transition>
</template>

<script>

import download from '../../api/download.js'

export default { 
  props:{
     images: {
         type: Array,
         required: true
     }
  }, 
  methods:{
      backgroundImage :function(src){
           return {
              backgroundImage: `url(${src})`
           }    
      },
      setImage: function(image){
          this.$emit("setImage",image)
      },
      download: function(url){
          download(url)    
      },
      dismiss: function(){
          this.$emit("dismiss")
      }
  }
}
</script>

<style lang="scss" scoped>
.gallery{
    position: fixed;
    width: 860px;
    height: 510px;
    top: 50%;
    left: 50%;
    margin-left: calc(-860px/2);
    margin-top: calc(-510px/2);
    background: #fff;
    padding: 20px;
    border-color: rgba(0,0,0,.2);
    color: #3b444f;
    box-shadow: 0 2px 10px rgba(0,0,0,.2);    
    border-radius: 4px;  
    z-index: 99;

}
.bing{
    font-size: 18px;
    margin-bottom: 24px;
    text-align: center;
    & >span{
        font-size: 14px
    }
    & >a{
        font-size: 14px;
        text-decoration: underline;
        color: rgba(0, 0, 0, 1)
    }
    & >.dismiss{
        position: absolute;
        top: 22px;
        right: 24px;
    }
}
.images{    
    display: flex;
    flex-wrap: wrap;  

    & >.item{
        width: 202px;
        height: 202px;
        margin-right: 2px;
        margin-bottom: 4px;

        background-size: cover;
        background-position: center;
        background-repeat: no-repeat;

        text-align: center;
        cursor: pointer;
        position: relative;
        transition: all .2s ease-in-out;
        &>.date{
           line-height: 162px;
           color: #fff;
           opacity: 0;
        }
        &>.copyright{       
            position: absolute;
            bottom: 10px;
            left: 0;
            width: 100%;
            max-height: 49%;

            font-size: 12px;
            color: #fff;
            opacity: 0
        }
        &>.options{
            position: absolute;
            right: 8px;
            top: 8px;
            display: inline-block;
            height: 18px;
            color: #fff;
            opacity: 0;
            z-index: 9;           
        }
        &:hover{
            opacity: .91;
            &>.date ,&>.copyright,&>.options{
                opacity: 1;
            }
        }
    }
}
</style>
