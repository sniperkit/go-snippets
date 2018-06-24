<template>
  <div class="popup">
      <img id="icon" :src="icon"/>
      <span>mino built by vuejs,</span>
      <a target="_blank" href="https://github.com/Lemonreds/chrome-extension-mino">@see more.</a>
      <span id="button"
            @click="newtab"
      >New Tab Page</span>
  </div>
</template>
<script>
export default {
  data(){
      return{
          icon : './static/icons/icon-128.png'
      }
  },
  methods:{
      newtab(){
        chrome.tabs.create({})
      }
  }
}
</script>

<style lang="scss" scoped>
.popup{
    width: 140px;
    height: 110px;
    font-size: 13px;   
    color: rgba(0,0,0,.74); 
    line-height: 18px;
    &>a{
        text-decoration: underline;
        color: rgba(0,0,0,1);
        transition: all .2s ease-in-out;
        &:hover{
            color: rgba(0,0,0,.84)
        }   
    }
    & >#icon{
        width: 42px;
        height: 42px;

        display: block;
        margin: 12px auto;
    }
    & >#button{
        width: 100px;

        font: 400 13px 'Arial';
        color: #fff;
        text-align: center;

        display: block;
        margin: 20px auto;
        padding: 4px;

        background: #46b8da;
        border: 1px solid #5bc0de;
        border-radius: 4px;

        cursor: pointer;
        transition: all .3s ease-in-out;
    }
}
</style>
