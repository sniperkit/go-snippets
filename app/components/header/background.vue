<template>
    <div>
          <h3 class="title">Background</h3> 
            <div class="item">       
                <span>Bing</span>
                <Choose 
                v-if='bing'
                :initial='bing.status'
                @hanlder='bingHanlder'/>                          
           </div>          
           <h3 class="title">Pure</h3>
           <ul class="pure">
               <li                      
                    v-for="(color,idx) in pure"
                    :key="idx"
                    :style="background(color)"
                    @click="colorHanlder(color)"
                    >
               </li>
           </ul>

    </div>
</template>

<script>
import Choose from "../common/choose.vue"
import { mapGetters } from 'vuex'
import { colors } from '../../api/setting.js'

export default {
    data(){
        return{                
            pure:[]
        }
    },
    created(){ 
        this.pure = colors
    },
    components:{
        Choose
    }, 
    computed:{              
        ...mapGetters(['bing']),
    },    
    methods:{
        background(color){
            return{
                background : color 
            }
        },
        bingHanlder(){     
            this.$store.dispatch('modifyBackgroundVisible',!this.bing.status)
        },
        colorHanlder(color){
            this.$store.dispatch('modifyColor',color)
        }       
    }    
}
</script>


<style lang="scss" scoped>
.pure{
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    margin: 2px 0;
    & >li{
        display: inline-block;
        width: 22px;
        height: 22px;
        margin: 0 4px 4px 4px;
        border: 1px solid #eee;
        border-radius: 2px;
    }
}
</style>
