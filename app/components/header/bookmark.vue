<template> 
    <div class="item">
          <i 
            @click="bookmarkHandler"
            class="material-icons"
            >{{option.icon}}</i>
            
     <transition name="top">
        <div class="bookmark dialog"
            v-if="visible">          

          <div class="search">
            <input type="text" autofocus 
            v-model="search_text"
            placeholder="search..."/> 
            <span class="button" @click="close">Dismiss</span>
            
          </div>   

          <div class="inner">
            <div class="category">
                  <p
                    v-for="category in categories"
                    :key="category.id"
                    @click="changeCategory(category.id)" >
                    {{ category.title }}
                  </p>
            </div>

            <div class="item">
                <a                 
                    v-for="(bookmark,idx) in current"  
                    :key="idx"
                    :href="bookmark.url"   
                    target="_blank"      
                > 
                <img :src="bookmark.icon"/>

                <span class="item-title">{{bookmark.title}}</span>  
                </a>
            </div>
            </div>

       
        </div>  
     </transition>

    </div>
</template>

<script>
import * as bookmarkAPI from "../../api/bookmarks"

export default {
    props:{
        option:{
            type: Object,
            require: true
        },
        visible:{
            type: Boolean,
            require: true
        }
    },
  data() {
    return {
      id: 1,
      categories: [],
      all: [],
      search_text: "",
      search:[],
      showSearch: false
    }
  },
  created() {
    bookmarkAPI.getCategories.then(data => {
      let i = 1
      // get all categories
      this.categories = data
      // get all bookmarks
      this.categories.forEach(category => {
        bookmarkAPI.getBookMarksById(category.id).then(data => {
          this.$set(this.all, i++, data)
        })
      })
    })
  },
  computed: {
    current() {      
      if(this.showSearch){                  
        return this.search
      }    
        return this.all[this.id]      
    }
  },
  methods: {
    close() {
      this.$emit("close", "bookmark")
    },
    bookmarkHandler(){
      this.$emit("bookmarkHandler")
    },
    changeCategory(id) {
      this.showSearch = false
      this.id = id     
    }
  },
  watch:{
    search_text(val){      
      if( this.search_text === ''){
        this.showSearch = false
      }else{   
        this.showSearch = true  
        this.search.length = 0

        this.all.forEach(cate=>{
            cate.forEach(bookmark=>{
                if(bookmark.title.toLowerCase().indexOf(val)!=-1){
                  this.search.push(bookmark)
                }
            })
        })      
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.bookmark {
  width: 680px;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  & > .search{
    width: 100%;
    height: 44px;
    padding-left: 160px;
    & > input{
      width: 360px;
      padding: 0 15px;
      margin-right: 32px;
      line-height: 36px;
      border: 1px solid #eee;
      border-radius: 4px;
    }
  }
  & > .inner{      
    display: flex;
    margin-top: 20px;
    & > .category {
      min-width: 140px;
      & > p {        
        width: 100%;
        padding: 2px 12px;
        margin-bottom: 30px;
        font-size: 14px;
        color: #3b444f;

        &:hover{
          color: rgba(0,0,0,.54)
        }
      }
    }
    & > .item {      
      display: flex;
      flex-wrap: wrap;  
      & > a {
        display: flex;
        align-items: center;
        width: 100%;        
        margin: 0 10px 12px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
        font-size: 14px;      
        &:hover{
          color: rgba(0,0,0,.54)
        }
        & >.item-title{
          display: inline-block;
          margin-left: 8px;
          text-align: center;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
          color: #252b32;
        }
      }
  }
  }
}
</style>
