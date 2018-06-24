<template>    
      <transition name="scale">
        <div class="dialog engines">
            <p>Choose a search engine.</p>
            <div class="group">
                <div class="item"
                    v-for="(engine,index) in engines"
                    :key="engine.name"
                    @click="changeEngine(index)">
                    <img 
                        :src="engine.icon"
                        :alt="engine.name">
                </div>    
            </div>  
            <div class="bottom">
                <span class="button" @click="dismiss">Dismiss</span>                
            </div>  
        </div>
      </transition>
</template>

<script>
export default {
  props: {
    engines: {
      type: Array,
      required: true
    }
  },
  methods: {
    changeEngine(index) {
        this.$emit('changeEngine',index)
    },
    dismiss() {
        this.$emit('closeEngines')
    }
  }
};
</script>

<style lang="scss" scoped>
.engines {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 420px;
  height: 200px;
  margin-top: -100px;
  margin-left: -210px;
  &::before{
    display: none;
  }
  & >p{
    color: #46505e;
  }
  & > .group {
    margin: 26px 0;
    & > .item {
      width: 48px;
      padding: 6px;
      margin-right: 10px;
      display: inline-block;

      border: 1px solid #eee;
      border-radius: 4px;

      text-align: center;
      cursor: pointer;
      & > img {
        width: 36px;
        height: 36px;
      }
    }
  }
  & >.bottom{
      margin: 10px 0 0 ;
      font-size: 14px;
      text-align: right;
      cursor: pointer;

  }
}
</style>
