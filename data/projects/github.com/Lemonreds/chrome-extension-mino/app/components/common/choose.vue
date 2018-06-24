<template>
    <span class="choose" 
          @click="_emit"
         :class="initial ? 'active' : 'inactive'">
        <span class="cirlce"/>
    </span>
</template>

<script>
export default {
  props:{
      initial:{
          type: Boolean,
          default: true
      },
      id:{       
          default: 0
      }
  },
  data(){
     return{    
          status: true  
     }
  },  
  methods:{
      _emit(){
          this.$emit('hanlder',this.id)
      }
  }
}
</script>

<style lang="scss" scoped>
.choose{
    width: 40px;
    height: 24px;
    display: inline-block;
    position: relative;
    border: 1px solid  #eee;
    border-radius: 999px;  
    transition: all .1s ease-in-out;
    & > .cirlce{
        position: absolute;      
        width: 50%;
        height: 100%;
        top: 0;
        right: 0;
        background: #f5f5f5;    
        border-radius: 50%;
        transition: all .1s ease-in-out;                
    }
}
.active{
  background: #5cb85c; 
  & >.cirlce{
    transform: translateX(0);  
  }
}
.inactive{
  background: #d9534f;  
  & >.cirlce{
    transform: translateX(-20px);     
  }
}
</style>