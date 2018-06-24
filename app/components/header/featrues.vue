<template>
  <div>
       <h3 class="title">Features</h3>     

            <div 
                class="item"
                v-for="(item,idx) in featrues"
                :key="item.name">   
                <span>{{item.text}}</span>
                <Choose
                    :id = "idx"
                    :initial="item.status"
                    @hanlder='featruesHanlder'
                ></Choose>
            </div>
  </div>
</template>

<script>
import Choose from "../common/choose.vue"
import {mapGetters} from 'vuex'

export default {   
    computed:{
      ...mapGetters(['featrues'])
    },
    methods:{        
        featruesHanlder(index) {
            let item = this.featrues[index]   
            this.$store.dispatch('modifyFeatrues',{type: item.name, status: !item.status})     
        }  
    },
    components:{
        Choose
    }
}
</script>
