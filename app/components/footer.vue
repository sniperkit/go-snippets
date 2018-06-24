<template>
    <footer class="footer">
        
        <div class="item">           
          <i @click="opGallery"
            class="material-icons"
            >images</i>
          
        </div>
        <div class="item"
                v-if="bing && bing.status">   
            <i @click="opImageInfo"
                class="material-icons">error_outline</i>
            <Info
                :show="info"
                :background = "background"
                @dismiss ="opImageInfo"
            />
        </div>

        <Background
            :background = "background"
         />
        <Gallery 
            v-show="gallery"
            :images = "images"
            @setImage ="setImage"  
            @dismiss = "opGallery"        
        />
    </footer>
  
</template>


<script>
import Background from './footer/background.vue'
import Gallery from './footer/gallery.vue'
import Info from './footer/info.vue'

import {mapGetters} from 'vuex'
import * as BingAPI from '../api/bing.js'

export default { 
  data(){
      return{
          images: [],        
          background: {}, 
          gallery: false,
          info: false
      }
  },
  created(){
     BingAPI.getBingImages.then( data =>{       
         this.images = data
     }),
     BingAPI.getBackground.then( data =>{       
         this.background = data     
    })
  },
  components:{
      Background,Gallery,Info
  },
  computed:{              
        ...mapGetters(['bing'])
  }, 
  methods:{
      opGallery(){
          this.info = false
          this.gallery = !this.gallery 
      },    
      opImageInfo(){
          this.gallery = false
          this.info = !this.info 
      },
      setImage(image){
          // modify vuex
          this.$store.dispatch('modifyBackground',image)   
          // set this.
          this.$set(this.background,"url",image.url)
          this.$set(this.background,"date",image.date)
          this.$set(this.background,"copyright",image.copyright)
             
      }
  }
}
</script>

<style lang="scss" scoped>
.footer{
    width: 100%;
    height: 60px;      
    padding: 16px 24px 24px 24px;
    
    display: flex;
    flex-direction: row-reverse;
    align-items: center;
    & >.item{
        position: relative;
        display: inline-block;
        width: 32px;
        height: 32px;
        margin: 0 6px;
        text-align: center;
        cursor: pointer;       
    }
}
</style>